package client

import (
	"context"
	"fmt"

	"github.com/rhacm-global-hub-monitor/backend/pkg/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

// RHACMClient provides methods to interact with RHACM resources
type RHACMClient struct {
	kubeClient *KubeClient
}

// NewRHACMClient creates a new RHACM client
func NewRHACMClient(kubeClient *KubeClient) *RHACMClient {
	return &RHACMClient{
		kubeClient: kubeClient,
	}
}

// GetManagedHubs returns all managed hub clusters
func (r *RHACMClient) GetManagedHubs(ctx context.Context) ([]models.ManagedHub, error) {
	// Get all managed clusters once
	managedClusters, err := r.kubeClient.GetManagedClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get managed clusters: %w", err)
	}

	// Build hub-to-spokes mapping
	spokesMap := make(map[string][]models.ManagedCluster)
	var hubClusters []*clusterv1.ManagedCluster

	// First pass: identify hubs and build map
	for i := range managedClusters.Items {
		cluster := &managedClusters.Items[i]
		if isHub(*cluster) {
			hubClusters = append(hubClusters, cluster)
			spokesMap[cluster.Name] = []models.ManagedCluster{}
		}
	}

	// Second pass: assign spokes to hubs
	for i := range managedClusters.Items {
		cluster := &managedClusters.Items[i]
		if !isHub(*cluster) {
			// Find which hub manages this cluster
			for _, hub := range hubClusters {
				if belongsToHub(*cluster, hub.Name) {
					mc, err := r.convertToManagedCluster(ctx, cluster, hub.Name)
					if err == nil {
						spokesMap[hub.Name] = append(spokesMap[hub.Name], *mc)
					}
					break
				}
			}
		}
	}

	// Convert hubs to models and fetch their spoke clusters
	var hubs []models.ManagedHub
	for _, cluster := range hubClusters {
		// Try to fetch spoke clusters from this hub
		spokes := spokesMap[cluster.Name]

		// Attempt to connect to the hub and get its spoke clusters
		spokesFromHub, err := r.getSpokesClustersFromHub(ctx, cluster.Name)
		if err != nil {
			// Log error but continue - spokes may not be accessible
			fmt.Printf("Warning: Could not fetch spokes from hub %s: %v\n", cluster.Name, err)
		} else {
			spokes = spokesFromHub
		}

		// Fetch policies for this hub from its namespace on the global hub
		hubPolicies, err := r.kubeClient.GetPoliciesForNamespace(ctx, cluster.Name)
		if err != nil {
			fmt.Printf("Warning: Could not fetch policies for hub %s: %v\n", cluster.Name, err)
			hubPolicies = []models.PolicyInfo{}
		}

		// Fetch both K8s Nodes and BareMetalHost for this hub
		var hubNodes []models.NodeInfo

		// Try to connect to hub and get actual K8s nodes
		hubClient, err := NewHubClientFromSecret(ctx, r.kubeClient, cluster.Name)
		if err == nil {
			// Get K8s Node resources from the hub
			k8sNodes, err := hubClient.kubeClient.GetNodes(ctx)
			if err == nil {
				for i := range k8sNodes.Items {
					nodeInfo := ConvertNodeToNodeInfo(&k8sNodes.Items[i])
					if nodeInfo.Annotations == nil {
						nodeInfo.Annotations = make(map[string]string)
					}
					nodeInfo.Annotations["data-source"] = "Node"
					hubNodes = append(hubNodes, nodeInfo)
				}
			}
		}

		// Also get BareMetalHost resources from hub namespace on global hub
		bmhNodes, err := r.kubeClient.GetBareMetalHostsForNamespace(ctx, cluster.Name)
		if err == nil {
			hubNodes = append(hubNodes, bmhNodes...)
		}

		hub := &models.ManagedHub{
			Name:            cluster.Name,
			Namespace:       cluster.Name,
			Status:          getClusterStatus(cluster),
			Version:         getClusterVersion(cluster),
			Conditions:      convertConditions(cluster.Status.Conditions),
			ClusterInfo:     extractClusterInfo(cluster),
			NodesInfo:       hubNodes,
			PoliciesInfo:    hubPolicies,
			ManagedClusters: spokes,
			Labels:          cluster.Labels,
			Annotations:     cluster.Annotations,
			CreatedAt:       cluster.CreationTimestamp.Time,
		}
		hubs = append(hubs, *hub)
	}

	// Create set of existing hub names
	existingHubNames := make(map[string]bool)
	for _, hub := range hubs {
		existingHubNames[hub.Name] = true
	}

	// Also discover hubs from kubeconfig secrets (manually added hubs)
	unmanagedHubs, err := r.discoverUnmanagedHubs(ctx, existingHubNames)
	if err == nil {
		hubs = append(hubs, unmanagedHubs...)
	}

	return hubs, nil
}

