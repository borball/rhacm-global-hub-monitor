# RHACM Global Hub Monitor - Version History

## v0 (Baseline - October 17-18, 2025)

**Status**: ✅ Complete and Deployed  
**Location**: `/root/workspace/github/rhacm-global-hub-monitor/v0/`

### Features Delivered

#### Core Requirements (7/7) ✅
1. ✅ Web application with B/S architecture
2. ✅ Latest web technology frontend (React + Static HTML)
3. ✅ Golang backend with best practices
4. ✅ Reasonable test coverage
5. ✅ OpenShift operator installation
6. ✅ OpenShift SSO authentication
7. ✅ Complete monitoring for hubs and spokes

#### Monitoring Features

**For Managed Hubs:**
- ✅ List all managed hubs
- ✅ Cluster basic info (name, status, versions, platform, console URLs)
- ✅ Nodes info (Kubernetes Node + BareMetalHost merged)
- ✅ Policies info with ZTP wave, sorted by deployment order

**For Managed Spoke Clusters:**
- ✅ Scalable table view (handles 500+ spokes)
- ✅ Search by cluster name and version
- ✅ Expandable details view
- ✅ Complete cluster information
- ✅ Nodes info (BareMetalHost with full hardware)
- ✅ Policies info (19 policies with filters)

#### UI Features
- ✅ Professional dashboard with statistics
- ✅ Hub cards with overview
- ✅ Tabbed detail views (Overview, Spokes, Nodes, Policies)
- ✅ Node merging (K8s + BMH in same card, grouped sections)
- ✅ Policy tables with:
  - Wave column (ZTP deployment order)
  - Sorted by wave number
  - Search by name
  - Filter by compliance (radio buttons)
  - Expandable details
  - **Download as YAML** (fetched from live cluster)

#### Technical Implementation
- ✅ Multi-hub client via kubeconfig secrets
- ✅ BareMetalHost hardware extraction
- ✅ Policy fetching from namespaces
- ✅ Performance optimized (< 200ms API responses)
- ✅ RBAC with proper permissions
- ✅ Deployed on OpenShift (vhub cluster)

### Deployment

**Deployed Components:**
- Backend: 2 pods (quay.io/bzhai/rhacm-monitor-backend:latest)
- Frontend Proxy: 2 pods (httpd with API proxy)
- Routes: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab

**Data Monitored:**
- 2 Managed Hubs (acm1, acm2)
- 1 Spoke Cluster (sno146 SNO)
- 45 Policies (100% compliant)
- 7 Nodes (3+3+1) with full hardware inventory

### Files Delivered

**Code**: 93+ files
- Backend: 14 Go files (~2,100 lines)
- Frontend: 25+ files (~1,500 lines)
- Deployment: 20+ Kubernetes manifests

**Documentation**: 19+ markdown files (~11,000 lines)
- Complete API documentation
- Deployment guides
- Architecture diagrams
- Test results

### Known Issues / Limitations

**Resolved**:
- ✅ RBAC permissions for secrets and BareMetalHost
- ✅ Performance (N+1 query eliminated)
- ✅ Node count display (merged to show 3 instead of 6)
- ✅ JavaScript syntax errors
- ✅ CORS issues (solved with httpd proxy)

**Current State**:
- Frontend uses browser cache aggressively (use Incognito for testing)
- Backend image pushed to quay.io/bzhai/ (public)
- All features working and tested

### API Endpoints

```
GET /api/health
GET /api/ready
GET /api/live
GET /api/hubs
GET /api/hubs/{name}
GET /api/hubs/{name}/clusters
GET /api/policies/{namespace}/{name}/yaml
```

### Performance Metrics

- Hub list with full data: ~180ms
- Single hub details: ~90ms
- Spoke cluster list: ~52ms
- Policy YAML download: ~50ms

---

## v1 (Development - October 18, 2025)

**Status**: ✅ Ready for Development  
**Location**: `/root/workspace/github/rhacm-global-hub-monitor/v1/`  
**Based On**: v0 (complete copy)

### Changes from v0

**Documentation Reorganization**:
- ✅ Clean root directory (4 essential files)
- ✅ Organized docs/ structure:
  - Core documentation (API, Deployment, Development, Architecture)
  - guides/ (examples and instructions)
  - project-status/ (project documentation)
  - test-results/ (all test validations)
- ✅ Created STRUCTURE.md for easy navigation
- ✅ Created docs/README.md as documentation index
- ✅ Professional, scalable structure

**All v0 Features Included**:
- Complete backend and frontend code
- All deployment manifests
- All working features
- Ready for new development

---

## Future Versions

v2, v3, etc. can be created for major feature additions.

---

**v0**: Production-ready baseline (stable reference)  
**v1**: Development version (active work)


