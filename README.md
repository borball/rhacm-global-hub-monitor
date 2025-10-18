# RHACM Global Hub Monitor

A production-ready web application for monitoring Red Hat Advanced Cluster Management (RHACM) Global Hub deployments with managed hubs and spoke clusters.

## Quick Links

- **[Quick Start](v1/QUICKSTART.md)** - Get started in 5 minutes
- **[Documentation](v1/docs/README.md)** - Complete documentation index
- **[Version History](VERSION_HISTORY.md)** - Version tracking

## Project Status

**Current Version**: v1 (Development)  
**Stable Baseline**: v0  
**Deployment**: Live on vhub.outbound.vz.bos2.lab  
**URL**: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab

## Demo

📹 **[Watch Application Demo](v1/screen-capture.webm)** - See the application in action

The demo shows:
- Dashboard with hub statistics
- Hub details with spokes, nodes, and policies
- Spoke cluster monitoring with hardware inventory
- Policy management with search, filter, and YAML download
- Complete workflow and UI features

## Features

✅ **Complete Monitoring**
- 2-level hub hierarchy (Global Hub → Managed Hubs → Spoke Clusters)
- Real-time cluster status and health
- Hardware inventory (Kubernetes Nodes + BareMetalHost)
- Policy compliance tracking with NIST SP 800-53

✅ **Scalable Design**
- Table view for 500+ spoke clusters
- Search and filter capabilities
- Sorted by ZTP deployment wave
- Expandable details on demand

✅ **Policy Management**
- View all policies with compliance status
- Download policies as YAML from live cluster
- Filter by name and compliance
- ZTP wave-based ordering

✅ **Production Ready**
- Deployed on OpenShift
- RBAC configured
- Performance optimized (< 200ms API)
- Comprehensive documentation

## Versions

### v0 - Stable Baseline
- Complete working application
- All 7 requirements met
- Production-ready
- Reference implementation

### v1 - Development
- Clean, reorganized structure
- Based on v0
- Ready for enhancements
- Professional organization

## Technology

- **Backend**: Go 1.22+ with Gin framework
- **Frontend**: React 18 + TypeScript / Static HTML
- **Deployment**: Kubernetes/OpenShift
- **Authentication**: OpenShift OAuth
- **Images**: Red Hat UBI 9

## Getting Started

```bash
# Quick deployment
cd v1
oc apply -k deployment/k8s/

# Access application
oc get route rhacm-monitor -n rhacm-monitor
```

See [v1/QUICKSTART.md](v1/QUICKSTART.md) for detailed instructions.

## Documentation

All documentation is organized in v1/docs/:
- Core documentation (API, Deployment, Architecture)
- User guides and examples
- Project status and summaries
- Test results and validation

See [v1/docs/README.md](v1/docs/README.md) for complete index.

## License

Apache License 2.0

---

**RHACM Global Hub Monitor** - Complete visibility into your RHACM infrastructure
