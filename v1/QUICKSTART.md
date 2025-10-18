# Quick Start Guide

Get the RHACM Global Hub Monitor up and running in 5 minutes!

## Prerequisites Checklist

- [ ] OpenShift 4.14+ cluster access
- [ ] RHACM 2.9+ installed
- [ ] Global Hub configured
- [ ] `oc` CLI installed and logged in

## Installation in 3 Steps

### Step 1: Create Namespace and RBAC

```bash
oc apply -f deployment/k8s/namespace.yaml
oc apply -f deployment/k8s/serviceaccount.yaml
oc apply -f deployment/k8s/clusterrole.yaml
oc apply -f deployment/k8s/clusterrolebinding.yaml
```

### Step 2: Deploy Application

```bash
# Deploy backend
oc apply -f deployment/k8s/backend-deployment.yaml
oc apply -f deployment/k8s/backend-service.yaml

# Deploy frontend
oc apply -f deployment/k8s/frontend-deployment.yaml
oc apply -f deployment/k8s/frontend-service.yaml

# Create route
oc apply -f deployment/k8s/route.yaml
```

### Step 3: Access the UI

```bash
# Get the route URL
ROUTE_URL=$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')
echo "Access the application at: https://$ROUTE_URL"

# Wait for pods to be ready
oc wait --for=condition=available --timeout=300s deployment/rhacm-monitor-backend -n rhacm-monitor
oc wait --for=condition=available --timeout=300s deployment/rhacm-monitor-frontend -n rhacm-monitor
```

Open the URL in your browser and log in with your OpenShift credentials!

## One-Line Installation (Alternative)

If you have Kustomize:

```bash
oc apply -k deployment/k8s/
```

## Verify Installation

Check that everything is running:

```bash
# Check pods
oc get pods -n rhacm-monitor

# Expected output:
# NAME                                      READY   STATUS    RESTARTS   AGE
# rhacm-monitor-backend-xxxxx-xxxxx        1/1     Running   0          2m
# rhacm-monitor-backend-xxxxx-xxxxx        1/1     Running   0          2m
# rhacm-monitor-frontend-xxxxx-xxxxx       1/1     Running   0          2m
# rhacm-monitor-frontend-xxxxx-xxxxx       1/1     Running   0          2m

# Test the API
ROUTE_URL=$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')
curl -k https://$ROUTE_URL/api/health

# Expected output:
# {"status":"healthy","version":"1.0.0","timestamp":"..."}
```

## First Steps

1. **Navigate to Dashboard**: View overall statistics of your hubs and clusters
2. **Browse Managed Hubs**: Click on "Managed Hubs" to see all your hub clusters
3. **View Hub Details**: Click on any hub to see its managed clusters, nodes, and policies
4. **Explore Clusters**: Click on any spoke cluster to see detailed information

## Configuration Options

### Disable Authentication (Development Only)

```bash
oc set env deployment/rhacm-monitor-backend ENABLE_AUTH=false -n rhacm-monitor
```

### Adjust Replicas

```bash
# Scale backend
oc scale deployment/rhacm-monitor-backend --replicas=3 -n rhacm-monitor

# Scale frontend
oc scale deployment/rhacm-monitor-frontend --replicas=3 -n rhacm-monitor
```

### Custom Route Hostname

```bash
oc patch route rhacm-monitor -n rhacm-monitor -p '{"spec":{"host":"my-custom-host.apps.example.com"}}'
```

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
oc get pods -n rhacm-monitor

# Check pod logs
oc logs -l component=backend -n rhacm-monitor
oc logs -l component=frontend -n rhacm-monitor

# Check events
oc get events -n rhacm-monitor --sort-by='.lastTimestamp'
```

### Cannot Access UI

```bash
# Verify route exists
oc get route rhacm-monitor -n rhacm-monitor

# Check if route is accessible
ROUTE_URL=$(oc get route rhacm-monitor -n rhacm-monitor -o jsonpath='{.spec.host}')
curl -k -I https://$ROUTE_URL
```

### No Hubs Showing

```bash
# Verify RBAC permissions
oc auth can-i get managedclusters --as=system:serviceaccount:rhacm-monitor:rhacm-monitor

# Check backend logs for errors
oc logs -l component=backend -n rhacm-monitor | grep -i error

# Verify ManagedCluster resources exist
oc get managedclusters
```

## Uninstall

Remove the application:

```bash
# Quick uninstall
oc delete namespace rhacm-monitor

# Clean up cluster-scoped resources
oc delete clusterrole rhacm-monitor
oc delete clusterrolebinding rhacm-monitor
```

## Next Steps

- Read the [Full Documentation](README.md)
- Check the [API Documentation](docs/API.md)
- See [Deployment Guide](docs/DEPLOYMENT.md) for advanced options
- Review [Development Guide](docs/DEVELOPMENT.md) to contribute

## Getting Help

- Check the logs: `oc logs -l app=rhacm-monitor -n rhacm-monitor`
- Review documentation in the `docs/` directory
- Open an issue on GitHub
- Contact Red Hat support

## Common Use Cases

### View All Managed Hubs
Navigate to **Managed Hubs** page to see all hub clusters with their status and cluster counts.

### Monitor Spoke Cluster Health
Click on a hub, then navigate to its managed clusters tab to see all spoke clusters and their health status.

### Check Policy Compliance
Click on any hub or cluster to view the **Policies** tab and see compliance status.

### View Node Resources
Navigate to the **Nodes** tab in any hub or cluster detail page to see node capacity and allocatable resources.

---

**That's it!** You now have a fully functional RHACM Global Hub Monitor running on your cluster. 🎉

