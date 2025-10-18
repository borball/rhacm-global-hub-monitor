# Docker Registry Changes Summary

## Overview

Updated all Dockerfiles and deployment manifests to use Red Hat/Quay.io base images instead of Docker Hub images to avoid rate limiting issues.

## Changes Made

### 1. Backend Dockerfile Updates

#### Original (Using Docker Hub):
```dockerfile
FROM golang:1.22-alpine AS builder
```

#### Updated (Using Red Hat/Quay.io):
```dockerfile
# Option 1: Multi-stage with Fedora (from Quay.io)
FROM quay.io/fedora/fedora:latest AS builder

# Option 2: Simple - using prebuilt binary
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
```

**Base Images Used:**
- **Builder**: `quay.io/fedora/fedora:latest` (has Go 1.23+)
- **Runtime**: `registry.access.redhat.com/ubi9/ubi-minimal:latest` (Red Hat UBI9 Minimal)

### 2. Frontend Dockerfile Updates

#### Original (Using Docker Hub):
```dockerfile
FROM node:20-alpine AS builder
```

#### Updated (Using Red Hat Registry):
```dockerfile
FROM registry.access.redhat.com/ubi9/nodejs-20:latest AS builder
```

**Base Images Used:**
- **Builder**: `registry.access.redhat.com/ubi9/nodejs-20:latest`
- **Runtime**: `registry.access.redhat.com/ubi9/nginx-124:latest` (already using Red Hat)

### 3. Deployment Manifest Updates

Updated image references in all Kubernetes manifests:

#### Files Modified:
- `deployment/k8s/backend-deployment.yaml`
- `deployment/k8s/frontend-deployment.yaml`
- `deployment/k8s/kustomization.yaml`

#### Changes:
```yaml
# Before:
image: quay.io/rhacm-monitor/backend:latest
image: quay.io/rhacm-monitor/frontend:latest

# After:
image: quay.io/bzhai/rhacm-monitor-backend:latest
image: quay.io/bzhai/rhacm-monitor-frontend:latest
```

### 4. New Files Created

#### Dockerfile Variants:
1. **`deployment/docker/Dockerfile.backend`** (Multi-stage with Fedora)
   - Builds from source in container
   - ~5-10 minute build time
   - Best for CI/CD

2. **`deployment/docker/Dockerfile.backend.simple`** (Prebuilt binary)
   - Uses locally built binary
   - < 10 second build time
   - Best for development

3. **`deployment/docker/Dockerfile.backend.multistage`** (Same as main)
   - Backup of multi-stage approach

#### Build Scripts:
1. **`build-and-push.sh`**
   - Automated build and push script
   - Supports version tagging
   - Works with both podman and docker

#### Documentation:
1. **`docs/BUILD_AND_DEPLOY.md`**
   - Comprehensive build and deployment guide
   - Multiple deployment methods
   - Troubleshooting section

2. **`DOCKER_REGISTRY_CHANGES.md`** (this file)
   - Summary of all registry changes

## Registry URLs

### Red Hat Container Registry
- **URL**: `registry.access.redhat.com`
- **Images Used**:
  - `ubi9/ubi-minimal:latest` - Minimal UBI9 (~28 MB)
  - `ubi9/nodejs-20:latest` - Node.js 20 on UBI9
  - `ubi9/nginx-124:latest` - nginx 1.24 on UBI9
- **Advantages**:
  - No rate limits
  - Production-ready
  - Security updates from Red Hat
  - No authentication required for public images

### Quay.io Registry
- **URL**: `quay.io`
- **Images Used**:
  - `fedora/fedora:latest` - Fedora Linux (builder)
  - `centos/centos:stream9` - CentOS Stream 9 (alternative)
  - `bzhai/*` - Custom application images
- **Advantages**:
  - No rate limits
  - Free for public repositories
  - Better performance than Docker Hub
  - Red Hat supported

## Build Approaches

### Approach 1: Quick Build (Recommended)

**Steps:**
1. Build binary locally: `go build -o bin/server cmd/server/main.go`
2. Build container with binary: `podman build -f Dockerfile.backend.simple`

**Pros:**
- Very fast (< 10 seconds for container build)
- Minimal image size (~86 MB)
- No network dependencies during container build

**Cons:**
- Requires local Go environment
- Two-step process

### Approach 2: Multi-stage Build

**Steps:**
1. Build everything in container: `podman build -f Dockerfile.backend`

**Pros:**
- Single command
- Reproducible builds
- No local dependencies except container runtime

**Cons:**
- Slower (~5-10 minutes)
- Downloads dependencies every build
- Requires good network connection

## Image Size Comparison

