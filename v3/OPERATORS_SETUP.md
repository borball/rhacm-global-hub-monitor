# Operators Feature Setup Guide

## Overview

The RHACM Global Hub Monitor v3 includes comprehensive operator monitoring for both hubs and spoke clusters with intelligent lazy loading for performance.

## Current Status

### Hub Operators ✅ Working
- **Implementation:** Fully functional
- **Data Source:** Direct connection to hub clusters
- **Display:** Full operators tab with search/filter
- **Example:** acm1 shows ~45 unique operators (304 total installations)
- **Grouping:** Multiple namespace installations consolidated
- **Performance:** Included in initial hub load

### Spoke Operators ✅ Working with Lazy Loading
- **Implementation:** Lazy loading for performance
- **Data Source:** Direct connection to spoke clusters via kubeconfig
- **Display:** Loaded on spoke expansion, shows in stat card + table
- **Example:** sno146 shows 7 operators when expanded
- **Performance:** Fetched on-demand only (~2s per spoke)

## Architecture

```
Global Hub (vhub)
  └─ Monitor Backend
      ├─ Hub (acm1)
      │   ├─ Hub Operators: ✅ Fetched via kubeconfig (included in hub load)
      │   └─ Spoke Clusters
      │       └─ sno146
      │           ├─ Operators installed ON sno146
      │           ├─ Kubeconfig stored in sno146 namespace on acm1
      │           └─ Fetched on-demand when spoke expanded
      └─ Hub (production-hub)
          └─ Similar structure
```

## Performance & Lazy Loading

### Why Lazy Loading?

**Problem:**
- Fetching operators for 1000+ spoke clusters takes too long
- Initial page load would be very slow
- Poor user experience

**Solution:**
- Hub operators: Fetched with hub data (acceptable delay)
- Spoke operators: Lazy loaded on expansion (on-demand)
- Scales efficiently to large deployments

### Performance Metrics

**Initial Homepage Load:**
- Hub data: ~10s (includes hub operators)
- Spoke data: Fast (no operators)
- Total: ~10s initial, ~50ms cached

**Spoke Expansion:**
- First expansion: ~2s (fetches operators)
- Operator data cached in DOM
- Subsequent expansions: Instant

**Scalability:**
- Works with 1000+ spoke clusters
- Only fetches operators for expanded spokes
- No performance degradation

## Setup for Spoke Operators

Spoke operators require kubeconfig secrets to connect to spoke clusters:

### Kubeconfig Location

The spoke's kubeconfig must be stored on the **hub cluster** (not the global hub):

```
Hub Cluster (e.g., acm1)
  └─ Namespace: sno146
      └─ Secret: sno146-admin-kubeconfig
          └─ Data: kubeconfig → connects to sno146 cluster
```

### Verification

Check if kubeconfig exists:

```bash
# On the hub cluster (e.g., acm1)
oc get secret -n sno146 sno146-admin-kubeconfig

# Expected output:
# NAME                      TYPE     DATA   AGE
# sno146-admin-kubeconfig   Opaque   2      5d
```

### If Kubeconfig Doesn't Exist

Create the kubeconfig secret:

```bash
# On the hub cluster (e.g., acm1)

# 1. Create namespace if it doesn't exist
oc create namespace sno146

# 2. Create secret with spoke kubeconfig
oc create secret generic sno146-admin-kubeconfig \
  --from-file=kubeconfig=/path/to/sno146/kubeconfig \
  -n sno146

# 3. Label it for identification
oc label secret sno146-admin-kubeconfig \
  created-by=rhacm-monitor \
  -n sno146
```

### Expected Results

With kubeconfig in place:

**Spoke Table:**
```
sno146 | Ready | 4.18.13 | config | ... | 7 | 18/19 | Details
                                          ↑
                                    Shows "..." initially
```

**Spoke Expansion:**
```
Operators: 7

Operator                    Version         Namespace                    Status
────────────────────────────────────────────────────────────────────────────────
OADP Operator              1.4.5           openshift-adp                ✅ Succeeded
Lifecycle Agent            4.18.0          openshift-lifecycle-agent    ✅ Succeeded
Local Storage              4.18.0-...      openshift-local-storage      ✅ Succeeded
... and 4 more
```

