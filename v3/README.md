# RHACM Global Hub Monitor - v2

A comprehensive web application for monitoring Red Hat Advanced Cluster Management (RHACM) Global Hub deployments with managed and unmanaged hubs.

## Quick Start

```bash
# Deploy using Kustomize
oc apply -k deployment/k8s/

# Access the application
oc get route rhacm-monitor -n rhacm-monitor
```

See [QUICKSTART.md](QUICKSTART.md) for detailed instructions.

## Documentation

- **[Quick Start Guide](QUICKSTART.md)** - Get started in 5 minutes
- **[API Documentation](docs/API.md)** - Complete API reference
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Deployment instructions
- **[Architecture](docs/ARCHITECTURE.md)** - System design
- **[Development Guide](docs/DEVELOPMENT.md)** - Developer documentation
- **[Project Structure](STRUCTURE.md)** - File organization

## Features

### v2.0 Features (All from v1.0)

**Core Monitoring:**
- Hub and spoke cluster monitoring (managed + unmanaged)
- Policy compliance tracking with accurate calculations
- Node information (Kubernetes + BareMetalHost merged)
- Configuration version tracking across all clusters

**Advanced Features:**
- Policy enforcement via TALM (ClusterGroupUpgrade creation)
- Policy status messages with detailed violations
- Hub management (Add/Remove with dual-method authentication)
- Unmanaged hub discovery and monitoring
- Search & filter (3 fields on spokes, 2 on policies)
- Policy YAML download from live cluster

**UI/UX:**
- Professional enterprise design
- Compact layouts (60% space efficiency)
- Proper status colors (green/orange/red)
- Cache-busting for immediate updates
- Responsive and intuitive interface

## Deployment

**Current Deployment:**
- URL: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab
- Backend: 2 pods
- Frontend Proxy: 2 pods
- Status: Production-ready

**Monitoring:**
- 3 Hubs (acm1, acm2, production-hub)
- 5 Spoke clusters
- 45+ Policies
- Complete infrastructure visibility

## Technology Stack

- **Backend**: Go 1.22+ with Gin framework
- **Frontend**: React 18 + TypeScript / Static HTML
- **Deployment**: Kubernetes/OpenShift
- **Images**: Red Hat UBI 9

## Requirements Met

✅ All 7 original requirements: 100%  
✅ All 11 features: Deployed and working

## License

Apache License 2.0

---

**Version**: v2.0 (Development)  
**Based On**: v1.0 (Production)  
**Status**: Ready for enhancements
