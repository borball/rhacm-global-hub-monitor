# RHACM Global Hub Monitor - Deployment Status

**Date**: October 20, 2025  
**Application URL**: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab  
**Repository**: github.com:borball/rhacm-global-hub-monitor.git

## Current Deployment

### ✅ Fully Deployed and Working

**v1 Features Active:**
- Hub and spoke monitoring (2 hubs, 1 spoke)
- Policy management with search, filter, details
- Policy YAML download from live cluster
- Policy enforcement (CGU/TALM integration)
- Configuration version tracking and filtering
- Policy status messages (latest violations/notifications)
- Correct compliance calculation (98%, not 100%)
- Improved UI layouts (compact, efficient)
- Node information (K8s + BMH merged)

**Components Running:**
- Backend: 2 pods (with CGU, policy status, config version)
- Frontend Proxy: 2 pods (with dual-method UI)
- Route: hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab

### ⏳ Code Complete, Awaiting Backend Deployment

**Add Hub Feature:**
- ✅ Code complete and committed to GitHub
- ✅ Frontend UI deployed (visible in Incognito)
- ❌ Backend endpoints not deployed (quay.io read-only)

**Backend Endpoints Pending:**
```
POST   /api/hubs/add          - Add new hub
DELETE /api/hubs/{name}/remove - Remove hub
```

**Reason**: Quay.io container registry is currently read-only. Cannot push updated backend image.

**When Quay.io Becomes Writable:**
```bash
cd /root/workspace/github/rhacm-global-hub-monitor/v1
podman push quay.io/bzhai/rhacm-monitor-backend:latest
oc rollout restart deployment/rhacm-monitor-backend -n rhacm-monitor
```

Then the Add Hub UI feature will work immediately.

## Workaround for Adding Hubs

Until backend is deployed, hubs can be added manually:

```bash
# Create namespace
oc create namespace {hub-name}

# Create kubeconfig secret
oc create secret generic {hub-name}-admin-kubeconfig \
  --from-file=kubeconfig={path-to-kubeconfig} \
  -n {hub-name}

# Hub will appear in application automatically
```

## Feature Summary

### Working Features (Deployed)
1. ✅ Hub monitoring (acm1, acm2)
2. ✅ Spoke monitoring (sno146) with scalable table
3. ✅ Policy compliance tracking (45 policies, 98% compliant)
4. ✅ Policy enforcement via CGU/TALM
5. ✅ Configuration version display and search
6. ✅ Policy status messages
7. ✅ Node information (7 nodes, K8s + BMH merged)
8. ✅ Search/filter (spokes: 3 fields, policies: 2 fields)
9. ✅ Policy YAML download
10. ✅ Improved UI/UX

### Ready Features (Code Complete)
11. ⏳ Add Hub via kubeconfig upload
12. ⏳ Add Hub via API credentials
13. ⏳ Unmanaged Hubs section
14. ⏳ Remove Hub functionality

## Next Steps

1. Wait for quay.io to become writable
2. Push backend image: `podman push quay.io/bzhai/rhacm-monitor-backend:latest`
3. Restart backend: `oc rollout restart deployment/rhacm-monitor-backend -n rhacm-monitor`
4. Test Add Hub feature in browser

---

**10 out of 14 features are live and working. 4 features are code-complete and awaiting deployment.**
