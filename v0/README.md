# RHACM Global Hub Monitor

A comprehensive web application for monitoring Red Hat Advanced Cluster Management (RHACM) global hub clusters and their managed hub and spoke clusters.

## Overview

This application provides a centralized monitoring solution for RHACM deployments in a hub-of-hubs architecture. It displays:

- **Managed Hubs**: All RHACM hub clusters managed by the global hub
- **Managed Clusters**: All spoke clusters (including SNO) managed by each hub
- **Cluster Information**: Basic cluster details, platform, version, and status
- **Nodes Information**: Node-level details including capacity, allocatable resources, and health
- **Policies Information**: Policy compliance status, violations, and remediation actions

## Architecture

The application follows a modern B/S (Browser/Server) architecture with:

- **Backend**: Go-based REST API server
  - RHACM client integration via Kubernetes API
  - OpenShift SSO authentication
  - RESTful API endpoints
  - Structured logging and error handling

- **Frontend**: React-based single-page application
  - TypeScript for type safety
  - PatternFly 5 UI components for consistent OpenShift look-and-feel
  - React Query for data fetching and caching
  - React Router for navigation

- **Operator**: Kubernetes operator for deployment and lifecycle management
  - Custom Resource Definition (CRD) for configuration
  - Automated deployment and updates
  - RBAC management

## Features

### 1. Dashboard
- Overview of all managed hubs and clusters
- Health statistics and metrics
- Quick access to recent activity

### 2. Managed Hubs View
- List all managed hub clusters
- Hub status and version information
- Number of managed clusters per hub
- Navigation to hub details

### 3. Hub Details
- Complete hub cluster information
- List of all managed spoke clusters
- Node information and health
- Policy compliance status

### 4. Cluster Details
- Detailed cluster information
- Node-level metrics and status
- Policy compliance and violations
- Console and API URLs

### 5. Authentication
- OpenShift SSO integration
- JWT token-based authentication
- RBAC support

## Prerequisites

- OpenShift 4.14+ cluster
- Red Hat Advanced Cluster Management (RHACM) 2.9+
- Global Hub configured and operational
- kubectl or oc CLI
- (Optional) Go 1.22+ for local development
- (Optional) Node.js 20+ and npm for frontend development

## Installation

### Using the Operator (Recommended)

1. **Install the CRD:**
```bash
oc apply -f operator/config/crd/rhacmmonitor-crd.yaml
```

2. **Install the Operator:**
```bash
oc apply -f operator/config/rbac/role.yaml
oc apply -f operator/config/manager/operator-deployment.yaml
```

3. **Create an instance:**
```bash
oc apply -f operator/config/samples/rhacmmonitor-sample.yaml
```

### Manual Installation

1. **Create the namespace:**
```bash
oc apply -f deployment/k8s/namespace.yaml
```

2. **Create RBAC resources:**
```bash
oc apply -f deployment/k8s/serviceaccount.yaml
oc apply -f deployment/k8s/clusterrole.yaml
oc apply -f deployment/k8s/clusterrolebinding.yaml
```

3. **Deploy backend:**
```bash
oc apply -f deployment/k8s/backend-deployment.yaml
oc apply -f deployment/k8s/backend-service.yaml
```

4. **Deploy frontend:**
```bash
oc apply -f deployment/k8s/frontend-deployment.yaml
oc apply -f deployment/k8s/frontend-service.yaml
```

5. **Create route:**
```bash
oc apply -f deployment/k8s/route.yaml
```

### Using Kustomize

```bash
oc apply -k deployment/k8s/
```

## Configuration

### Backend Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `KUBECONFIG` | Path to kubeconfig (optional for in-cluster) | - |
| `ENABLE_AUTH` | Enable OpenShift SSO authentication | `true` |
| `OAUTH_ISSUER_URL` | OAuth issuer URL | `https://kubernetes.default.svc` |
| `OAUTH_CLIENT_ID` | OAuth client ID | `rhacm-monitor` |
| `LOG_LEVEL` | Log level (info, debug, error) | `info` |