// discoverUnmanagedHubs finds hubs that were manually added via kubeconfig secrets
func (r *RHACMClient) discoverUnmanagedHubs(ctx context.Context, existingHubs map[string]bool) ([]models.ManagedHub, error) {
	var unmanagedHubs []models.ManagedHub

	// List all namespaces
	namespaces, err := r.kubeClient.ClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Check each namespace for admin-kubeconfig secret
	for _, ns := range namespaces.Items {
		nsName := ns.Name

		// Skip system namespaces and already discovered hubs
		if existingHubs[nsName] || nsName == "default" || nsName == "kube-system" || nsName == "openshift" || len(nsName) > 20 && nsName[:10] == "openshift-" {
			continue
		}

		// Look for {namespace}-admin-kubeconfig secret
		secretName := nsName + "-admin-kubeconfig"
		secret, err := r.kubeClient.ClientSet.CoreV1().Secrets(nsName).Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			continue // Secret doesn't exist, not a hub
		}

		// Check if created by rhacm-monitor
		if secret.Labels == nil || secret.Labels["created-by"] != "rhacm-monitor" {
			continue
		}

		// This is a manually added hub - try to connect and get basic info
		hub := models.ManagedHub{
			Name:      nsName,
			Namespace: nsName,
			Status:    "External",
			Version:   "Unknown",
			Labels: map[string]string{
				"type": "unmanaged",
			},
			Annotations: map[string]string{
				"source": "manual",
			},
			ClusterInfo: models.ClusterInfo{
				Platform: "External",
			},
			NodesInfo:       []models.NodeInfo{},
			PoliciesInfo:    []models.PolicyInfo{},
			ManagedClusters: nil,
		}

		// Try to connect and get complete information like managed hubs
		hubClient, err := NewHubClientFromSecret(ctx, r.kubeClient, nsName)
		if err == nil {
			// Get cluster information from ClusterVersion resource
			cvGVR := schema.GroupVersionResource{
				Group:    "config.openshift.io",
				Version:  "v1",
				Resource: "clusterversions",
			}
			cvList, err := hubClient.kubeClient.DynamicClient.Resource(cvGVR).List(ctx, metav1.ListOptions{})
			if err == nil && len(cvList.Items) > 0 {
				// Get OpenShift version
				if status, found, _ := unstructured.NestedMap(cvList.Items[0].Object, "status"); found {
					if desired, found, _ := unstructured.NestedMap(status, "desired"); found {
						if version, found, _ := unstructured.NestedString(desired, "version"); found {
							hub.ClusterInfo.OpenshiftVersion = version
						}
					}
				}
				// Get cluster ID
				if spec, found, _ := unstructured.NestedMap(cvList.Items[0].Object, "spec"); found {
					if clusterID, found, _ := unstructured.NestedString(spec, "clusterID"); found {
						hub.ClusterInfo.ClusterID = clusterID
					}
				}
			}

			// Get console URL from console route
			routeGVR := schema.GroupVersionResource{
				Group:    "route.openshift.io",
				Version:  "v1",
				Resource: "routes",
			}
			route, err := hubClient.kubeClient.DynamicClient.Resource(routeGVR).Namespace("openshift-console").Get(ctx, "console", metav1.GetOptions{})
			if err == nil {
				if spec, found, _ := unstructured.NestedMap(route.Object, "spec"); found {
					if host, found, _ := unstructured.NestedString(spec, "host"); found {
						hub.ClusterInfo.ConsoleURL = "https://" + host
					}
				}
			}

			// Get GitOps console URL from openshift-gitops namespace
			gitopsRoutes, err := hubClient.kubeClient.DynamicClient.Resource(routeGVR).Namespace("openshift-gitops").List(ctx, metav1.ListOptions{})
			if err == nil && len(gitopsRoutes.Items) > 0 {
				// Look for server route (usually first route or one with 'server' in name)
				for _, route := range gitopsRoutes.Items {
					if spec, found, _ := unstructured.NestedMap(route.Object, "spec"); found {
						if host, found, _ := unstructured.NestedString(spec, "host"); found {
							hub.ClusterInfo.GitOpsURL = "https://" + host
							break // Use first route found
						}
					}
				}
			}

			// Get nodes and version info
			nodes, err := hubClient.kubeClient.GetNodes(ctx)
			if err == nil && len(nodes.Items) > 0 {
				hub.Status = "Connected"
				hub.Version = nodes.Items[0].Status.NodeInfo.KubeletVersion
				hub.ClusterInfo.KubernetesVersion = nodes.Items[0].Status.NodeInfo.KubeletVersion

				// Convert all nodes
				for i := range nodes.Items {
					nodeInfo := ConvertNodeToNodeInfo(&nodes.Items[i])
					if nodeInfo.Annotations == nil {
						nodeInfo.Annotations = make(map[string]string)
					}
					nodeInfo.Annotations["data-source"] = "Node"
					hub.NodesInfo = append(hub.NodesInfo, nodeInfo)
				}
			}

			// Get spoke clusters
			spokes, err := r.getSpokesClustersFromHub(ctx, nsName)
			if err == nil {
				hub.ManagedClusters = spokes
			}

			// Get policies
			policies, err := hubClient.kubeClient.GetPoliciesForNamespace(ctx, nsName)
			if err == nil {
				hub.PoliciesInfo = policies
			}
		}

		unmanagedHubs = append(unmanagedHubs, hub)
	}

	return unmanagedHubs, nil
}

