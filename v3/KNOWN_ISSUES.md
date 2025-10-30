# Known Issues - v3.0.0

## Unmanaged Hub Authentication

**Issue:** Unmanaged hubs using basic auth connect as "system:anonymous"

**Status:** Under Investigation

**Details:**
- Basic auth credentials are detected and set correctly
- REST config shows Username and Password
- Kubernetes API server rejects basic auth (deprecated in K8s 1.19+)
- Falls back to anonymous access

**Workaround:**
- Use managed hubs (auto-discovered from ManagedCluster resources)
- Or use kubeconfig with token or certificate auth

**Impact:**
- Unmanaged hubs with basic auth: Limited functionality
- Cannot fetch Console/GitOps URLs, nodes, or full cluster info
- Managed hubs: Full functionality (not affected)

**Future Fix:**
v4 will investigate token generation from basic auth credentials or
alternative authentication methods.

## Current Functionality

### ✅ Works Perfectly
- Managed hubs (acm1, acm2, etc.)
- All v3 features
- Dark/Light mode
- Operators monitoring
- Performance caching

### ⚠️ Limitations
- Unmanaged hubs with basic auth have limited data
- Client-go basic auth support varies

## Recommendations

**For Production:**
1. Use managed hubs (auto-discovered)
2. If adding unmanaged hubs, use token or certificate auth in kubeconfig
3. Avoid basic auth in kubeconfigs (deprecated)

**For Development:**
1. Generate service account token for unmanaged clusters
2. Use that token in kubeconfig instead of username/password
3. Full functionality will work

---

*v3 is production-ready for managed hub monitoring.*
