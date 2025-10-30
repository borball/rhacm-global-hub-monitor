# RHACM Global Hub Monitor v4.0.0 - Roadmap

**Status:** 🚧 In Development  
**Target Release:** Q1 2026

## Overview

v4 focuses on making the monitor work in ANY environment (with or without RHACM) and adding Global Hub visibility.

## Primary Goals

### 1. Non-RHACM Environment Support (HIGH PRIORITY)

**Feature:** Make monitor work on clusters without RHACM installed

**Current State:**
- Monitor assumes RHACM is installed
- Looks for ManagedCluster resources
- Fails gracefully but shows empty state

**v4 Implementation:**

**Detection Logic:**
```go
func isRHACMInstalled(ctx context.Context, client *KubeClient) bool {
    // Check if RHACM CRDs exist
    _, err := client.DiscoveryClient.ServerResourcesForGroupVersion(
        "cluster.open-cluster-management.io/v1")
    return err == nil
}
```

**UI Adaptation:**
```javascript
// In fetchHubs()
if (data.data.rhacmInstalled === false) {
    // Show only "Unmanaged Hubs" section
    // Hide "Managed Hubs" section
    // Show message: "RHACM not detected - using manual hub management mode"
} else {
    // Show both sections as normal
}
```

**API Response:**
```json
{
  "success": true,
  "rhacmInstalled": false,
  "data": [
    {"name": "hub1", "source": "manual", ...},
    {"name": "hub2", "source": "manual", ...}
  ]
}
```

**Benefits:**
- ✅ Works on any Kubernetes cluster
- ✅ Doesn't require RHACM installation
- ✅ Still useful for monitoring external hubs
- ✅ Graceful adaptation to environment

### 2. Global Hub Dashboard (HIGH PRIORITY)

**Feature:** Show Global Hub information and topology

**Global Hub Section:**

**Location:** New section at top of page (above hub cards)

**Information Displayed:**

1. **Global Hub Details:**
   - Cluster name
   - OpenShift version
   - Kubernetes version
   - Node count
   - Status

2. **Topology Overview:**
   - Total managed hubs
   - Total spoke clusters (SNO)
   - Total policies
   - Total operators

3. **Managed Hubs Summary:**
   - List of managed hubs with status
   - Spoke count per hub
   - Quick links to hub details

**Implementation:**

**Backend:**
```go
// New endpoint
func (h *HubHandler) GetGlobalHubInfo(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Get global hub cluster info (where monitor is deployed)
    globalHub := models.GlobalHub{
        Name: getClusterName(),
        Version: getKubernetesVersion(),
        OpenshiftVersion: getOpenshiftVersion(),
        NodeCount: getNodeCount(),
        ManagedHubs: len(getManagedHubs()),
        TotalSpokes: getTotalSpokes(),
        TotalPolicies: getTotalPolicies(),
    }
    
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Data: globalHub,
    })
}

// GET /api/global-hub
```

**Frontend:**
```javascript
// New section at top
function renderGlobalHubSection(globalHub) {
    return `
        <div class="global-hub-section">
            <h2>🌐 Global Hub: ${globalHub.name}</h2>
            <div class="global-stats">
                <div class="stat-card">
                    <div class="stat-label">Managed Hubs</div>
                    <div class="stat-value">${globalHub.managedHubs}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Spoke Clusters</div>
                    <div class="stat-value">${globalHub.totalSpokes}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Policies</div>
                    <div class="stat-value">${globalHub.totalPolicies}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Operators</div>
                    <div class="stat-value">${globalHub.totalOperators}</div>
                </div>
            </div>
            <div class="topology-view">
                <h3>Hub Topology</h3>
                <!-- Visual tree of hubs → spokes -->
            </div>
        </div>
    `;
}
```

**Visual Topology:**
```
Global Hub (vhub)
├── acm1 (Connected)
│   ├── sno146 (Ready)
│   ├── sno132 (Ready)
│   └── sno133 (Ready)
├── acm2 (Connected)
│   ├── sno171 (Ready)
│   └── sno180 (Ready)
└── acm3 (External - Unmanaged)
```