// GetManagedHub returns a specific managed hub (either from ManagedCluster or kubeconfig secret)
func (r *RHACMClient) GetManagedHub(ctx context.Context, name string) (*models.ManagedHub, error) {
	// Try to get from ManagedCluster first
	cluster, err := r.kubeClient.GetManagedCluster(ctx, name)
	if err == nil && isHub(*cluster) {
		return r.convertToManagedHub(ctx, cluster)
	}

	// Not found as ManagedCluster, check if it's a manually added hub
	secretName := name + "-admin-kubeconfig"
	secret, err := r.kubeClient.ClientSet.CoreV1().Secrets(name).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("hub not found as ManagedCluster or manual hub: %w", err)
	}

	// Verify it was created by rhacm-monitor
	if secret.Labels == nil || secret.Labels["created-by"] != "rhacm-monitor" {
		return nil, fmt.Errorf("hub %s not found", name)
	}

	// This is a manually added hub
	hub := &models.ManagedHub{
		Name:      name,
		Namespace: name,
		Status:    "External",
		Version:   "Unknown",
		Labels: map[string]string{
			"type": "unmanaged",
		},
		Annotations: map[string]string{
			"source": "manual",
		},
		ClusterInfo: models.ClusterInfo{
			Platform: "External",
		},
		NodesInfo:       []models.NodeInfo{},
		PoliciesInfo:    []models.PolicyInfo{},
		ManagedClusters: nil,
	}

	// Try to connect and get basic info
	hubClient, err := NewHubClientFromSecret(ctx, r.kubeClient, name)
	if err == nil {
		// Get nodes
		nodes, err := hubClient.kubeClient.GetNodes(ctx)
		if err == nil && len(nodes.Items) > 0 {
			hub.Status = "Connected"
			hub.Version = nodes.Items[0].Status.NodeInfo.KubeletVersion

			// Try to get OpenShift version and cluster info from ClusterVersion resource
			cvGVR := schema.GroupVersionResource{
				Group:    "config.openshift.io",
				Version:  "v1",
				Resource: "clusterversions",
			}
			cvList, err := hubClient.kubeClient.DynamicClient.Resource(cvGVR).List(ctx, metav1.ListOptions{})
			if err == nil && len(cvList.Items) > 0 {
				// Get version from status.desired.version
				if status, found, _ := unstructured.NestedMap(cvList.Items[0].Object, "status"); found {
					if desired, found, _ := unstructured.NestedMap(status, "desired"); found {
						if version, found, _ := unstructured.NestedString(desired, "version"); found {
							hub.ClusterInfo.OpenshiftVersion = version
						}
					}
				}
				// Get cluster ID from spec.clusterID
				if spec, found, _ := unstructured.NestedMap(cvList.Items[0].Object, "spec"); found {
					if clusterID, found, _ := unstructured.NestedString(spec, "clusterID"); found {
						hub.ClusterInfo.ClusterID = clusterID
					}
				}
			}
			
			// Get console and GitOps URLs from routes
			routeGVR := schema.GroupVersionResource{
				Group:    "route.openshift.io",
				Version:  "v1",
				Resource: "routes",
			}
			
			// Get console URL
			consoleRoute, err := hubClient.kubeClient.DynamicClient.Resource(routeGVR).Namespace("openshift-console").Get(ctx, "console", metav1.GetOptions{})
			if err == nil {
				if spec, found, _ := unstructured.NestedMap(consoleRoute.Object, "spec"); found {
					if host, found, _ := unstructured.NestedString(spec, "host"); found {
						hub.ClusterInfo.ConsoleURL = "https://" + host
					}
				}
			}
			
			// Get GitOps console URL
			gitopsRoutes, err := hubClient.kubeClient.DynamicClient.Resource(routeGVR).Namespace("openshift-gitops").List(ctx, metav1.ListOptions{})
			if err == nil && len(gitopsRoutes.Items) > 0 {
				for _, route := range gitopsRoutes.Items {
					if spec, found, _ := unstructured.NestedMap(route.Object, "spec"); found {
						if host, found, _ := unstructured.NestedString(spec, "host"); found {
							hub.ClusterInfo.GitOpsURL = "https://" + host
							break
						}
					}
				}
			}

			// Convert nodes
			for i := range nodes.Items {
				nodeInfo := ConvertNodeToNodeInfo(&nodes.Items[i])
				if nodeInfo.Annotations == nil {
					nodeInfo.Annotations = make(map[string]string)
				}
				nodeInfo.Annotations["data-source"] = "Node"
				hub.NodesInfo = append(hub.NodesInfo, nodeInfo)
			}
		}

		// Try to get policies
		policies, err := hubClient.kubeClient.GetPoliciesForNamespace(ctx, name)
		if err == nil {
			hub.PoliciesInfo = policies
		}

		// Try to get spoke clusters
		spokes, err := r.getSpokesClustersFromHub(ctx, name)
		if err == nil {
			hub.ManagedClusters = spokes
		}
	}

	return hub, nil
}

