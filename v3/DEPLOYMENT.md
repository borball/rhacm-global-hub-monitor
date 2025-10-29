# RHACM Global Hub Monitor v3.0.0 - Deployment Guide

## Overview

This guide covers deploying the RHACM Global Hub Monitor v3.0.0 to your OpenShift/Kubernetes cluster.

## Prerequisites

- OpenShift 4.x or Kubernetes 1.24+
- RHACM (Red Hat Advanced Cluster Management) installed
- Access to managed hub clusters
- Kubeconfig secrets for hubs to monitor
- Cluster-admin or sufficient RBAC permissions

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│ Global Hub (vhub)                                               │
│                                                                 │
│  ┌──────────────────┐      ┌─────────────────┐                │
│  │  RHACM Monitor   │      │  RHACM Monitor  │                │
│  │  Backend (x2)    │◄────►│  Frontend (x2)  │                │
│  │  Port: 8080      │      │  Nginx Proxy    │                │
│  └──────────────────┘      └─────────────────┘                │
│         │                           │                           │
│         │                           │                           │
│         ▼                           ▼                           │
│  ┌──────────────────┐      ┌─────────────────┐                │
│  │  Hub Kubeconfigs │      │  Static Files   │                │
│  │  (Secrets)       │      │  (ConfigMap)    │                │
│  └──────────────────┘      └─────────────────┘                │
└─────────────────────────────────────────────────────────────────┘
                    │
                    ├──► Hub: acm1 (monitors spoke clusters)
                    ├──► Hub: acm2 (monitors spoke clusters)  
                    └──► Hub: production-hub (unmanaged)
```

## Container Images

### Backend Image

**Repository:** `quay.io/bzhai/rhacm-monitor-backend`

**Available Tags:**
- `latest` - Latest stable version (recommended)
- `v3` - v3 branch latest

**Image Contents:**
- Go binary (backend server)
- Port: 8080
- Health endpoints: /api/health, /api/ready, /api/live

### Frontend Files

**Deployment:** ConfigMap with static files
- `index.html` - Main HTML page
- `app.js` - Application JavaScript (v=20251107)
- `styles.css` - Styles with dark/light mode

## Deployment Steps

### 1. Create Namespace

```bash
oc create namespace rhacm-monitor
```

### 2. Create Service Account

```bash
cat <<EOF | oc apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: rhacm-monitor
  namespace: rhacm-monitor
EOF
```

### 3. Create RBAC Permissions

```bash
cat <<EOF | oc apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: rhacm-monitor
rules:
- apiGroups: [""]
  resources: ["nodes", "secrets", "namespaces"]
  verbs: ["get", "list"]
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "list", "create", "delete"]
- apiGroups: ["cluster.open-cluster-management.io"]
  resources: ["managedclusters"]
  verbs: ["get", "list"]
- apiGroups: ["policy.open-cluster-management.io"]
  resources: ["policies"]
  verbs: ["get", "list"]
- apiGroups: ["operators.coreos.com"]
  resources: ["clusterserviceversions"]
  verbs: ["get", "list"]
- apiGroups: ["metal3.io"]
  resources: ["baremetalhosts"]
  verbs: ["get", "list"]
- apiGroups: ["route.openshift.io"]
  resources: ["routes"]
  verbs: ["get", "list"]
- apiGroups: ["config.openshift.io"]
  resources: ["clusterversions"]
  verbs: ["get", "list"]
