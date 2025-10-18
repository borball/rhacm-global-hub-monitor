package client

import (
	"context"

	"github.com/rhacm-global-hub-monitor/backend/pkg/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// GetPoliciesForNamespace fetches policies from a specific namespace
func (k *KubeClient) GetPoliciesForNamespace(ctx context.Context, namespace string) ([]models.PolicyInfo, error) {
	// Query Policy resources using dynamic client
	policyList, err := k.DynamicClient.Resource(PolicyGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		// Policies may not be available, return empty list
		return []models.PolicyInfo{}, nil
	}

	var policies []models.PolicyInfo
	for _, item := range policyList.Items {
		policy, err := convertUnstructuredToPolicy(&item)
		if err != nil {
			// Skip policies we can't parse
			continue
		}
		policies = append(policies, policy)
	}

	return policies, nil
}

// convertUnstructuredToPolicy converts an unstructured policy to PolicyInfo
func convertUnstructuredToPolicy(obj *unstructured.Unstructured) (models.PolicyInfo, error) {
	policy := models.PolicyInfo{
		Name:        obj.GetName(),
		Namespace:   obj.GetNamespace(),
		Labels:      obj.GetLabels(),
		Annotations: obj.GetAnnotations(),
		CreatedAt:   obj.GetCreationTimestamp().Time,
	}

	// Extract spec fields
	spec, found, err := unstructured.NestedMap(obj.Object, "spec")
	if err == nil && found {
		if remediationAction, ok := spec["remediationAction"].(string); ok {
			policy.RemediationAction = remediationAction
		}
		if disabled, ok := spec["disabled"].(bool); ok {
			policy.Disabled = disabled
		}
	}

	// Extract status fields
	status, found, err := unstructured.NestedMap(obj.Object, "status")
	if err == nil && found {
		if compliant, ok := status["compliant"].(string); ok {
			policy.ComplianceState = compliant
		}
	}

	// Extract annotations for additional info
	if annotations := obj.GetAnnotations(); annotations != nil {
		// Try to get severity from annotations
		if severity, ok := annotations["policy.open-cluster-management.io/severity"]; ok {
			policy.Severity = severity
		} else {
			policy.Severity = "medium" // default
		}

		// Try to get categories
		if categories, ok := annotations["policy.open-cluster-management.io/categories"]; ok {
			policy.Categories = []string{categories}
		}

		// Try to get standards
		if standards, ok := annotations["policy.open-cluster-management.io/standards"]; ok {
			policy.Standards = []string{standards}
		}

		// Try to get controls
		if controls, ok := annotations["policy.open-cluster-management.io/controls"]; ok {
			policy.Controls = []string{controls}
		}
	}

	// Count violations (if status has details)
	if status, found, err := unstructured.NestedMap(obj.Object, "status"); err == nil && found {
		if details, ok := status["details"].([]interface{}); ok {
			policy.Violations = len(details)
		}
	}

	return policy, nil
}
