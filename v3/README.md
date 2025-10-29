# RHACM Global Hub Monitor v3.0.0

**One-stop monitoring solution for Red Hat Advanced Cluster Management (RHACM) multi-cluster environments.**

## Quick Start

### 1. Deploy (One Command)

```bash
cd v3/deployment
chmod +x deploy.sh
./deploy.sh
```

That's it! The script will:
- ✅ Create namespace and RBAC
- ✅ Deploy backend (2 replicas)
- ✅ Deploy frontend (Nginx proxy)
- ✅ Create route
- ✅ Configure session affinity

### 2. Setup Hub Access (One-Time)

For each hub cluster you want to monitor:

```bash
# Create namespace
oc create namespace <hub-name>

# Create kubeconfig secret
oc create secret generic <hub-name>-admin-kubeconfig \
  --from-file=kubeconfig=/path/to/<hub-name>/kubeconfig \
  -n <hub-name>
```

**Example:**
```bash
oc create namespace acm1
oc create secret generic acm1-admin-kubeconfig \
  --from-file=kubeconfig=~/.kube/acm1-kubeconfig \
  -n acm1
```

### 3. Access

Get your URL:
```bash
oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}'
```

Open in browser: `https://<route-url>`

## Features

### v3.0.0 Highlights

🎨 **Dark/Light Mode**
- Toggle in header (top-right)
- Professional GitHub-inspired themes
- Persists across sessions

🔧 **Operators Monitoring**
- View all installed operators
- Hub operators: Full table with search
- Spoke operators: Lazy loaded for performance
- Smart grouping by name

⚡ **Performance**
- 300x faster with caching
- Lazy loading for 1000+ spokes
- Session affinity for consistency

🔄 **Per-Hub Refresh**
- 🔄 button on each hub card
- Refresh specific hub only
- Instant navigation

## What's Included

```
v3/
├── deployment/
│   ├── deploy.sh              # One-command deploy script
│   ├── k8s/
│   │   └── all-in-one.yaml    # All K8s resources
│   └── docker/
│       └── Dockerfile.backend.simple
├── frontend-static/
│   ├── index.html
│   ├── app.js
│   └── styles.css
└── backend/
    └── (Go source code)
```

## Container Images

**Backend:**
- `quay.io/bzhai/rhacm-monitor-backend:latest` (recommended - latest stable)
- `quay.io/bzhai/rhacm-monitor-backend:v3` (v3 branch)

**Frontend:**
- Static files served via Nginx
- No container build required

## Management

### Update Application

**Backend update:**
```bash
oc set image deployment/rhacm-monitor-backend \
  rhacm-monitor-backend=quay.io/bzhai/rhacm-monitor-backend:latest \
  -n rhacm-monitor
```

**Frontend update:**
```bash
cd v3/deployment
oc delete configmap frontend-html -n rhacm-monitor
oc create configmap frontend-html --from-file=../frontend-static/ -n rhacm-monitor
oc delete pods -l component=proxy -n rhacm-monitor
```

### Add Hub to Monitor

**Option 1: Managed Hub (Auto-discovered)**
- Create namespace with hub name
- Add kubeconfig secret
- Monitor auto-discovers it

**Option 2: Unmanaged Hub (Manual)**
- Click "➕ Add Hub" in UI
- Paste kubeconfig
- Submit

### Remove Deployment

```bash
oc delete namespace rhacm-monitor
oc delete clusterrole rhacm-monitor
oc delete clusterrolebinding rhacm-monitor
```

## Monitoring

### Check Status

```bash
# All pods
oc get pods -n rhacm-monitor

# Backend logs
oc logs -l component=backend -n rhacm-monitor -f

# Cache performance
oc logs -l component=backend -n rhacm-monitor | grep Cache
```

### Verify Performance

```bash
ROUTE_URL=$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')

# First request (cache miss)
time curl -k https://$ROUTE_URL/api/hubs

# Second request (should be fast)
time curl -k https://$ROUTE_URL/api/hubs
```

Expected:
- First: ~10 seconds
- Second: < 100ms

## Troubleshooting

### No Hubs Showing

**Check kubeconfig secrets:**
```bash
oc get secrets --all-namespaces | grep admin-kubeconfig
```

**Verify RBAC:**
```bash
oc auth can-i list managedclusters --as=system:serviceaccount:rhacm-monitor:rhacm-monitor
```

### Operators Not Showing

**Hub operators:** Should work automatically

**Spoke operators:** Need spoke kubeconfig secrets
```bash
# On hub cluster, for each spoke:
oc create namespace sno146
oc create secret generic sno146-admin-kubeconfig \
  --from-file=kubeconfig=/path/to/sno146/kubeconfig \
  -n sno146
```

See `OPERATORS_SETUP.md` for details.

### Slow Performance

**Check cache:**
```bash
oc logs -l component=backend -n rhacm-monitor | grep -E "Cache HIT|Cache MISS"
```

**Verify session affinity:**
```bash
oc get svc rhacm-monitor-backend -n rhacm-monitor -o yaml | grep sessionAffinity
```

Should show `sessionAffinity: ClientIP`

## Configuration

### Scale Replicas

```bash
# Backend
oc scale deployment/rhacm-monitor-backend --replicas=3 -n rhacm-monitor

# Frontend
oc scale deployment/rhacm-monitor-proxy --replicas=3 -n rhacm-monitor
```

### Adjust Resources

```bash
oc set resources deployment/rhacm-monitor-backend \
  -n rhacm-monitor \
  --requests=cpu=1,memory=1Gi \
  --limits=cpu=2,memory=2Gi
```

### Change Cache TTL

Edit `backend/pkg/handlers/hubs.go`:
```go
sharedCache := cache.NewCache(60 * time.Minute) // Change to 60 minutes
```

Rebuild and push image, then update deployment.

## Documentation

- **DEPLOYMENT.md** - Complete deployment guide
- **VERSION.md** - Feature documentation  
- **OPERATORS_SETUP.md** - Operators configuration
- **CACHE_ARCHITECTURE.md** - Cache design
- **FINAL_STATUS.md** - Project status
- **SPRINT_HISTORY.md** - Development journey

## Support

**Repository:** https://github.com/borball/rhacm-global-hub-monitor

**Issues:**
- Check logs: `oc logs -l component=backend -n rhacm-monitor`
- Browser console: F12 for JavaScript errors
- Review documentation above

## Summary

**Deploy in 2 steps:**
1. Run `./deploy.sh`
2. Create hub kubeconfig secrets

**Access:** via Route URL  
**Features:** Dark mode, Operators, Caching, Refresh  
**Status:** Production Ready ✅

---

*RHACM Global Hub Monitor v3.0.0 - Simple deployment, powerful monitoring.*
