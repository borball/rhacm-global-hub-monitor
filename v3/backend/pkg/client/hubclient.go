package client

import (
	"context"
	"fmt"

	"github.com/rhacm-global-hub-monitor/backend/pkg/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

// HubClient represents a client for a specific managed hub
type HubClient struct {
	kubeClient *KubeClient
	hubName    string
}

// GetKubeClient returns the underlying KubeClient
func (h *HubClient) GetKubeClient() *KubeClient {
	return h.kubeClient
}

// NewHubClientFromSecret creates a new hub client from a kubeconfig secret
func NewHubClientFromSecret(ctx context.Context, globalHubClient *KubeClient, hubName string) (*HubClient, error) {
	// Get the kubeconfig secret from the hub's namespace
	secretName := fmt.Sprintf("%s-admin-kubeconfig", hubName)
	secret, err := globalHubClient.ClientSet.CoreV1().Secrets(hubName).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig secret %s/%s: %w", hubName, secretName, err)
	}

	// Get the kubeconfig data
	kubeconfigData, ok := secret.Data["kubeconfig"]
	if !ok {
		return nil, fmt.Errorf("kubeconfig not found in secret %s/%s", hubName, secretName)
	}

	// Build config from kubeconfig data using proper clientcmd
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %w", err)
	}
	
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}

	// Create Kubernetes client with the hub's kubeconfig
	kubeClient, err := NewKubeClientFromConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for hub %s: %w", hubName, err)
	}

	return &HubClient{
		kubeClient: kubeClient,
		hubName:    hubName,
	}, nil
}

// NewSpokeClientFromKubeconfig creates a spoke client from kubeconfig bytes
func NewSpokeClientFromKubeconfig(kubeconfigData []byte, spokeName string) (*HubClient, error) {
	// Build config from kubeconfig data using proper clientcmd
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %w", err)
	}
	
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}

	// Create Kubernetes client
	kubeClient, err := NewKubeClientFromConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for spoke %s: %w", spokeName, err)
	}

	return &HubClient{
		kubeClient: kubeClient,
		hubName:    spokeName,
	}, nil
}

// NewKubeClientFromConfig creates a KubeClient from an existing rest.Config
func NewKubeClientFromConfig(config *rest.Config) (*KubeClient, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &KubeClient{
		ClientSet:     clientSet,
		DynamicClient: dynamicClient,
		RestConfig:    config,
	}, nil
}

// GetManagedClusters returns all managed clusters (spokes) on this hub
func (h *HubClient) GetManagedClusters(ctx context.Context) (*clusterv1.ManagedClusterList, error) {
	return h.kubeClient.GetManagedClusters(ctx)
}

// GetNodesForCluster returns node information for a cluster
func (h *HubClient) GetNodesForCluster(ctx context.Context, namespace string) ([]models.NodeInfo, error) {
	var allNodes []models.NodeInfo

	// Get Kubernetes Node resources if this is the cluster itself
	// For spoke clusters, we query nodes in the spoke's namespace which gives us BMH
	nodes, err := h.kubeClient.GetNodes(ctx)
	if err == nil && namespace == "" {
		// This is the hub cluster itself, use actual nodes
		for i := range nodes.Items {
			nodeInfo := ConvertNodeToNodeInfo(&nodes.Items[i])
			if nodeInfo.Annotations == nil {
				nodeInfo.Annotations = make(map[string]string)
			}
			nodeInfo.Annotations["data-source"] = "Node"
			allNodes = append(allNodes, nodeInfo)
		}
	}

	// Also get BareMetalHost resources for hardware details
	if namespace != "" {
		bmhNodes, err := h.kubeClient.GetBareMetalHostsForNamespace(ctx, namespace)
		if err == nil {
			allNodes = append(allNodes, bmhNodes...)
		}
	}

	return allNodes, nil
}
