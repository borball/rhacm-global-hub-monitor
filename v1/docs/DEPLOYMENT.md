# Deployment Guide

This guide provides detailed instructions for deploying the RHACM Global Hub Monitor in various environments.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Deployment Methods](#deployment-methods)
3. [Post-Installation](#post-installation)
4. [Configuration](#configuration)
5. [Upgrades](#upgrades)
6. [Uninstallation](#uninstallation)

## Prerequisites

### Cluster Requirements

- **OpenShift Version**: 4.14 or later
- **RHACM Version**: 2.9 or later
- **Architecture**: amd64
- **Resources**:
  - Minimum 2 CPU cores
  - Minimum 2 GB memory
  - Persistent storage not required

### Required Permissions

You need cluster-admin or equivalent permissions to:
- Create namespaces
- Create ClusterRoles and ClusterRoleBindings
- Create CustomResourceDefinitions (for operator method)
- Deploy applications

### Network Requirements

- Access to container registry (quay.io or internal registry)
- Outbound HTTPS access (if pulling images from external registries)
- Access to Kubernetes API server
- Access to OpenShift OAuth endpoints

## Deployment Methods

### Method 1: Operator-Based Deployment (Recommended)

The operator-based deployment provides automated lifecycle management.

#### Step 1: Install the Operator

```bash
# Login to your OpenShift cluster
oc login --token=<token> --server=https://api.cluster.example.com:6443

# Install the CRD
oc apply -f operator/config/crd/rhacmmonitor-crd.yaml

# Install operator RBAC
oc apply -f operator/config/rbac/role.yaml
oc apply -f operator/config/rbac/role-binding.yaml
oc apply -f operator/config/rbac/service-account.yaml

# Install the operator
oc apply -f operator/config/manager/operator-deployment.yaml
```

#### Step 2: Create an Instance

Create a custom resource file `my-rhacm-monitor.yaml`:

```yaml
apiVersion: apps.redhat.com/v1alpha1
kind: RHACMMonitor
metadata:
  name: rhacm-monitor
  namespace: rhacm-monitor
spec:
  backend:
    image: quay.io/rhacm-monitor/backend:v1.0.0
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
    image: quay.io/rhacm-monitor/frontend:v1.0.0
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
    host: rhacm-monitor.apps.cluster.example.com
    tlsEnabled: true
```

Apply the configuration:

```bash
oc apply -f my-rhacm-monitor.yaml
```

#### Step 3: Verify Installation

```bash
# Check operator status
oc get pods -n rhacm-monitor-operator

# Check application status
oc get rhacmmonitor -n rhacm-monitor
oc get pods -n rhacm-monitor

# Get route URL
oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}'
```

### Method 2: Manual Deployment with Kustomize

This method uses Kustomize for streamlined deployment.

#### Step 1: Customize Configuration

Edit `deployment/k8s/kustomization.yaml` to customize:
- Image tags
- Resource limits
- Replica counts
- Route hostname

#### Step 2: Deploy

```bash
# Deploy all resources
oc apply -k deployment/k8s/

# Or using kustomize CLI
kustomize build deployment/k8s/ | oc apply -f -
```

#### Step 3: Verify

```bash
oc get all -n rhacm-monitor
oc get route -n rhacm-monitor
```

### Method 3: Manual Deployment with Individual Manifests

For maximum control, deploy each component individually.

#### Step 1: Create Namespace

```bash
oc apply -f deployment/k8s/namespace.yaml
```

#### Step 2: Setup RBAC

```bash
oc apply -f deployment/k8s/serviceaccount.yaml
oc apply -f deployment/k8s/clusterrole.yaml
oc apply -f deployment/k8s/clusterrolebinding.yaml
```

#### Step 3: Deploy Backend

```bash
# Update image tag in backend-deployment.yaml if needed
oc apply -f deployment/k8s/backend-deployment.yaml
oc apply -f deployment/k8s/backend-service.yaml
```

Wait for backend to be ready:

```bash
oc wait --for=condition=available --timeout=300s deployment/rhacm-monitor-backend -n rhacm-monitor
```

#### Step 4: Deploy Frontend

```bash
# Update image tag in frontend-deployment.yaml if needed
oc apply -f deployment/k8s/frontend-deployment.yaml
oc apply -f deployment/k8s/frontend-service.yaml
```

Wait for frontend to be ready:

```bash
oc wait --for=condition=available --timeout=300s deployment/rhacm-monitor-frontend -n rhacm-monitor
```

#### Step 5: Create Route

```bash
# Update hostname in route.yaml if needed
oc apply -f deployment/k8s/route.yaml
```

## Post-Installation

### Verify Deployment

1. **Check Pod Status:**

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

2. **Check Services:**

```bash
oc get svc -n rhacm-monitor
```

3. **Check Route:**

```bash
oc get route rhacm-monitor -n rhacm-monitor
```

4. **Test API Health:**

```bash
ROUTE=$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')
curl -k https://$ROUTE/api/health
```

Expected response:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Access the Application

1. Get the route URL:

```bash
echo "https://$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')"
```

2. Open the URL in a web browser
3. Login with your OpenShift credentials

## Configuration

### Configuring Authentication

#### Enable/Disable Authentication

To disable authentication (development only):

```bash
oc set env deployment/rhacm-monitor-backend ENABLE_AUTH=false -n rhacm-monitor
```

To enable authentication:

```bash
oc set env deployment/rhacm-monitor-backend ENABLE_AUTH=true -n rhacm-monitor
```

#### Configure OAuth Client

Create an OAuthClient for better integration:

```yaml
apiVersion: oauth.openshift.io/v1
kind: OAuthClient
metadata:
  name: rhacm-monitor
redirectURIs:
- https://rhacm-monitor.apps.cluster.example.com/callback
grantMethod: auto
```

### Scaling

Scale backend:

```bash
oc scale deployment/rhacm-monitor-backend --replicas=3 -n rhacm-monitor
```

Scale frontend:

```bash
oc scale deployment/rhacm-monitor-frontend --replicas=3 -n rhacm-monitor
```

### Resource Limits

Update resource limits:

```bash
oc set resources deployment/rhacm-monitor-backend \
  --limits=cpu=2,memory=2Gi \
  --requests=cpu=1,memory=1Gi \
  -n rhacm-monitor
```

### Custom Route

To use a custom hostname:

```bash
oc patch route rhacm-monitor -n rhacm-monitor -p '{"spec":{"host":"custom.example.com"}}'
```

## Upgrades

### Upgrade Strategy

The application supports rolling updates with zero downtime.

### Upgrade via Operator

Update the image tags in the RHACMMonitor CR:

```bash
oc patch rhacmmonitor rhacm-monitor -n rhacm-monitor --type=merge -p '
{
  "spec": {
    "backend": {
      "image": "quay.io/rhacm-monitor/backend:v1.1.0"
    },
    "frontend": {
      "image": "quay.io/rhacm-monitor/frontend:v1.1.0"
    }
  }
}'
```

### Manual Upgrade

Update backend:

```bash
oc set image deployment/rhacm-monitor-backend \
  backend=quay.io/rhacm-monitor/backend:v1.1.0 \
  -n rhacm-monitor
```

Update frontend:

```bash
oc set image deployment/rhacm-monitor-frontend \
  frontend=quay.io/rhacm-monitor/frontend:v1.1.0 \
  -n rhacm-monitor
```

### Rollback

If an upgrade fails:

```bash
# Rollback backend
oc rollout undo deployment/rhacm-monitor-backend -n rhacm-monitor

# Rollback frontend
oc rollout undo deployment/rhacm-monitor-frontend -n rhacm-monitor
```

## Uninstallation

### Operator-Based Uninstallation

```bash
# Delete the instance
oc delete rhacmmonitor rhacm-monitor -n rhacm-monitor

# Delete the operator
oc delete -f operator/config/manager/operator-deployment.yaml

# Delete RBAC
oc delete -f operator/config/rbac/

# Delete CRD
oc delete -f operator/config/crd/rhacmmonitor-crd.yaml
```

### Manual Uninstallation

```bash
# Delete all resources
oc delete -k deployment/k8s/

# Or delete namespace (removes everything)
oc delete namespace rhacm-monitor

# Clean up ClusterRole and ClusterRoleBinding
oc delete clusterrole rhacm-monitor
oc delete clusterrolebinding rhacm-monitor
```

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for detailed troubleshooting steps.

