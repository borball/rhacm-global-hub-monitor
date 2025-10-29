# Known Issues - v3.0.0

## Unmanaged Hub Authentication

**Issue:** Unmanaged hubs (added via UI) connect as "system:anonymous"

**Symptoms:**
```
Warning: Could not fetch console route for dev: routes.route.openshift.io "console" is forbidden: 
User "system:anonymous" cannot get resource "routes" in API group "route.openshift.io"
```

**Affected Operations:**
- Fetching Console URL
- Fetching GitOps URL  
- Fetching nodes
- Fetching routes

**Root Cause:**
The kubeconfig authentication from manually added hubs isn't being properly applied when creating the REST client.

**Workaround:**
Use managed hubs (auto-discovered from ManagedCluster resources) instead of manually adding hubs. Managed hubs authenticate correctly.

**Status:** Under investigation

**Impact:**
- Managed hubs: ✅ Work perfectly
- Unmanaged hubs: ⚠️ Limited functionality (basic info only)

---

## Cache TTL

**Issue:** Cache configured for 30 minutes but appears to expire every ~2 minutes

**Evidence:**
```
18:00:04 - 10.665s (cache miss)
18:02:22 - 10.535s (cache miss) ← Only 2min 18sec later
```

**Investigation:**
- Cache logging added for debugging
- May be related to pod restarts or multiple clients
- Session affinity configured to help consistency

**Impact:**
- Performance benefit still exists (3-4ms cached vs 10s uncached)
- Just not lasting full 30 minutes as expected

**Status:** Under investigation with logging

---

## Recommendations

**For Production Use:**
1. Use managed hubs (ManagedCluster resources)
2. Enable RBAC with proper service accounts
3. Monitor cache behavior with logs
4. Use refresh buttons for on-demand updates

**Managed Hubs Work Perfectly:**
- Full functionality
- All features operational
- Console/GitOps URLs
- Operators monitoring
- Dark/light mode
- Performance caching
