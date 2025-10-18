# Spoke Cluster Discovery Test Results

## Test Date
**October 17, 2025**

## Test Objective
Test if the API, when connected to the **Global Hub (vhub)**, can return all managed spoke clusters for acm1.

## Environment
- **Global Hub**: vhub.outbound.vz.bos2.lab
- **Managed Hubs**: acm1, acm2
- **acm1 Managed Spokes**: 1 cluster (sno146)

## Architecture Understanding

### RHACM Global Hub Architecture

```
┌─────────────────────────────────────────────────────────┐
│            Global Hub (vhub)                            │
│  - Manages hub clusters (acm1, acm2)                    │
│  - Does NOT have spoke clusters registered directly     │
└──────────────┬────────────────────────────┬──────────────┘
               │                            │
               ▼                            ▼
    ┌──────────────────┐         ┌──────────────────┐
    │   Hub: acm1      │         │   Hub: acm2      │
    │  ManagedClusters:│         │  ManagedClusters:│
    │  - sno146        │         │  - (spokes)      │
    └──────────────────┘         └──────────────────┘
```

### Key Finding: Spoke Clusters Are NOT Visible from Global Hub

**Important Discovery**: Spoke clusters (like sno146) are **NOT** registered as `ManagedCluster` resources on the global hub. They are only registered on their managing hub (acm1).

## Test Results

### 1. From Global Hub (vhub)

**Command**: 
```bash
oc get managedclusters
```

**Result**:
```
NAME   HUB ACCEPTED   MANAGED CLUSTER URLS                         JOINED   AVAILABLE   AGE
acm1   true           https://api.acm1.outbound.vz.bos2.lab:6443   True     True        2d
acm2   true           https://api.acm2.outbound.vz.bos2.lab:6443   True     True        123m
```

**Observation**: Only hub clusters (acm1, acm2) are visible. No spoke clusters.

### 2. From Managed Hub (acm1)

**Command**:
```bash
acm1
oc get managedclusters
```

**Result**:
```
NAME     HUB ACCEPTED   MANAGED CLUSTER URLS                           JOINED   AVAILABLE   AGE
sno146   true           https://api.sno146.outbound.vz.bos2.lab:6443   True     True        46h
```

**Spoke Cluster Details**:
- **Name**: sno146
- **Status**: Available (True)
- **Kubernetes Version**: v1.31.8
- **Platform**: Other
- **Type**: Single Node OpenShift (SNO)

### 3. API Test from Global Hub

**API Endpoint**: `GET /api/hubs/acm1`

**Result**:
```json
{
  "success": true,
  "hub": {
    "name": "acm1",
    "status": "Ready",
    "version": "v1.31.13",
    "consoleURL": "https://console-openshift-console.apps.acm1.outbound.vz.bos2.lab",
    "labels": 17,
    "conditions": 5,
    "managedClusters": 0  // ← Shows 0 because spokes not visible from vhub
  }
}
```

## Current Limitation

### Why Spoke Clusters Are Not Returned

The API currently queries `ManagedCluster` resources from the cluster it's connected to. When connected to vhub:

1. **What we see**: acm1, acm2 (hub clusters)
2. **What we don't see**: sno146 (spoke cluster managed by acm1)
3. **Why**: Spoke clusters are NOT registered as `ManagedCluster` on vhub

### To Access Spoke Clusters

To retrieve spoke clusters from the global hub, the application would need to:

1. **Retrieve Hub Credentials**: Get the kubeconfig/credentials for each managed hub from vhub
   - These are typically stored as secrets in the hub's namespace
   - Secret name pattern: `{hub-name}-admin-kubeconfig` or similar

2. **Connect to Each Hub**: Create a Kubernetes client for each hub using its credentials

3. **Query Each Hub**: Query `ManagedCluster` resources from each hub

4. **Aggregate Results**: Combine all spoke clusters with their hub association

## Implementation Options

### Option 1: Multi-Cluster Client (Recommended)

Add functionality to fetch hub credentials and create clients for each hub:

