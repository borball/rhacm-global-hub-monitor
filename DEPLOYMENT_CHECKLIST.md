# Deployment Checklist - Clean Deployment Procedure

## For Every Code Change

Follow these steps in order to ensure clean deployment:

### 1. Build Backend

```bash
cd v3/backend
go build -o bin/server cmd/server/main.go
```

**Verify:** No build errors

### 2. Build Docker Image

```bash
cd v3
podman build -t quay.io/bzhai/rhacm-monitor-backend:latest \
  -f deployment/docker/Dockerfile.backend.simple .
```

**Verify:** Image built successfully

### 3. Save Image

```bash
podman save quay.io/bzhai/rhacm-monitor-backend:latest \
  -o /tmp/rhacm-backend-deploy.tar
```

### 4. Deploy to All Cluster Nodes

```bash
for ip in 192.168.58.47 192.168.58.48 192.168.58.49; do
  scp -i ~/.ssh/hub_id.rsa /tmp/rhacm-backend-deploy.tar core@$ip:/tmp/b.tar
  ssh -i ~/.ssh/hub_id.rsa core@$ip \
    "sudo podman load -i /tmp/b.tar && sudo rm /tmp/b.tar"
  echo "✅ $ip"
done
```

**Verify:** All 3 nodes updated

### 5. Update Frontend (if changed)

```bash
cd v3
oc delete configmap frontend-html -n rhacm-monitor
oc create configmap frontend-html --from-file=frontend-static/ -n rhacm-monitor
```

### 6. Force Pod Restart

```bash
# Delete backend pods
oc delete pods -l component=backend -n rhacm-monitor

# Delete frontend pods (if frontend changed)
oc delete pods -l component=proxy -n rhacm-monitor

# Wait for new pods
sleep 40
```

### 7. Verify Deployment

```bash
# Check pods are running
oc get pods -n rhacm-monitor

# Check image ID matches
oc get pod -l component=backend -n rhacm-monitor \
  -o jsonpath='{.items[0].status.containerStatuses[0].imageID}'

# Test health endpoint
curl -k https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab/api/health
```

**Verify:** 
- All pods Running
- Health returns {"status":"healthy"}
- No errors in logs

### 8. Test Application

```bash
# Test API
curl -k https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab/api/hubs

# Test frontend
curl -k https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab/
```

**Manual UI Test:**
- Open browser
- Hard refresh (Ctrl+Shift+R)
- Test changed functionality

## Quick Deploy Script

Create `quick-deploy.sh`:

```bash
#!/bin/bash
set -e

echo "🚀 Clean Deployment Starting..."

# 1. Build
echo "📦 Building backend..."
cd v3/backend && go build -o bin/server cmd/server/main.go
cd ../..

# 2. Create image
echo "🐳 Building Docker image..."
cd v3
podman build -t quay.io/bzhai/rhacm-monitor-backend:latest \
  -f deployment/docker/Dockerfile.backend.simple . -q

# 3. Save
echo "💾 Saving image..."
podman save quay.io/bzhai/rhacm-monitor-backend:latest \
  -o /tmp/rhacm-deploy.tar

# 4. Deploy to nodes
echo "📤 Deploying to cluster nodes..."
for ip in 192.168.58.47 192.168.58.48 192.168.58.49; do
  scp -i ~/.ssh/hub_id.rsa -q /tmp/rhacm-deploy.tar core@$ip:/tmp/b.tar
  ssh -i ~/.ssh/hub_id.rsa -q core@$ip \
    "sudo podman load -i /tmp/b.tar && sudo rm /tmp/b.tar"
  echo "  ✅ $ip"
done

# 5. Update frontend
echo "🎨 Updating frontend..."
oc delete configmap frontend-html -n rhacm-monitor 2>/dev/null || true
oc create configmap frontend-html --from-file=frontend-static/ -n rhacm-monitor

# 6. Restart pods
echo "🔄 Restarting pods..."
oc delete pods -l component=backend -n rhacm-monitor
oc delete pods -l component=proxy -n rhacm-monitor

echo "⏳ Waiting for pods to start..."
sleep 40

# 7. Verify
echo "✅ Verifying deployment..."
oc get pods -n rhacm-monitor

echo ""
echo "🧪 Testing health endpoint..."
curl -k -s https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab/api/health | jq '.'

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║                                                            ║"
echo "║    ✅ DEPLOYMENT COMPLETE!                                ║"
echo "║                                                            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "URL: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab"
echo ""
echo "Remember to hard refresh your browser (Ctrl+Shift+R)!"
```

## Common Issues

### Issue: Old Image Still Running

**Symptom:** Changes not taking effect

**Solution:**
```bash
# Force image pull by changing tag
podman tag quay.io/bzhai/rhacm-monitor-backend:latest \
  quay.io/bzhai/rhacm-monitor-backend:v3-$(date +%s)

# Then follow deployment steps
```

### Issue: Pods Not Restarting

**Symptom:** Pods stuck in old state

**Solution:**
```bash
# Force delete with grace period 0
oc delete pods -l component=backend -n rhacm-monitor \
  --force --grace-period=0
```

### Issue: ConfigMap Not Updating

**Symptom:** Frontend showing old files

**Solution:**
```bash
# Delete and recreate (not just update)
oc delete configmap frontend-html -n rhacm-monitor
oc create configmap frontend-html --from-file=frontend-static/ -n rhacm-monitor

# Force proxy restart
oc delete pods -l component=proxy -n rhacm-monitor
```

## Best Practices

1. **Always build fresh** - Don't rely on cached layers
2. **Tag images with timestamp** - Ensure uniqueness
3. **Delete pods, don't just restart** - Force image pull
4. **Wait adequate time** - Allow 40s for pod startup
5. **Verify before testing** - Check pods are Running
6. **Hard refresh browser** - Clear client-side cache

## Verification Commands

```bash
# Check backend image ID
oc get pod -l component=backend -n rhacm-monitor \
  -o jsonpath='{.items[0].status.containerStatuses[0].imageID}'

# Compare with local image
podman images quay.io/bzhai/rhacm-monitor-backend:latest \
  --format "{{.ID}}"

# Should match (first 12 chars)

# Check logs for startup messages
oc logs -l component=backend -n rhacm-monitor --tail=20

# Test API immediately
curl -k https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab/api/health
```

---

**Always follow this checklist for clean deployments!**