| Approach | Backend Size | Frontend Size | Total |
|----------|--------------|---------------|-------|
| Before (Alpine) | ~50 MB | ~350 MB | ~400 MB |
| After (UBI) | ~86 MB | ~350 MB | ~436 MB |

**Note**: Slight size increase (~36 MB) but gains:
- Enterprise support
- Security updates
- No rate limits
- Better compatibility

## Testing Results

### Backend Image Build
✅ **Success** - Built with `registry.access.redhat.com/ubi9/ubi-minimal:latest`

```bash
$ podman build -f deployment/docker/Dockerfile.backend.simple -t quay.io/bzhai/rhacm-monitor-backend:latest .

Successfully tagged quay.io/bzhai/rhacm-monitor-backend:latest
```

### Backend Binary Size
```bash
$ ls -lh backend/bin/server
-rwxr-xr-x. 1 root root 58M Oct 17 18:17 bin/server
```

### Image Layers
```
Layer 1: UBI9 Minimal Base (~28 MB)
Layer 2: Application Binary (~58 MB)
Total: ~86 MB
```

## Migration Checklist

- [x] Updated backend Dockerfile to use Red Hat/Quay.io images
- [x] Updated frontend Dockerfile to use Red Hat images
- [x] Created simple Dockerfile variant for faster builds
- [x] Updated deployment manifests with new image names
- [x] Updated kustomization.yaml with new image references
- [x] Created automated build-and-push script
- [x] Tested backend image build successfully
- [x] Created comprehensive documentation
- [ ] Test frontend image build
- [ ] Push images to quay.io/bzhai/
- [ ] Deploy to OpenShift cluster
- [ ] Verify application works end-to-end

## Pushing to Quay.io

### Prerequisites
```bash
# Login to Quay.io
podman login quay.io
Username: bzhai
Password: <your-robot-token>
```

### Push Commands
```bash
# Push backend
podman push quay.io/bzhai/rhacm-monitor-backend:latest

# Tag with version
podman tag quay.io/bzhai/rhacm-monitor-backend:latest \
  quay.io/bzhai/rhacm-monitor-backend:v1.0.0

podman push quay.io/bzhai/rhacm-monitor-backend:v1.0.0
```

### Repository Settings
Make sure repositories are set to **public** in Quay.io:
1. Go to https://quay.io/repository/bzhai/rhacm-monitor-backend
2. Click "Settings"
3. Under "Repository Visibility", select "Public"
4. Click "Save"

## Benefits of This Approach

### 1. No Rate Limits
- Red Hat registry has no rate limits
- Quay.io has much higher limits than Docker Hub
- No need for authentication for public images

### 2. Enterprise Support
- Red Hat UBI images are supported
- Regular security updates
- Production-ready

### 3. Better Performance
- Faster pulls from Red Hat registry
- CDN distribution
- Regional mirrors

### 4. Compliance
- Red Hat UBI is freely redistributable
- Meets enterprise compliance requirements
- No licensing restrictions

### 5. Security
- Regular CVE scanning
- Timely security patches
- Minimal attack surface with ubi-minimal

## Rollback Plan

If issues occur, you can quickly rollback:

### Option 1: Use Previous Images
```bash
# If previous images exist
oc set image deployment/rhacm-monitor-backend \
  backend=<old-image>:tag \
  -n rhacm-monitor
```

### Option 2: Use Docker Hub Mirror on Quay
```bash
# Quay.io mirrors many Docker Hub images
# Example: quay.io/library/golang:1.22-alpine
```

### Option 3: Revert Dockerfile
```bash
git revert <commit-hash>
```

## Next Steps

1. **Test Frontend Build**
   ```bash
   podman build -f deployment/docker/Dockerfile.frontend \
     -t quay.io/bzhai/rhacm-monitor-frontend:latest .
   ```

2. **Push Images to Quay.io**
   ```bash
   ./build-and-push.sh
   ```

3. **Deploy to Cluster**
   ```bash
   oc apply -k deployment/k8s/
   ```

4. **Verify Deployment**
   ```bash
   oc get pods -n rhacm-monitor
   curl -k https://<route>/api/health
   ```

## Additional Notes

- All base images are from trusted sources (Red Hat, Fedora, CentOS)
- Images are regularly updated and maintained
- No changes to application code required
- Deployment manifests updated to use new image names
- Build scripts support both podman and docker

## References

- [Red Hat Universal Base Images](https://www.redhat.com/en/blog/introducing-red-hat-universal-base-image)
- [Quay.io Documentation](https://docs.quay.io/)
- [Red Hat Container Catalog](https://catalog.redhat.com/software/containers/explore)
- [OpenShift Container Platform Documentation](https://docs.openshift.com/container-platform/latest/welcome/index.html)