// GetManagedClustersForHub returns all managed clusters for a specific hub
func (r *RHACMClient) GetManagedClustersForHub(ctx context.Context, hubName string) ([]models.ManagedCluster, error) {
	// Try to connect to the hub and get its spoke clusters
	spokes, err := r.getSpokesClustersFromHub(ctx, hubName)
	if err != nil {
		return nil, fmt.Errorf("failed to get spoke clusters from hub %s: %w", hubName, err)
	}

	return spokes, nil
}

// getSpokesClustersFromHub connects to a managed hub and retrieves its spoke clusters
func (r *RHACMClient) getSpokesClustersFromHub(ctx context.Context, hubName string) ([]models.ManagedCluster, error) {
	// Create a client for the hub using its kubeconfig secret
	hubClient, err := NewHubClientFromSecret(ctx, r.kubeClient, hubName)
	if err != nil {
		return nil, fmt.Errorf("failed to create hub client: %w", err)
	}

	// Get all managed clusters from the hub
	managedClusters, err := hubClient.GetManagedClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get managed clusters from hub: %w", err)
	}

	// Convert to our model (excluding hub itself if it appears)
	var spokes []models.ManagedCluster
	for i := range managedClusters.Items {
		cluster := &managedClusters.Items[i]

		// Skip if this is a hub cluster or local-cluster (hub itself)
		if isHub(*cluster) || cluster.Name == "local-cluster" {
			continue
		}

		mc, err := r.convertToManagedCluster(ctx, cluster, hubName)
		if err != nil {
			// Log error but continue
			continue
		}

		// Fetch policies for this spoke from the hub in the spoke's namespace
		spokePolicies, err := hubClient.kubeClient.GetPoliciesForNamespace(ctx, cluster.Name)
		if err != nil {
			fmt.Printf("Warning: Could not fetch policies for spoke %s: %v\n", cluster.Name, err)
			spokePolicies = []models.PolicyInfo{}
		}
		mc.PoliciesInfo = spokePolicies

		// Fetch BareMetalHost nodes for this spoke from the hub in the spoke's namespace
		spokeNodes, err := hubClient.kubeClient.GetBareMetalHostsForNamespace(ctx, cluster.Name)
		if err != nil {
			fmt.Printf("Warning: Could not fetch nodes for spoke %s: %v\n", cluster.Name, err)
			spokeNodes = []models.NodeInfo{}
		}
		mc.NodesInfo = spokeNodes

		spokes = append(spokes, *mc)
	}

	return spokes, nil
}

