# Issues from v3 to Fix in v4

## Critical Issues

### Issue #1: Unmanaged Hub Authentication

**Priority:** 🔴 Critical  
**Status:** Needs Fix  
**Assigned to:** v4.0.0

**Problem:**
Unmanaged hubs (manually added via UI) connect as `system:anonymous` instead of using kubeconfig credentials.

**Symptoms:**
```
Warning: User "system:anonymous" cannot get resource "routes"
Warning: User "system:anonymous" cannot list resource "nodes"  
```

**Impact:**
- Cannot fetch Console/GitOps URLs
- Cannot get cluster version
- Cannot list nodes
- Status shows "Unknown"
- Limited functionality

**Root Cause Analysis:**

Current code in `hubclient.go`:
```go
// This doesn't properly handle all auth methods
config, err := clientcmd.Load(kubeconfigData)
config, err := clientcmd.NewDefaultClientConfig(*kubeConfig, ...).ClientConfig()
```

**Why it fails:**
1. `clientcmd` functions expect file paths or use default rules
2. Auth context might not be selected correctly
3. Client certificate/token data not being applied to REST config
4. Exec plugins not being executed

**Proposed Fix for v4:**

```go
func NewHubClientFromKubeconfigData(kubeconfigData []byte, hubName string) (*HubClient, error) {
    // Parse kubeconfig
    config, err := clientcmd.Load(kubeconfigData)
    if err != nil {
        return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
    }
    
    // Determine which context to use
    contextName := config.CurrentContext
    if contextName == "" {
        // Use first available context
        for name := range config.Contexts {
            contextName = name
            break
        }
    }
    
    context := config.Contexts[contextName]
    cluster := config.Clusters[context.Cluster]
    authInfo := config.AuthInfos[context.AuthInfo]
    
    // Build REST config with explicit auth
    restConfig := &rest.Config{
        Host: cluster.Server,
        TLSClientConfig: rest.TLSClientConfig{
            Insecure:   cluster.InsecureSkipTLSVerify,
            ServerName: cluster.TLSServerName,
        },
    }
    
    // Handle CA certificate
    if len(cluster.CertificateAuthorityData) > 0 {
        restConfig.TLSClientConfig.CAData = cluster.CertificateAuthorityData
    } else if cluster.CertificateAuthority != "" {
        restConfig.TLSClientConfig.CAFile = cluster.CertificateAuthority
    }
    
    // Handle client authentication
    if len(authInfo.ClientCertificateData) > 0 && len(authInfo.ClientKeyData) > 0 {
        // Client certificate auth
        restConfig.TLSClientConfig.CertData = authInfo.ClientCertificateData
        restConfig.TLSClientConfig.KeyData = authInfo.ClientKeyData
    } else if authInfo.Token != "" {
        // Bearer token auth
        restConfig.BearerToken = authInfo.Token
    } else if authInfo.TokenFile != "" {
        restConfig.BearerTokenFile = authInfo.TokenFile
    } else if authInfo.Username != "" {
        // Basic auth
        restConfig.Username = authInfo.Username
        restConfig.Password = authInfo.Password
    } else if authInfo.Exec != nil {
        // Exec plugin - need to execute and get token
        token, err := executeAuthPlugin(authInfo.Exec)
        if err != nil {
            return nil, fmt.Errorf("exec plugin failed: %w", err)
        }
        restConfig.BearerToken = token
    }
    
    // Create client
    kubeClient, err := NewKubeClientFromConfig(restConfig)
    if err != nil {
        return nil, err
    }
    
    return &HubClient{kubeClient: kubeClient, hubName: hubName}, nil
}
```

**Testing Plan:**
1. Create test kubeconfigs with different auth methods
2. Test each auth type separately
3. Verify no anonymous authentication
4. Check RBAC works correctly

**Acceptance Criteria:**
- ✅ No "system:anonymous" in logs
- ✅ Routes fetched successfully
- ✅ Nodes listed correctly
- ✅ Console/GitOps URLs displayed
- ✅ Proper status shown

---

## High Priority Issues

### Issue #2: Cache TTL Not Working

**Priority:** 🟠 High  
**Status:** Under Investigation  
**Assigned to:** v4.0.0

**Problem:**
Cache configured for 30 minutes but expires every ~2 minutes.

**Evidence:**
```
Configured: cache.NewCache(30 * time.Minute)
Actual:     ~2 minute expiration based on logs
```

**Possible Causes:**

