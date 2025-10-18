# RHACM Global Hub Monitor - v1

A comprehensive web application for monitoring Red Hat Advanced Cluster Management (RHACM) Global Hub deployments with managed hubs and spoke clusters.

## Demo

📹 **[Watch Application Demo](screen-capture.webm)** - See the application in action

## Quick Links

- **[Quick Start Guide](QUICKSTART.md)** - Get started in 5 minutes
- **[API Documentation](docs/API.md)** - Complete API reference
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Detailed deployment instructions
- **[Architecture](docs/ARCHITECTURE.md)** - System design and architecture

## Overview

The RHACM Global Hub Monitor provides complete visibility into your RHACM infrastructure:

- **2-Level Hub Monitoring** - Global Hub → Managed Hubs → Spoke Clusters
- **Hardware Inventory** - Complete node information (Kubernetes + BareMetalHost)
- **Policy Compliance** - Track policies across all clusters with ZTP wave ordering
- **Scalable UI** - Handles 500+ spoke clusters with search/filter
- **Policy Export** - Download policies as YAML from live cluster

## Features

### Monitoring Capabilities

✅ **Hub Monitoring**
- List all managed hubs
- Cluster information (versions, platform, console URLs)
- Node details (merged Kubernetes + BareMetalHost data)
- Policy compliance with NIST SP 800-53 tracking

✅ **Spoke Cluster Monitoring**
- Scalable table view for hundreds of clusters
- Search by cluster name and version
- Complete cluster information
- Hardware inventory with BMC, CPU, RAM, Storage
- Policy compliance tracking

✅ **Policy Management**
- Sorted by ZTP deployment wave
- Search by name, filter by compliance status
- Expandable details view
- Download as YAML from live cluster
- Cluster-prefixed filenames

✅ **Node Information**
- Kubernetes Node data (status, roles, versions)
- BareMetalHost data (CPU, RAM, Storage, BMC)
- Merged display (same physical node shows both)
- Complete hardware inventory

### UI Features

- Professional dashboard with statistics
- Tabbed interface for organized navigation
- Real-time search and filtering
- Expandable detail views
- Policy and cluster export capabilities
- Responsive design

## Quick Start

```bash
# Deploy using Kustomize
oc apply -k deployment/k8s/

# Access the application
oc get route rhacm-monitor -n rhacm-monitor
```

See [QUICKSTART.md](QUICKSTART.md) for detailed instructions.

## Documentation

### Getting Started
- [Quick Start Guide](QUICKSTART.md) - 5-minute setup
- [Deployment Guide](docs/DEPLOYMENT.md) - Complete deployment instructions
- [Build Guide](docs/BUILD_AND_DEPLOY.md) - Build and deploy from source

### Reference
- [API Documentation](docs/API.md) - API endpoints and examples
- [Architecture](docs/ARCHITECTURE.md) - System design
- [Development Guide](docs/DEVELOPMENT.md) - Developer documentation

### Additional Resources
- [API Examples](docs/guides/API_EXAMPLES.md) - Real API usage examples
- [Access Instructions](docs/guides/ACCESS_INSTRUCTIONS.md) - How to use the application
- [Docker Registry Info](docs/guides/DOCKER_REGISTRY_CHANGES.md) - Container image details

### Project Status
- [Project Summary](docs/project-status/PROJECT_SUMMARY.md) - Complete overview
- [Success Report](docs/project-status/SUCCESS_REPORT.md) - Test results
- [Test Results](docs/test-results/) - All validation results

## Technology Stack

### Backend
- **Language**: Go 1.22+
- **Framework**: Gin (HTTP)
- **Kubernetes**: client-go v0.30.0
- **RHACM**: open-cluster-management.io/api v0.14.0

### Frontend
- **Framework**: React 18 + TypeScript (production)
- **Static Version**: HTML5 + CSS3 + Vanilla JS (deployed)
- **UI Components**: PatternFly 5 design system
- **Build Tool**: Vite 5

### Infrastructure
- **Container Images**: Red Hat UBI 9
- **Deployment**: Kubernetes/OpenShift 4.14+
- **Authentication**: OpenShift OAuth (JWT)

## Deployment

**Current Deployment**:
- **URL**: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab
- **Cluster**: vhub.outbound.vz.bos2.lab (Global Hub)
- **Status**: ✅ Running and operational

**Components**:
- Backend: 2 pods (quay.io/bzhai/rhacm-monitor-backend:latest)
- Frontend Proxy: 2 pods (httpd with API proxy)
- Namespace: rhacm-monitor

## Requirements Met

All 7 original requirements are fully implemented:

1. ✅ Web application with B/S architecture
2. ✅ Latest web technology frontend
3. ✅ Golang backend with best practices
4. ✅ Reasonable test coverage
5. ✅ OpenShift operator installation
6. ✅ OpenShift SSO authentication
7. ✅ Complete monitoring for hubs and spokes

## Support

For issues, questions, or contributions, see:
- [Development Guide](docs/DEVELOPMENT.md)
- [Project Documentation](docs/)
- Version History in [VERSION_HISTORY.md](../VERSION_HISTORY.md)

---

**Version**: v1  
**Status**: Ready for development  
**Based on**: v0 (Production-ready baseline)  
**License**: Apache 2.0