```go
// Pseudocode
func (r *RHACMClient) GetSpokeClustersFo rHub(ctx context.Context, hubName string) ([]models.ManagedCluster, error) {
    // 1. Get hub kubeconfig secret from global hub
    secret, err := r.kubeClient.GetSecret(hubName, "{hubName}-admin-kubeconfig")
    
    // 2. Create client for the hub
    hubClient, err := client.NewClientFromKubeconfig(secret.Data["kubeconfig"])
    
    // 3. Query managed clusters from the hub
    spokeClusters, err := hubClient.GetManagedClusters(ctx)
    
    // 4. Convert and return
    return convertToSpokeClusters(spokeClusters, hubName)
}
```

**Pros**:
- Can retrieve spoke cluster information
- Follows RHACM architecture
- Works with existing permissions

**Cons**:
- More complex implementation
- Requires access to hub credentials
- Additional API calls per hub

### Option 2: Global Hub API

Use RHACM Global Hub APIs if available:

```go
// Check if global hub APIs expose spoke cluster information
// Example: multicluster.x-k8s.io resources
```

**Pros**:
- May be more efficient
- Designed for this use case

**Cons**:
- Requires specific RHACM Global Hub features
- May not be available in all versions

### Option 3: Status Quo with Documentation

Document that spoke clusters must be viewed by connecting to individual hubs:

```bash
# To see spokes for acm1
acm1
oc get managedclusters

# Or use the API connected to acm1
# Start server with KUBECONFIG pointing to acm1
```

**Pros**:
- Simple, already works
- No code changes needed

**Cons**:
- Requires switching contexts
- Not a unified view

## Recommendation

**Implement Option 1** in a future version to provide a complete global view.

**Current Workaround**:

To view spoke clusters for acm1:

```bash
# Switch to acm1
acm1

# Start API server connected to acm1  
cd backend
ENABLE_AUTH=false PORT=8080 ./bin/server &

# Query API (but it won't show "hubs", it will show spokes as managed clusters)
curl http://localhost:8080/api/hubs

# Better: Connect to acm1 directly
oc get managedclusters --context acm1
```

## Summary

| View From | What API Returns | Spoke Clusters Visible? |
|-----------|------------------|------------------------|
| vhub (Global Hub) | acm1, acm2 as managed hubs | ❌ No (not registered on vhub) |
| acm1 (Hub) | sno146 as managed cluster | ✅ Yes (registered on acm1) |
| acm2 (Hub) | Its spokes | ✅ Yes (registered on acm2) |

## Test Conclusion

✅ **API Working as Designed**: The API correctly returns all `ManagedCluster` resources visible from the connected cluster.

⚠️ **Architecture Limitation**: Spoke clusters are not visible from the global hub without additional implementation to connect to each hub.

📝 **Feature Request**: Add multi-hub client support to retrieve spoke clusters from managed hubs.

## Next Steps

### Short Term (Current Capabilities)
- ✅ Use API to view hub clusters from global hub
- ✅ Switch to individual hubs to view their spokes
- ✅ Document the architecture for users

### Medium Term (Enhancement)
- 🔄 Implement hub credential retrieval
- 🔄 Add multi-cluster client support
- 🔄 Aggregate spoke cluster information

### Long Term (Full Solution)
- 🔄 Real-time spoke cluster monitoring
- 🔄 Unified dashboard showing all levels
- 🔄 Metrics and status aggregation across hubs

## Appendix: Test Commands

```bash
# View from global hub
vhub
oc get managedclusters  # Shows: acm1, acm2

# View from acm1
acm1
oc get managedclusters  # Shows: sno146

# Check spoke details
oc get managedcluster sno146 -o yaml

# Test API from global hub
curl http://localhost:8080/api/hubs | jq .

# Test API from acm1
acm1
cd backend && ENABLE_AUTH=false PORT=8080 ./bin/server &
curl http://localhost:8080/api/hubs | jq .
```

## Documentation References

- RHACM Architecture: Multi-level hub management
- Global Hub Configuration: Hub-of-hubs pattern
- ManagedCluster Resources: Scoped to local hub