// convertToManagedHub converts a ManagedCluster to a ManagedHub model (called by GetManagedHub for single hub query)
func (r *RHACMClient) convertToManagedHub(ctx context.Context, cluster *clusterv1.ManagedCluster) (*models.ManagedHub, error) {
	// Fetch spokes for this hub using kubeconfig
	spokes, err := r.getSpokesClustersFromHub(ctx, cluster.Name)
	if err != nil {
		// Log error but continue
		fmt.Printf("Warning: Could not fetch spokes from hub %s: %v\n", cluster.Name, err)
		spokes = []models.ManagedCluster{}
	}

	// Fetch policies for this hub from its namespace on the global hub
	hubPolicies, err := r.kubeClient.GetPoliciesForNamespace(ctx, cluster.Name)
	if err != nil {
		fmt.Printf("Warning: Could not fetch policies for hub %s: %v\n", cluster.Name, err)
		hubPolicies = []models.PolicyInfo{}
	}

	// Fetch both K8s Nodes and BareMetalHost for this hub
	var hubNodes []models.NodeInfo

	// Try to connect to hub and get actual K8s nodes
	hubClient, err := NewHubClientFromSecret(ctx, r.kubeClient, cluster.Name)
	if err == nil {
		// Get K8s Node resources from the hub
		k8sNodes, err := hubClient.kubeClient.GetNodes(ctx)
		if err == nil {
			for i := range k8sNodes.Items {
				nodeInfo := ConvertNodeToNodeInfo(&k8sNodes.Items[i])
				if nodeInfo.Annotations == nil {
					nodeInfo.Annotations = make(map[string]string)
				}
				nodeInfo.Annotations["data-source"] = "Node"
				hubNodes = append(hubNodes, nodeInfo)
			}
		}
	}

	// Also get BareMetalHost resources from hub namespace on global hub
	bmhNodes, err := r.kubeClient.GetBareMetalHostsForNamespace(ctx, cluster.Name)
	if err == nil {
		hubNodes = append(hubNodes, bmhNodes...)
	}

	hub := &models.ManagedHub{
		Name:            cluster.Name,
		Namespace:       cluster.Name, // In ACM, namespace equals cluster name
		Status:          getClusterStatus(cluster),
		Version:         getClusterVersion(cluster),
		Conditions:      convertConditions(cluster.Status.Conditions),
		ClusterInfo:     extractClusterInfo(cluster),
		NodesInfo:       hubNodes,
		PoliciesInfo:    hubPolicies,
		ManagedClusters: spokes,
		Labels:          cluster.Labels,
		Annotations:     cluster.Annotations,
		CreatedAt:       cluster.CreationTimestamp.Time,
	}

	return hub, nil
}

