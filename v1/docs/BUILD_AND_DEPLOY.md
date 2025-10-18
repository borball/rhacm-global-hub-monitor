# Build and Deploy Guide

This guide explains how to build container images and deploy the RHACM Global Hub Monitor application.

## Container Images

All images use Red Hat/Quay.io base images to avoid Docker Hub rate limits:
- **Backend**: Uses `registry.access.redhat.com/ubi9/ubi-minimal:latest` (Red Hat Universal Base Image)
- **Frontend**: Uses `registry.access.redhat.com/ubi9/nodejs-20` and `registry.access.redhat.com/ubi9/nginx-124`

## Building Images

### Option 1: Quick Build (Recommended for Development)

Build the backend binary locally, then create a minimal container:

```bash
# Build backend binary
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server cmd/server/main.go

# Build container image (uses simple Dockerfile)
cd ..
podman build -f deployment/docker/Dockerfile.backend.simple \
  -t quay.io/bzhai/rhacm-monitor-backend:latest .
```

**Advantages:**
- Very fast (< 10 seconds for container build)
- Uses minimal Red Hat UBI base image (~28 MB)
- No dependency downloads in container build
- Great for iterative development

### Option 2: Multi-stage Build (For CI/CD)

Use the multi-stage Dockerfile that builds everything in the container:

```bash
podman build -f deployment/docker/Dockerfile.backend.multistage \
  -t quay.io/bzhai/rhacm-monitor-backend:latest .
```

**Note:** This uses Fedora (from quay.io) as the builder image, which has Go 1.22+. Build time: ~5-10 minutes.

### Frontend Build

```bash
podman build -f deployment/docker/Dockerfile.frontend \
  -t quay.io/bzhai/rhacm-monitor-frontend:latest .
```

## Automated Build Script

Use the provided script to build and push all images:

```bash
./build-and-push.sh
```

To build with a specific version tag:

```bash
VERSION=v1.0.0 ./build-and-push.sh
```

The script will:
1. Build backend image
2. Build frontend image
3. Tag images with both version and `:latest`
4. Push to `quay.io/bzhai/`

## Dockerfile Variants

### Backend Dockerfiles

1. **Dockerfile.backend** - Current: Multi-stage with Fedora builder
   - Uses: `quay.io/fedora/fedora:latest` (builder) + `registry.access.redhat.com/ubi9/ubi-minimal:latest` (runtime)
   - Build time: ~5-10 minutes
   - Best for: CI/CD pipelines

2. **Dockerfile.backend.simple** - Recommended for development
   - Uses: `registry.access.redhat.com/ubi9/ubi-minimal:latest`
   - Requires: Prebuilt binary
   - Build time: < 10 seconds
   - Best for: Local development and testing

3. **Dockerfile.backend.multistage** - Same as Dockerfile.backend
   - Multi-stage build with all dependencies

### Frontend Dockerfile

- **Dockerfile.frontend**
  - Builder: `registry.access.redhat.com/ubi9/nodejs-20:latest`
  - Runtime: `registry.access.redhat.com/ubi9/nginx-124:latest`
  - Build time: ~2-5 minutes

## Pushing to Quay.io

### Login to Quay.io

```bash
podman login quay.io
# Enter username: bzhai
# Enter password: <your-token>
```

### Push Images

```bash
podman push quay.io/bzhai/rhacm-monitor-backend:latest
podman push quay.io/bzhai/rhacm-monitor-frontend:latest
```

## Deploying to OpenShift

### Prerequisites

```bash
# Ensure you're connected to the correct cluster
oc whoami
oc project

# Switch to vhub (if needed)
vhub
```

### Method 1: Using Kustomize (Recommended)

```bash
# Deploy everything
oc apply -k deployment/k8s/

# Check deployment status
oc get pods -n rhacm-monitor
oc get route -n rhacm-monitor
```

### Method 2: Manual Deployment

```bash
# Create namespace and RBAC
oc apply -f deployment/k8s/namespace.yaml
oc apply -f deployment/k8s/serviceaccount.yaml
oc apply -f deployment/k8s/clusterrole.yaml
oc apply -f deployment/k8s/clusterrolebinding.yaml

# Deploy backend
oc apply -f deployment/k8s/backend-deployment.yaml
oc apply -f deployment/k8s/backend-service.yaml

# Wait for backend to be ready
oc wait --for=condition=available --timeout=300s \
  deployment/rhacm-monitor-backend -n rhacm-monitor

# Deploy frontend
oc apply -f deployment/k8s/frontend-deployment.yaml
oc apply -f deployment/k8s/frontend-service.yaml

# Create route
oc apply -f deployment/k8s/route.yaml
```

### Method 3: Using Operator (Future)

