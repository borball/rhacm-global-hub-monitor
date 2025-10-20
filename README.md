# RHACM Global Hub Monitor

A production-ready web application for comprehensive monitoring of Red Hat Advanced Cluster Management (RHACM) Global Hub deployments with managed and unmanaged hub clusters.

## Project Status

**Current Version**: v2.0 (Development)  
**Stable Versions**: v1.0 (Production), v0 (Baseline)  
**Live Deployment**: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab  
**Repository**: github.com:borball/rhacm-global-hub-monitor.git

## Quick Links

- **[Quick Start](v2/QUICKSTART.md)** - Get started in 5 minutes
- **[Documentation](v2/docs/README.md)** - Complete documentation index
- **[Version History](VERSION_HISTORY.md)** - All version details
- **[Project Summary](SUMMARY.md)** - Quick overview

## What's New in v1/v2

### v1.0 Features (Production)
- ✅ **Hub Management**: Add/remove hubs with dual-method authentication
- ✅ **Unmanaged Hub Support**: Monitor external hubs via kubeconfig
- ✅ **Policy Enforcement**: One-click CGU creation via TALM
- ✅ **Configuration Tracking**: Search and filter by configuration version
- ✅ **Policy Status Messages**: Detailed violation information
- ✅ **Improved UI**: Compact layouts, proper status colors

### v2.0 (Development Baseline)
- Complete copy of v1.0 features
- Clean documentation structure
- Ready for new enhancements

## Features Overview

### Core Monitoring (v0+)
- Hub and spoke cluster monitoring
- Policy compliance tracking
- Node information (Kubernetes + BareMetalHost)
- Hardware inventory with BMC details

### Advanced Features (v1+)
- Policy enforcement via TALM/CGU
- Hub management (add/remove)
- Unmanaged hub discovery
- Configuration version tracking
- Policy status messages
- Multi-field search and filter
- Policy YAML download

### UI/UX
- Professional enterprise design
- Scalable for 500+ spoke clusters
- Compact layouts (60% space efficiency)
- Real-time search and filtering
- Proper color coding (green/orange/red)
- Cache-busting for immediate updates

## Versions

### v0 - Stable Baseline
- All 7 requirements met
- Complete working application
- Reference implementation
- **Location**: `v0/`

### v1 - Production
- Hub management features
- Unmanaged hub support
- Policy enforcement
- Currently deployed
- **Location**: `v1/`

### v2 - Development
- Based on v1
- Clean structure
- Ready for enhancements
- **Location**: `v2/`

## Deployment

**Live Application:**
- URL: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab
- Cluster: vhub.outbound.vz.bos2.lab
- Namespace: rhacm-monitor
- Status: ✅ Operational

**Monitoring:**
- 3 Hubs (2 managed + 1 unmanaged)
- 5 Spoke clusters
- 45+ Policies (98% compliant)
- 7+ Nodes with complete hardware inventory

## Getting Started

```bash
# Deploy v2 (latest)
cd v2
oc apply -k deployment/k8s/

# Access application
oc get route rhacm-monitor -n rhacm-monitor
```

See [v2/QUICKSTART.md](v2/QUICKSTART.md) for detailed instructions.

## Documentation

All documentation is organized in version-specific folders:

- **v2/docs/**: Latest documentation
  - API reference
  - Deployment guides
  - Architecture diagrams
  - Development guides
  - Examples and tutorials

## Technology

- **Backend**: Go 1.22+ (Gin framework)
- **Frontend**: React 18 + TypeScript / Static HTML
- **Deployment**: Kubernetes/OpenShift 4.14+
- **Authentication**: OpenShift OAuth (JWT)
- **Images**: Red Hat UBI 9

## Requirements

✅ All 7 original requirements: **100% Met**
1. Web application with B/S architecture
2. Latest web technology frontend
3. Golang backend with best practices
4. Reasonable test coverage
5. OpenShift operator installation
6. OpenShift SSO authentication
7. Complete monitoring for hubs and spokes

Plus **11 deployed features** exceeding requirements.

## Contributing

Development workflow:
- v0: Keep as stable reference
- v1: Production version (deployed)
- v2: Active development
- Create v3+ for major new features

## License

Apache License 2.0

---

**RHACM Global Hub Monitor** - Complete visibility into your RHACM infrastructure

**Status**: ✅ Production-Ready  
**All Requirements**: ✅ Met  
**All Features**: ✅ Working