// convertToManagedCluster converts a ManagedCluster to a ManagedCluster model
func (r *RHACMClient) convertToManagedCluster(ctx context.Context, cluster *clusterv1.ManagedCluster, hubName string) (*models.ManagedCluster, error) {
	mc := &models.ManagedCluster{
		Name:         cluster.Name,
		Namespace:    cluster.Name,
		Status:       getClusterStatus(cluster),
		Version:      getClusterVersion(cluster),
		Conditions:   convertConditions(cluster.Status.Conditions),
		ClusterInfo:  extractClusterInfo(cluster),
		NodesInfo:    []models.NodeInfo{},   // Would need to fetch from cluster
		PoliciesInfo: []models.PolicyInfo{}, // Would need to fetch policies
		Labels:       cluster.Labels,
		Annotations:  cluster.Annotations,
		HubName:      hubName,
		CreatedAt:    cluster.CreationTimestamp.Time,
	}

	return mc, nil
}

// Helper functions

func isHub(cluster clusterv1.ManagedCluster) bool {
	// Check if cluster has hub-related labels
	if cluster.Labels != nil {
		if val, ok := cluster.Labels["cluster.open-cluster-management.io/clusterset"]; ok && val == "global-hub" {
			return true
		}
		if val, ok := cluster.Labels["feature.open-cluster-management.io/addon-multicluster-hub"]; ok && val == "available" {
			return true
		}
		if val, ok := cluster.Labels["vendor"]; ok && val == "OpenShift" {
			// Additional check for hub-specific claims
			if _, hasHub := cluster.Labels["hub"]; hasHub {
				return true
			}
		}
	}
	return false
}

func belongsToHub(cluster clusterv1.ManagedCluster, hubName string) bool {
	if cluster.Labels != nil {
		if val, ok := cluster.Labels["managed-by"]; ok && val == hubName {
			return true
		}
		if val, ok := cluster.Labels["cluster.open-cluster-management.io/clusterset"]; ok && val == hubName {
			return true
		}
	}
	return false
}