**After Expansion (stat card updates):**
```
Operators: 7  (instead of "...")
```

## Lazy Loading API

### Endpoint

```
GET /api/hubs/:name/spokes/:spoke/operators
```

**Example:**
```bash
curl https://monitor/api/hubs/acm1/spokes/sno146/operators

Response:
{
  "success": true,
  "data": [
    {
      "name": "oadp-operator.v1.4.5",
      "displayName": "OADP Operator",
      "version": "1.4.5",
      "namespace": "openshift-adp",
      "phase": "Succeeded",
      "provider": "Red Hat",
      "createdAt": "2025-10-16T13:05:17Z"
    },
    ...
  ]
}
```

### Frontend Integration

**Automatic on Spoke Expansion:**
1. User expands spoke details
2. Frontend checks if operators already loaded
3. If not, calls lazy loading endpoint
4. Updates UI with operators data
5. Data cached for that spoke

**No User Action Required:**
- Transparent lazy loading
- Automatic on first expansion
- Cached after first load

## Troubleshooting

**Spoke shows 0 operators (or "..."):**

1. **Check if kubeconfig secret exists:**
   ```bash
   # On hub cluster (context: acm1)
   oc get secret -n sno146 sno146-admin-kubeconfig
   ```

2. **Check backend logs:**
   ```bash
   oc logs -l component=backend -n rhacm-monitor | grep sno146
   ```
   
   Look for:
   - ✅ "Successfully fetched X operators from spoke sno146"
   - ⚠️ "No kubeconfig secret for spoke sno146 on hub"
   - ❌ "Could not fetch operators from spoke sno146: ..."

3. **Verify spoke is accessible:**
   ```bash
   # From hub cluster
   oc --kubeconfig=/path/to/sno146/kubeconfig get csv -A
   ```

4. **Check browser console (F12):**
   - Look for fetch errors
   - Check network tab for API calls
   - Verify lazy loading is triggered

**Secret exists but still shows 0:**

1. Verify kubeconfig content:
   ```bash
   oc get secret -n sno146 sno146-admin-kubeconfig -o jsonpath='{.data.kubeconfig}' | base64 -d | head -10
   ```

2. Check RBAC permissions:
   - Backend needs access to secrets in spoke namespaces
   - Service account needs list/get permissions

3. Verify spoke cluster is reachable:
   - Network connectivity from hub to spoke
   - API server accessible
   - Certificates valid

**Operators not loading on expansion:**

1. Hard refresh browser (Ctrl+Shift+R)
2. Check browser console for JavaScript errors
3. Verify API endpoint returns data:
   ```bash
   curl -k https://monitor/api/hubs/acm1/spokes/sno146/operators
   ```
4. Check if version is current (v=20251107)

## Advanced Configuration

### Custom Kubeconfig Secrets

If your kubeconfig uses a different secret name:

The code expects: `{spoke-name}-admin-kubeconfig`

To use a different name, you would need to modify:
- `backend/pkg/client/rhacm.go` (line ~416)
- Look for: `cluster.Name + "-admin-kubeconfig"`

### Batch Setup Script

For multiple spokes:

```bash
#!/bin/bash
# setup-spoke-operators.sh

SPOKES=("sno146" "sno132" "sno133")

for spoke in "${SPOKES[@]}"; do
  echo "Setting up $spoke..."
  oc create namespace $spoke 2>/dev/null || true
  oc create secret generic ${spoke}-admin-kubeconfig \
    --from-file=kubeconfig=/kubeconfigs/${spoke}-kubeconfig \
    -n $spoke
  oc label secret ${spoke}-admin-kubeconfig \
    created-by=rhacm-monitor -n $spoke
done

echo "✅ Setup complete for ${#SPOKES[@]} spokes"
```

## Summary

- **Hub operators:** ✅ Working out of the box
- **Spoke operators:** ✅ Working with lazy loading (requires kubeconfig setup)
- **Performance:** ✅ Optimized for large scale (1000+ spokes)
- **Code:** ✅ Ready to fetch when secrets exist
- **Setup:** Create `{spoke}-admin-kubeconfig` secrets on hub

**Lazy loading ensures the application remains fast and responsive even with thousands of spoke clusters.**

---

*For questions or issues, check backend logs and browser console for detailed error messages.*