## Implementation Plan

### Phase 1: Environment Detection (Week 1)

**Tasks:**
1. Add RHACM detection logic
2. Update API to return rhacmInstalled flag
3. Adapt UI based on environment
4. Test on non-RHACM cluster

**Deliverable:** Monitor works on any cluster

### Phase 2: Global Hub Dashboard (Weeks 2-3)

**Tasks:**
1. Create GlobalHub model
2. Implement /api/global-hub endpoint
3. Design Global Hub UI section
4. Add topology visualization
5. Show aggregated statistics

**Deliverable:** Global Hub overview dashboard

### Phase 3: Testing & Polish (Week 4)

**Tasks:**
1. Test on RHACM environment
2. Test on non-RHACM environment
3. Test topology with multiple hubs
4. Documentation updates
5. Performance testing

**Deliverable:** v4.0.0 release

## Additional v4 Features (If Time Permits)

### 3. Improved Authentication

**Current Issue:** Basic auth not working with client-go

**Investigation:** 
- Test with custom RoundTripper
- Try token generation from credentials
- Research client-go auth mechanisms

**Priority:** Medium (after main features)

### 4. Enhanced Caching

**Current Issue:** Cache expires faster than configured

**Options:**
- Redis integration for shared cache
- Better in-memory cache debugging
- Cache metrics

**Priority:** Low (works acceptably now)

## Success Criteria

### Must Have
- ✅ Works on non-RHACM clusters
- ✅ Global Hub dashboard visible
- ✅ Topology view functional
- ✅ All v3 features maintained

### Should Have
- ✅ Visual topology tree
- ✅ Aggregated statistics
- ✅ Clean UI design

### Nice to Have
- Basic auth working
- Redis cache
- Historical data

## Technical Approach

### Environment Detection

```go
// pkg/client/environment.go
package client

func (k *KubeClient) DetectEnvironment(ctx context.Context) EnvironmentInfo {
    env := EnvironmentInfo{
        IsOpenShift: false,
        HasRHACM: false,
        ClusterName: "",
        Version: "",
    }
    
    // Check for OpenShift
    _, err := k.DynamicClient.Resource(clusterVersionGVR).List(ctx, metav1.ListOptions{})
    env.IsOpenShift = (err == nil)
    
    // Check for RHACM
    _, err = k.DynamicClient.Resource(managedClusterGVR).List(ctx, metav1.ListOptions{})
    env.HasRHACM = (err == nil)
    
    // Get cluster name from infrastructure
    infra, _ := k.DynamicClient.Resource(infrastructureGVR).Get(ctx, "cluster", metav1.GetOptions{})
    if infra != nil {
        env.ClusterName, _, _ = unstructured.NestedString(infra.Object, "status", "infrastructureName")
    }
    
    return env
}
```

### Global Hub Info

```go
// pkg/models/global_hub.go
type GlobalHub struct {
    Name              string
    Version           string
    OpenshiftVersion  string
    NodeCount         int
    ManagedHubCount   int
    SpokeCount        int
    PolicyCount       int
    OperatorCount     int
    RHACMInstalled    bool
    Topology          HubTopology
}

type HubTopology struct {
    Hubs []HubNode
}

type HubNode struct {
    Name   string
    Status string
    Spokes []SpokeNode
}

type SpokeNode struct {
    Name   string
    Status string
}
```

## Timeline

| Week | Milestone |
|------|-----------|
| 1 | Environment detection complete |
| 2 | Global Hub info endpoint |
| 3 | UI with Global Hub section |
| 4 | Testing & documentation |

## Migration from v3

**Backward Compatible:** v4 includes all v3 features

**Upgrade:**
```bash
# Update image
oc set image deployment/rhacm-monitor-backend \
  rhacm-monitor-backend=quay.io/bzhai/rhacm-monitor-backend:v4.0.0 \
  -n rhacm-monitor

# No config changes needed
```

**New Features:**
- Automatically adapts to environment
- Global Hub section appears automatically
- All v3 features continue to work

---

*v4 makes the monitor universal - works anywhere, shows everything.*