- apiGroups: ["ran.openshift.io"]
  resources: ["clustergroupupgrades"]
  verbs: ["get", "list", "create"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: rhacm-monitor
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: rhacm-monitor
subjects:
- kind: ServiceAccount
  name: rhacm-monitor
  namespace: rhacm-monitor
EOF
```

### 4. Create Hub Kubeconfig Secrets

For each hub you want to monitor:

```bash
# Example for acm1 hub
oc create namespace acm1

oc create secret generic acm1-admin-kubeconfig \
  --from-file=kubeconfig=/path/to/acm1/kubeconfig \
  -n acm1

oc label secret acm1-admin-kubeconfig \
  created-by=rhacm-monitor \
  -n acm1
```

Repeat for each hub (acm2, production-hub, etc.)

### 5. Deploy Backend

```bash
cat <<EOF | oc apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rhacm-monitor-backend
  namespace: rhacm-monitor
  labels:
    app: rhacm-monitor
    component: backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: rhacm-monitor
      component: backend
  template:
    metadata:
      labels:
        app: rhacm-monitor
        component: backend
    spec:
      serviceAccountName: rhacm-monitor
      containers:
      - name: rhacm-monitor-backend
        image: quay.io/bzhai/rhacm-monitor-backend:v3.0.0
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        env:
        - name: PORT
          value: "8080"
        - name: ENABLE_AUTH
          value: "false"
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            cpu: 500m
            memory: 512Mi
          limits:
            cpu: 1
            memory: 1Gi
        livenessProbe:
          httpGet:
            path: /api/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          capabilities:
            drop:
            - ALL
EOF
```

### 6. Create Backend Service (with Session Affinity)

```bash
cat <<EOF | oc apply -f -
apiVersion: v1
kind: Service
metadata:
  name: rhacm-monitor-backend
  namespace: rhacm-monitor
  labels:
    app: rhacm-monitor
    component: backend
spec:
  selector:
    app: rhacm-monitor
    component: backend
  ports:
  - name: http
    port: 8080
    protocol: TCP
    targetPort: 8080
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800  # 3 hours for consistent cache
  type: ClusterIP
EOF
```

**Note:** Session affinity is critical for consistent cache performance with multiple backend replicas.

### 7. Deploy Frontend Files

Create ConfigMap with static frontend files:

```bash
cd /path/to/rhacm-global-hub-monitor/v3

oc create configmap frontend-html \
  --from-file=frontend-static/ \
  -n rhacm-monitor
```

### 8. Deploy Nginx Proxy

```bash
cat <<EOF | oc apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rhacm-monitor-proxy
  namespace: rhacm-monitor
  labels:
    app: rhacm-monitor
    component: proxy
spec:
  replicas: 2
  selector:
    matchLabels:
      app: rhacm-monitor
      component: proxy
  template:
    metadata:
      labels:
        app: rhacm-monitor
        component: proxy
    spec:
      containers:
      - name: nginx
        image: registry.access.redhat.com/ubi9/nginx-120:latest
        ports:
        - containerPort: 8080
          name: http
        volumeMounts:
        - name: html
          mountPath: /var/www/html
        - name: nginx-conf
          mountPath: /etc/nginx/nginx.conf
          subPath: nginx.conf
      volumes:
      - name: html
        configMap:
          name: frontend-html
      - name: nginx-conf
        configMap:
          name: nginx-config
EOF
```

### 9. Create Nginx Configuration

```bash
cat <<EOF | oc apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-config
  namespace: rhacm-monitor
data:
  nginx.conf: |
    worker_processes auto;
    error_log /dev/stderr info;
    pid /tmp/nginx.pid;

    events {
        worker_connections 1024;
    }

    http {
        include /etc/nginx/mime.types;
        default_type application/octet-stream;
        
        access_log /dev/stdout;
        
        sendfile on;
        keepalive_timeout 65;
        
        server {
            listen 8080;
            server_name _;
            
            root /var/www/html;
            index index.html;
            
            location / {
                try_files \$uri \$uri/ /index.html;
            }
            
            location /api/ {
                proxy_pass http://rhacm-monitor-backend:8080;
                proxy_http_version 1.1;
                proxy_set_header Upgrade \$http_upgrade;
                proxy_set_header Connection 'upgrade';
                proxy_set_header Host \$host;
                proxy_cache_bypass \$http_upgrade;
                proxy_set_header X-Real-IP \$remote_addr;
                proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
                proxy_read_timeout 300s;
                proxy_connect_timeout 75s;
            }
        }
    }
EOF
```

### 10. Create Frontend Service

```bash
cat <<EOF | oc apply -f -
apiVersion: v1
kind: Service
metadata:
  name: rhacm-monitor-frontend
  namespace: rhacm-monitor
  labels:
    app: rhacm-monitor
    component: proxy
spec:
  selector:
    app: rhacm-monitor
    component: proxy
  ports:
  - name: http
    port: 8080
    protocol: TCP
    targetPort: 8080
  type: ClusterIP
EOF
```

### 11. Create Route (OpenShift)

```bash
cat <<EOF | oc apply -f -
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: rhacm-monitor
  namespace: rhacm-monitor
  labels:
    app: rhacm-monitor
spec:
  to:
    kind: Service
    name: rhacm-monitor-frontend
    weight: 100
  port:
    targetPort: http
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
  wildcardPolicy: None
EOF
```

## Post-Deployment Configuration

### Setup Spoke Operator Monitoring

For spoke clusters to show operators (optional):

```bash
# On each hub cluster (e.g., acm1)
# For each spoke (e.g., sno146)

oc create namespace sno146

oc create secret generic sno146-admin-kubeconfig \
  --from-file=kubeconfig=/path/to/sno146/kubeconfig \
  -n sno146

oc label secret sno146-admin-kubeconfig \
  created-by=rhacm-monitor \
  -n sno146
```

See `OPERATORS_SETUP.md` for detailed instructions.

## Verification

### 1. Check Pod Status

```bash
oc get pods -n rhacm-monitor
```

Expected output:
```
NAME                                      READY   STATUS    RESTARTS   AGE
rhacm-monitor-backend-xxx                 1/1     Running   0          2m
rhacm-monitor-backend-yyy                 1/1     Running   0          2m
rhacm-monitor-proxy-xxx                   1/1     Running   0          2m
rhacm-monitor-proxy-yyy                   1/1     Running   0          2m
```

### 2. Check Services

```bash
oc get svc -n rhacm-monitor
```

Verify session affinity:
```bash
oc get svc rhacm-monitor-backend -n rhacm-monitor -o yaml | grep sessionAffinity
```

Should show:
```yaml
sessionAffinity: ClientIP
```

### 3. Test Backend Health

```bash
oc get route -n rhacm-monitor
ROUTE_URL=$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')

curl -k https://$ROUTE_URL/api/health
```

Expected response:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2025-10-21T18:00:00Z"
}
```

### 4. Test Frontend

Visit in browser: `https://<route-url>`

You should see:
- Statistics cards at top
- Managed Hubs section
- Unmanaged Hubs section
- Theme toggle button (top-right)

### 5. Test Cache

```bash
# First request (slow)
time curl -k https://$ROUTE_URL/api/hubs

# Second request (should be fast)
time curl -k https://$ROUTE_URL/api/hubs
```

Expected:
- First: ~10 seconds
- Second: < 100ms

## Updating Deployment

### Update Backend

```bash
# Using latest image
oc set image deployment/rhacm-monitor-backend \
  rhacm-monitor-backend=quay.io/bzhai/rhacm-monitor-backend:latest \
  -n rhacm-monitor

# Verify rollout
oc rollout status deployment/rhacm-monitor-backend -n rhacm-monitor
```

### Update Frontend

```bash
cd /path/to/rhacm-global-hub-monitor/v3

# Delete old configmap
oc delete configmap frontend-html -n rhacm-monitor

# Create new configmap
oc create configmap frontend-html \
  --from-file=frontend-static/ \
  -n rhacm-monitor

# Restart proxy pods
oc delete pods -l component=proxy -n rhacm-monitor
```

## Configuration

### Environment Variables (Backend)

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Backend API port |
| ENABLE_AUTH | false | Enable JWT authentication |
| OAUTH_ISSUER_URL | https://kubernetes.default.svc | OAuth issuer |
| LOG_LEVEL | info | Logging level (debug/info/warn/error) |

### Cache Configuration

**Cache TTL:** 30 minutes (hardcoded in backend)
- Location: `backend/pkg/handlers/hubs.go`
- Setting: `cache.NewCache(30 * time.Minute)`

**Session Affinity:** 3 hours
- Location: Backend service spec
- Setting: `sessionAffinityConfig.clientIP.timeoutSeconds: 10800`

### Frontend Version

**Cache Busting:** Update version in `index.html`
```html
<script src="app.js?v=20251107"></script>
```

Change version number to force browser reload of JavaScript.

## Monitoring

### Backend Logs

```bash
# View logs
oc logs -f deployment/rhacm-monitor-backend -n rhacm-monitor

# Cache logging
oc logs -l component=backend -n rhacm-monitor | grep Cache

# Performance monitoring
oc logs -l component=backend -n rhacm-monitor | grep "GET.*hubs"
```

### Cache Performance

**Fast requests (cache hit):**
```
[GIN] 18:09:28 | 200 | 3.656ms | GET "/api/hubs"
2025/10/21 18:09:28 Cache HIT for hubs:list
```

**Slow requests (cache miss):**
```
[GIN] 18:00:04 | 200 | 10.665s | GET "/api/hubs"
2025/10/21 18:00:04 Cache MISS for hubs:list - fetching fresh data
```

### Resource Usage

Monitor pod resources:
```bash
oc adm top pods -n rhacm-monitor
```

Expected:
- Backend: ~300-500Mi memory, ~0.2-0.5 CPU
- Proxy: ~50-100Mi memory, ~0.1 CPU

## Troubleshooting

### Backend Pods CrashLooping

**Check logs:**
```bash
oc logs -l component=backend -n rhacm-monitor --tail=100
```

**Common causes:**
- Missing hub kubeconfig secrets
- RBAC permissions insufficient
- Image pull errors

### 404 Errors

**Symptom:** API calls return 404

**Check:**
- Route is created: `oc get route -n rhacm-monitor`
- Backend service exists: `oc get svc rhacm-monitor-backend -n rhacm-monitor`
- Pods are running: `oc get pods -n rhacm-monitor`

### Slow Performance

**Symptom:** Page loads slowly consistently

**Checks:**
1. Verify cache is enabled (check logs for "Cache HIT")
2. Check session affinity: `oc get svc rhacm-monitor-backend -o yaml | grep sessionAffinity`
3. Verify imagePullPolicy: Always (pulls latest code)
4. Check pod resource limits

**Debug:**
```bash
# Check cache behavior
oc logs -l component=backend -n rhacm-monitor | grep -E "Cache|GET.*hubs" | tail -30
```

### Operators Not Showing

**For Hub Operators:**
- Should work out of the box
- Check RBAC for ClusterServiceVersion access

**For Spoke Operators:**
- Requires spoke kubeconfig secrets
- See `OPERATORS_SETUP.md`
- Check: `oc get secret -n {spoke-name} {spoke-name}-admin-kubeconfig`

### Dark Mode Issues

**Symptom:** White backgrounds in dark mode

**Check:**
- Hard refresh browser: Ctrl+Shift+R
- Verify version: v=20251107 in page source
- Check styles.css is loaded

## Scaling

### Horizontal Scaling

**Backend:**
```bash
oc scale deployment/rhacm-monitor-backend --replicas=3 -n rhacm-monitor
```

**Note:** Session affinity ensures cache consistency per client.

**Frontend:**
```bash
oc scale deployment/rhacm-monitor-proxy --replicas=3 -n rhacm-monitor
```

### Resource Tuning

For large deployments (10+ hubs, 1000+ spokes):

```bash
oc set resources deployment/rhacm-monitor-backend \
  -n rhacm-monitor \
  --requests=cpu=1,memory=1Gi \
  --limits=cpu=2,memory=2Gi
```

## Backup & Recovery

### Backup Unmanaged Hub Configurations

```bash
# Export hub kubeconfig secrets
oc get secrets -n rhacm-monitor -l created-by=rhacm-monitor -o yaml > hub-secrets-backup.yaml
```

### Restore

```bash
oc apply -f hub-secrets-backup.yaml
```

## Upgrades

### From v2 to v3

v3 is backward compatible with v2. Simple image update:

```bash
# Update backend
oc set image deployment/rhacm-monitor-backend \
  rhacm-monitor-backend=quay.io/bzhai/rhacm-monitor-backend:v3.0.0 \
  -n rhacm-monitor

# Update frontend
oc delete configmap frontend-html -n rhacm-monitor
oc create configmap frontend-html --from-file=v3/frontend-static/ -n rhacm-monitor
oc delete pods -l component=proxy -n rhacm-monitor
```

**New in v3:**
- Dark/light mode (automatic)
- Operators tab (automatic for hubs)
- Per-hub refresh buttons
- Session affinity recommended

## Advanced Configuration

### Enable Authentication

Set `ENABLE_AUTH=true` in backend deployment.

Requires:
- OAuth/OIDC provider
- JWT token validation
- Frontend login flow

(Authentication is optional - disabled by default)

### Custom Cache TTL

Modify `backend/pkg/handlers/hubs.go`:
```go
sharedCache := cache.NewCache(60 * time.Minute) // 60 minutes
```

Rebuild and redeploy backend image.

### Custom Session Affinity Timeout

Modify backend service:
```yaml
sessionAffinityConfig:
  clientIP:
    timeoutSeconds: 21600  # 6 hours
```

## Security Considerations

1. **Service Account Permissions:** Principle of least privilege applied
2. **Network Policies:** Consider adding for production
3. **TLS:** Route uses edge termination (recommended)
4. **Secrets:** Kubeconfigs stored as Kubernetes secrets (encrypted at rest)
5. **Container Security:** Runs as non-root, drops all capabilities

## Performance Tuning

### Cache Strategy

**Current:**
- TTL: 30 minutes
- Session affinity: 3 hours
- In-memory per pod

**Tuning:**
- Increase TTL if data freshness less critical
- Decrease TTL if you need more current data
- Use refresh buttons for on-demand updates

### Lazy Loading

Spoke operators use lazy loading by default. No configuration needed.

**Behavior:**
- Table shows "..." initially
- Operators fetched on spoke expansion
- ~2 seconds per spoke
- Cached in browser after first load

## Quick Reference

### Useful Commands

```bash
# View all resources
oc get all -n rhacm-monitor

# Restart backend (pull latest image)
oc rollout restart deployment/rhacm-monitor-backend -n rhacm-monitor

# Update frontend
oc delete configmap frontend-html -n rhacm-monitor && \
oc create configmap frontend-html --from-file=v3/frontend-static/ -n rhacm-monitor && \
oc delete pods -l component=proxy -n rhacm-monitor

# Check cache performance
oc logs -l component=backend -n rhacm-monitor | grep Cache | tail -20

# Monitor requests
oc logs -l component=backend -n rhacm-monitor -f | grep "GET.*hubs"
```

### Access URLs

```bash
# Get route URL
oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}'

# Full URL
echo "https://$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')"
```

## Support

### Documentation

- `VERSION.md` - Complete feature documentation
- `OPERATORS_SETUP.md` - Operators feature setup
- `CACHE_ARCHITECTURE.md` - Caching design
- `FINAL_STATUS.md` - Project completion status
- `SPRINT_HISTORY.md` - Development journey

### Logs

**Backend logs:**
```bash
oc logs -l component=backend -n rhacm-monitor --tail=100
```

**Frontend logs:**
```bash
oc logs -l component=proxy -n rhacm-monitor --tail=100
```

**Browser console:** F12 in browser for JavaScript errors

---

## Summary

**v3.0.0 Deployment:**
- 2 backend pods (HA with session affinity)
- 2 frontend pods (Nginx proxy)
- ConfigMap for static files
- Session affinity for cache consistency
- imagePullPolicy: Always for easy updates

**Key Features:**
- Dark/light mode theming
- Operators monitoring (hub + spoke)
- 30-minute caching
- Per-hub refresh
- Lazy loading for performance

**Production Ready:** ✅

---

*For issues or questions, check logs and documentation, or refer to SPRINT_HISTORY.md for historical context.*