```bash
# Install CRD
oc apply -f operator/config/crd/rhacmmonitor-crd.yaml

# Install operator
oc apply -f operator/config/rbac/
oc apply -f operator/config/manager/

# Create instance
oc apply -f operator/config/samples/rhacmmonitor-sample.yaml
```

## Verifying Deployment

### Check Pods

```bash
oc get pods -n rhacm-monitor
```

Expected output:
```
NAME                                      READY   STATUS    RESTARTS   AGE
rhacm-monitor-backend-xxxxx-xxxxx        1/1     Running   0          2m
rhacm-monitor-backend-xxxxx-xxxxx        1/1     Running   0          2m
rhacm-monitor-frontend-xxxxx-xxxxx       1/1     Running   0          2m
rhacm-monitor-frontend-xxxxx-xxxxx       1/1     Running   0          2m
```

### Check Backend Logs

```bash
oc logs -l component=backend -n rhacm-monitor --tail=50
```

### Test API

```bash
# Get route URL
ROUTE=$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')

# Test health endpoint
curl -k https://$ROUTE/api/health

# Test hubs endpoint
curl -k https://$ROUTE/api/hubs | jq .
```

### Access UI

```bash
# Get and open URL
echo "https://$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')"
```

## Updating Deployment

### Update Backend Image

```bash
oc set image deployment/rhacm-monitor-backend \
  backend=quay.io/bzhai/rhacm-monitor-backend:v1.1.0 \
  -n rhacm-monitor
```

### Update Frontend Image

```bash
oc set image deployment/rhacm-monitor-frontend \
  frontend=quay.io/bzhai/rhacm-monitor-frontend:v1.1.0 \
  -n rhacm-monitor
```

### Rollback

```bash
# Rollback backend
oc rollout undo deployment/rhacm-monitor-backend -n rhacm-monitor

# Rollback frontend
oc rollout undo deployment/rhacm-monitor-frontend -n rhacm-monitor
```

## Building with OpenShift BuildConfig (Alternative)

If you prefer to build images directly in OpenShift:

```bash
# Create build config
oc new-build --name=rhacm-monitor-backend \
  --binary=true \
  --strategy=docker \
  -n rhacm-monitor

# Start build from local directory
oc start-build rhacm-monitor-backend \
  --from-dir=. \
  --follow \
  -n rhacm-monitor
```

## Image Size Comparison

| Image | Base | Size |
|-------|------|------|
| Backend (simple) | UBI9 Minimal | ~86 MB |
| Backend (multistage) | UBI9 Minimal | ~86 MB |
| Frontend | UBI9 nginx | ~350 MB |

## Troubleshooting

### Build Issues

**Problem**: "Error: short-name resolution enforced"
```bash
# Solution: Use fully qualified image names
podman build --pull -f deployment/docker/Dockerfile.backend.simple .
```

**Problem**: Build timeout
```bash
# Solution: Increase timeout or use simple Dockerfile
timeout 600 podman build ...
```

**Problem**: "go.mod requires go >= 1.22"
```bash
# Solution: Use Fedora-based builder or build binary locally
cd backend && go build -o bin/server cmd/server/main.go
```

### Deployment Issues

**Problem**: ImagePullBackOff
```bash
# Check if images are public or create pull secret
oc create secret docker-registry quay-pull-secret \
  --docker-server=quay.io \
  --docker-username=bzhai \
  --docker-password=<token> \
  -n rhacm-monitor

# Link secret to service account
oc secrets link rhacm-monitor quay-pull-secret --for=pull -n rhacm-monitor
```

**Problem**: CrashLoopBackOff
```bash
# Check logs
oc logs -l component=backend -n rhacm-monitor --tail=100

# Common issues:
# - Missing RBAC permissions
# - Can't connect to Kubernetes API
# - Missing required environment variables
```

## Best Practices

1. **Use version tags** for production deployments, not `:latest`
2. **Build locally first** for faster iteration during development
3. **Test images** before pushing to registry
4. **Use BuildConfig** in OpenShift for automated builds from Git
5. **Keep images small** by using minimal base images
6. **Scan images** for security vulnerabilities before deploying

## CI/CD Integration

Example GitHub Actions workflow:

```yaml
name: Build and Push Images

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Login to Quay.io
        uses: docker/login-action@v2
        with:
          registry: quay.io
          username: ${{ secrets.QUAY_USERNAME }}
          password: ${{ secrets.QUAY_PASSWORD }}
      
      - name: Build and Push
        run: |
          VERSION=${GITHUB_REF#refs/tags/}
          export VERSION
          ./build-and-push.sh
```

## Additional Resources

- [Red Hat Container Catalog](https://catalog.redhat.com/software/containers/explore)
- [UBI Documentation](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/building_running_and_managing_containers/)
- [OpenShift Build Documentation](https://docs.openshift.com/container-platform/latest/cicd/builds/understanding-image-builds.html)