### Frontend Configuration

The frontend is configured via build-time environment variables or runtime configuration:

- API proxy is handled by nginx configuration
- Authentication tokens are stored in browser localStorage

### Operator Configuration

Customize the RHACMMonitor CR:

```yaml
apiVersion: apps.redhat.com/v1alpha1
kind: RHACMMonitor
metadata:
  name: rhacm-monitor
  namespace: rhacm-monitor
spec:
  backend:
    image: quay.io/rhacm-monitor/backend:latest
    replicas: 2
    resources:
      limits:
        cpu: "1"
        memory: 1Gi
      requests:
        cpu: 500m
        memory: 512Mi
    enableAuth: true
  frontend:
    image: quay.io/rhacm-monitor/frontend:latest
    replicas: 2
    resources:
      limits:
        cpu: 500m
        memory: 512Mi
      requests:
        cpu: 250m
        memory: 256Mi
  route:
    enabled: true
    host: rhacm-monitor.apps.example.com
    tlsEnabled: true
```

## Building from Source

### Backend

```bash
cd backend
go mod download
go build -o bin/server cmd/server/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run build
```

### Docker Images

Build backend:
```bash
docker build -f deployment/docker/Dockerfile.backend -t quay.io/rhacm-monitor/backend:latest .
docker push quay.io/rhacm-monitor/backend:latest
```

Build frontend:
```bash
docker build -f deployment/docker/Dockerfile.frontend -t quay.io/rhacm-monitor/frontend:latest .
docker push quay.io/rhacm-monitor/frontend:latest
```

## Development

### Backend Development

```bash
cd backend
export KUBECONFIG=/path/to/kubeconfig
export ENABLE_AUTH=false  # Disable auth for local dev
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080/api`

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

The UI will be available at `http://localhost:3000`

### Running Tests

Backend tests:
```bash
cd backend
make test
```

Frontend tests:
```bash
cd frontend
npm test
```

## API Documentation

### Endpoints

#### Health Endpoints

- `GET /api/health` - Health check
- `GET /api/ready` - Readiness check
- `GET /api/live` - Liveness check

#### Hub Endpoints

- `GET /api/hubs` - List all managed hubs
- `GET /api/hubs/{name}` - Get hub details
- `GET /api/hubs/{name}/clusters` - List clusters for a hub

### Authentication

All API requests (except health endpoints) require authentication via Bearer token:

```bash
curl -H "Authorization: Bearer <token>" https://rhacm-monitor.apps.example.com/api/hubs
```

## Troubleshooting

### Backend Issues

1. **Cannot connect to Kubernetes API:**
   - Verify ServiceAccount has correct RBAC permissions
   - Check if RHACM CRDs are installed

2. **Authentication failures:**
   - Verify OAuth configuration
   - Check token validity
   - Ensure ENABLE_AUTH is set correctly

### Frontend Issues

1. **API calls failing:**
   - Check nginx proxy configuration
   - Verify backend service is running
   - Check browser console for CORS errors

2. **Page not loading:**
   - Verify route is configured correctly
   - Check frontend pod logs

### Common Issues

1. **No hubs showing up:**
   - Ensure managed clusters have correct labels
   - Verify RBAC permissions to read ManagedCluster resources

2. **Empty cluster information:**
   - Check if hub kubeconfig secrets are accessible
   - Verify ACM agent is running on managed clusters

## Security Considerations

1. **Authentication**: Always enable authentication in production
2. **RBAC**: Follow principle of least privilege
3. **TLS**: Enable TLS termination on routes
4. **Network Policies**: Restrict pod-to-pod communication
5. **Image Security**: Use signed and scanned container images

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

Copyright Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Support

For issues and questions:
- Open an issue in the GitHub repository
- Contact the Red Hat support team

## Roadmap

- [ ] Real-time metrics and monitoring
- [ ] Alert management integration
- [ ] Multi-tenancy support
- [ ] Advanced filtering and search
- [ ] Export and reporting features
- [ ] Integration with Observability operators
- [ ] Grafana dashboards
- [ ] Policy template management