1. **Cleanup Goroutine Too Aggressive:**
```go
// Current code
func (c *Cache) cleanup() {
    ticker := time.NewTicker(60 * time.Second)
    for range ticker.C {
        // Deletes expired items every 60 seconds
        // Could this be deleting too much?
    }
}
```

2. **Pod Restarts:**
- Kubernetes might be restarting pods
- In-memory cache lost on restart
- No persistence

3. **Multiple Instances:**
- Despite "shared" cache in main.go
- Could handlers be creating new instances?
- Race conditions?

**Investigation Steps:**

1. Add detailed logging to cache operations
2. Monitor pod uptime vs cache misses
3. Check if multiple cache instances exist
4. Profile memory usage

**Proposed Solutions:**

**Short-term (v4.0):**
- Fix in-memory cache with extensive logging
- Identify exact cause
- Patch cleanup logic

**Long-term (v4.1):**
- Migrate to Redis
- Persistent cache across pod restarts
- Shared across all replicas

**Testing:**
```bash
# Test cache TTL
time curl https://monitor/api/hubs  # ~10s
sleep 5
time curl https://monitor/api/hubs  # Should be < 100ms
sleep 120
time curl https://monitor/api/hubs  # Should still be < 100ms
sleep 1800
time curl https://monitor/api/hubs  # After 30 min, should be ~10s again
```

---

## Medium Priority Issues

### Issue #3: Page Not Auto-Refreshing After Add Hub

**Priority:** 🟡 Medium  
**Status:** Fix Committed, Needs Deployment  
**Assigned to:** v4.0.0

**Problem:**
After adding an unmanaged hub via UI, page doesn't automatically refresh to show it.

**Current State:**
- Fix coded: `delete window.cachedHubsData; fetchHubs();`
- Committed to git
- Not deployed (old image running)

**Solution:**
Already implemented, just needs fresh deployment with new code.

**Testing:**
1. Add unmanaged hub via UI
2. Verify page refreshes automatically
3. Verify hub appears immediately

---

### Issue #4: Operators Lazy Loading Indicator

**Priority:** 🟡 Medium  
**Status:** Working but Could Improve  
**Assigned to:** v4.0.0

**Current:**
- Shows "..." in spoke table
- Shows "Loading..." on expansion
- Works but could be better

**Improvement:**
- Animated spinner while loading
- Progress indicator for multiple spokes
- Pre-fetch on hover (speculative)

**Implementation:**
```javascript
// On hover, start fetching
onmouseenter="prefetchSpokeOperators('${spoke.name}')"

function prefetchSpokeOperators(spokeName) {
    if (!operatorsCache[spokeName]) {
        // Start fetch in background
        fetch(`${API_BASE}/hubs/${hubName}/spokes/${spokeName}/operators`)
            .then(r => r.json())
            .then(data => operatorsCache[spokeName] = data.data);
    }
}
```

---

## Low Priority Issues

### Issue #5: No Bulk Operations

**Priority:** 🟢 Low  
**Status:** Future Enhancement  
**Assigned to:** v4.1 or later

**Feature Request:**
- Refresh all hubs at once
- Bulk policy enforcement
- Export all data

**Implementation:**
```javascript
// Refresh all button
async function refreshAllHubs() {
    for (const hub of hubs) {
        await refreshHub(hub.name);
    }
}
```

---

### Issue #6: Limited Search/Filter

**Priority:** 🟢 Low  
**Status:** Future Enhancement  
**Assigned to:** v4.1 or later

**Current:**
- Search operators only
- No global search

**Wanted:**
- Search across all hubs
- Filter by status, version
- Search spoke clusters
- Search policies

---

## Lessons Learned from v3

### What Worked Well
1. ✅ Iterative development
2. ✅ User feedback driven
3. ✅ Simple technology choices
4. ✅ Comprehensive documentation

### What to Improve in v4
1. 🔧 Test authentication methods upfront
2. 🔧 Prototype cache solution earlier
3. 🔧 Set up CI/CD pipeline
4. 🔧 Automated testing

### Technical Debt to Avoid
1. Don't let cache issues linger
2. Fix auth properly first time
3. Add tests for complex features
4. Monitor technical debt continuously

---

## Conclusion

v4 will address v3's technical debt while adding valuable monitoring features. The focus is on:
- **Reliability:** Fix auth and cache issues
- **Observability:** Add metrics and logging
- **Usability:** Maintain simplicity while adding power

**Primary Goal:** Make unmanaged hubs work as well as managed hubs.

---

*v4 planning is complete. Ready to start implementation.*