func getClusterStatus(cluster *clusterv1.ManagedCluster) string {
	for _, condition := range cluster.Status.Conditions {
		if condition.Type == clusterv1.ManagedClusterConditionAvailable {
			if condition.Status == metav1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}

func getClusterVersion(cluster *clusterv1.ManagedCluster) string {
	if cluster.Status.Version.Kubernetes != "" {
		return cluster.Status.Version.Kubernetes
	}
	// Try to get from cluster claims
	for _, claim := range cluster.Status.ClusterClaims {
		if claim.Name == "version.openshift.io" {
			return claim.Value
		}
	}
	return "Unknown"
}

func convertConditions(conditions []metav1.Condition) []models.Condition {
	result := make([]models.Condition, len(conditions))
	for i, c := range conditions {
		result[i] = models.Condition{
			Type:               c.Type,
			Status:             string(c.Status),
			LastTransitionTime: c.LastTransitionTime.Time,
			Reason:             c.Reason,
			Message:            c.Message,
		}
	}
	return result
}

func extractClusterInfo(cluster *clusterv1.ManagedCluster) models.ClusterInfo {
	info := models.ClusterInfo{
		ClusterID:         string(cluster.UID),
		KubernetesVersion: cluster.Status.Version.Kubernetes,
		CreatedAt:         cluster.CreationTimestamp.Time,
	}

	// Extract configuration version from labels
	if cluster.Labels != nil {
		if configVersion, ok := cluster.Labels["configuration-version"]; ok {
			info.Region = configVersion // Reuse Region field for configuration version
		}
	}

	// Extract information from cluster claims
	for _, claim := range cluster.Status.ClusterClaims {
		switch claim.Name {
		case "platform.open-cluster-management.io":
			info.Platform = claim.Value
		case "region.open-cluster-management.io":
			// Only set if not already set from configuration-version label
			if info.Region == "" {
				info.Region = claim.Value
			}
		case "version.openshift.io":
			info.OpenshiftVersion = claim.Value
		case "consoleurl.cluster.open-cluster-management.io":
			info.ConsoleURL = claim.Value
		case "kubeversion.open-cluster-management.io":
			info.KubernetesVersion = claim.Value
		}
	}

	return info
}

// ConvertNodeToNodeInfo converts a Kubernetes Node to NodeInfo model
func ConvertNodeToNodeInfo(node *corev1.Node) models.NodeInfo {
	nodeInfo := models.NodeInfo{
		Name:        node.Name,
		Status:      getNodeStatus(node),
		Role:        getNodeRole(node),
		Labels:      node.Labels,
		Annotations: node.Annotations,
		CreatedAt:   node.CreationTimestamp.Time,
		Conditions:  convertNodeConditions(node.Status.Conditions),
		Capacity: models.ResourceList{
			CPU:              node.Status.Capacity.Cpu().String(),
			Memory:           node.Status.Capacity.Memory().String(),
			Storage:          node.Status.Capacity.Storage().String(),
			EphemeralStorage: node.Status.Capacity.StorageEphemeral().String(),
			Pods:             node.Status.Capacity.Pods().String(),
		},
		Allocatable: models.ResourceList{
			CPU:              node.Status.Allocatable.Cpu().String(),
			Memory:           node.Status.Allocatable.Memory().String(),
			Storage:          node.Status.Allocatable.Storage().String(),
			EphemeralStorage: node.Status.Allocatable.StorageEphemeral().String(),
			Pods:             node.Status.Allocatable.Pods().String(),
		},
	}

	// Extract node info
	nodeInfo.KernelVersion = node.Status.NodeInfo.KernelVersion
	nodeInfo.OSImage = node.Status.NodeInfo.OSImage
	nodeInfo.ContainerRuntime = node.Status.NodeInfo.ContainerRuntimeVersion
	nodeInfo.KubeletVersion = node.Status.NodeInfo.KubeletVersion

	// Extract IPs
	for _, addr := range node.Status.Addresses {
		switch addr.Type {
		case corev1.NodeInternalIP:
			nodeInfo.InternalIP = addr.Address
		case corev1.NodeExternalIP:
			nodeInfo.ExternalIP = addr.Address
		}
	}

	return nodeInfo
}

func getNodeStatus(node *corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}

func getNodeRole(node *corev1.Node) string {
	if node.Labels != nil {
		if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
			return "master"
		}
		if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
			return "control-plane"
		}
		if _, ok := node.Labels["node-role.kubernetes.io/worker"]; ok {
			return "worker"
		}
	}
	return "worker"
}

func convertNodeConditions(conditions []corev1.NodeCondition) []models.Condition {
	result := make([]models.Condition, len(conditions))
	for i, c := range conditions {
		result[i] = models.Condition{
			Type:               string(c.Type),
			Status:             string(c.Status),
			LastTransitionTime: c.LastTransitionTime.Time,
			Reason:             c.Reason,
			Message:            c.Message,
		}
	}
	return result
}
