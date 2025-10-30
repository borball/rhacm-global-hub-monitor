# RHACM Global Hub Monitor

**One-stop monitoring solution for Red Hat Advanced Cluster Management (RHACM) multi-cluster environments.**

## Quick Start

```bash
# Clone repository
git clone https://github.com/borball/rhacm-global-hub-monitor.git
cd rhacm-global-hub-monitor/v3/deployment

# Deploy
./deploy.sh

# Get URL
oc get route hubs -n rhacm-monitor -o jsonpath='{.spec.host}'
```

Access at: `https://<route-url>`

## Features

### v3.0.0 (Latest)

- 🎨 **Dark/Light Mode** - Professional themes with toggle
- 🔧 **Operators Monitoring** - View installed operators, lazy loading for spokes
- ⚡ **Performance** - 300x faster with caching, session affinity
- 🔄 **Per-Hub Refresh** - Granular control with refresh buttons
- 📦 **Simplified Deployment** - One command setup
- 📱 **Instant Navigation** - Client-side caching for smooth UX

## Documentation

- **v3/README.md** - Quick start guide
- **v3/DEPLOYMENT.md** - Complete deployment instructions
- **v3/VERSION.md** - Feature documentation
- **v3/OPERATORS_SETUP.md** - Operators configuration
- **v3/KNOWN_ISSUES.md** - Known limitations
- **SPRINT_HISTORY.md** - Complete development story

## Requirements

- OpenShift 4.x or Kubernetes 1.24+
- RHACM installed
- Hub kubeconfig secrets

## What's Included

**For Managed Hubs (Auto-discovered):**
- Full cluster information
- Console and GitOps URLs
- Operators monitoring (hub + spokes)
- Policy compliance tracking
- Spoke cluster details

**For Unmanaged Hubs (Manually added):**
- Basic monitoring
- Stored in rhacm-monitor namespace
- Add via UI

## Version History

- **v0** - Foundation (baseline monitoring)
- **v1** - Hub management + Policy enforcement
- **v2** - Performance (350x improvement with caching)
- **v3** - Modern UX (dark mode, operators, simplified deployment)
- **v4** - Planned (authentication improvements, metrics)

## Support

**Repository:** https://github.com/borball/rhacm-global-hub-monitor  
**Issues:** Check logs and documentation  
**License:** Apache 2.0

---

*Production-ready RHACM monitoring for multi-cluster environments.*
