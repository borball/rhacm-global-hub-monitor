# RHACM Global Hub Monitor - v1

## Version Information

**Version**: v1.0  
**Release Date**: October 18, 2025  
**Based On**: v0 (October 17-18, 2025)  
**Status**: ✅ Production Ready

## New Features in v1

### 1. Policy Enforcement via TALM ⚡
- **CGU (ClusterGroupUpgrade) Creation**: One-click policy enforcement
- **Enforce Button**: Appears on non-compliant policies
- **Auto-configuration**: 
  - Namespace: Always `ztp-install`
  - Short CGU names: `{cluster}-{6-digit-timestamp}` (under 63 chars)
  - Removes `ztp-vdu.` prefix from policy names
  - Confirmation dialog before creation
- **TALM Integration**: Creates CGU resources for automated remediation

### 2. Policy Status Messages 📋
- **Latest Status Display**: Shows most recent policy status message
- **Timestamp-based**: Correctly identifies latest message from history
- **Detailed Information**: Displays violations, notifications, and compliance details
- **Troubleshooting**: Helps identify exactly what's wrong with non-compliant policies

### 3. Configuration Version Tracking ⚙️
- **Configuration Display**: Shows on all clusters (hubs and spokes)
- **Extracted from**: ManagedCluster label `configuration-version`
- **Examples**: `hub-418-v1`, `hub-418-v2`, `vdu2-4.18-p3a5`
- **Visible in**:
  - Hub cards (dashboard/list)
  - Hub overview page
  - Spoke cluster table (new column)
  - Spoke cluster detail cards
- **Searchable**: New configuration filter on spoke clusters page

### 4. Improved UI/UX 🎨
- **Redesigned Policy Details**:
  - 4 summary cards at top (Namespace, Compliance, Remediation, Violations)
  - Full-width Latest Status section (prominent, color-coded)
  - Removed Standards/Categories/Controls (cleaner)
  - Better visual hierarchy
- **Compact Spoke Detail Page**:
  - 6 info cards in single row (was 2 columns)
  - Compact hardware grid (4 columns)
  - Inline filter controls
  - 60% less vertical space
- **Enhanced Spoke Table**:
  - Added Configuration column
  - 3-field search (name, version, configuration)
  - Better space efficiency

### 5. Correct Data Display 🔧
- **Violation Counting**: Only counts actual non-compliant details (not all details)
- **Latest Message**: Timestamp-based detection (not array index)
- **Node Count**: Shows merged count (3 not 6) everywhere
- **Cluster-prefixed Filenames**: Policy YAML downloads include cluster name

## Improvements from v0

### Documentation
- ✅ Clean root directory (4 essential files)
- ✅ Organized docs/ structure (guides, project-status, test-results)
- ✅ Professional folder organization
- ✅ Created STRUCTURE.md and docs/README.md

### Performance
- ✅ Same excellent performance (< 200ms API responses)
- ✅ Additional endpoint for CGU creation

### Features
- ✅ All v0 features retained
- ✅ Plus 5 major new features listed above

## Bug Fixes

- Fixed violation counting (was counting all details, now only non-compliant)
- Fixed latest message detection (now timestamp-based, not array order)
- Fixed node count display (shows 3 not 6 everywhere)
- Fixed CGU name length (now under 63 chars for TALM compatibility)

## API Changes

### New Endpoints
```
POST /api/cgu/create              Create ClusterGroupUpgrade
GET  /api/policies/{ns}/{name}/yaml?hub={hub}  Download policy with hub parameter
```

### Enhanced Features
- Policy YAML endpoint supports hub parameter for spoke policies
- Configuration version extracted from cluster labels

## Deployment

**Current Deployment**:
- URL: https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab
- Backend: quay.io/bzhai/rhacm-monitor-backend:latest (v1)
- Status: ✅ Running

## Technical Details

### Backend Changes
- Added `cgu.go` handler for ClusterGroupUpgrade creation
- Enhanced policy client to extract status messages
- Configuration version extraction from labels
- Improved violation counting logic

### Frontend Changes
- New compact layouts for spoke details
- CGU enforce button with confirmation
- Configuration column and filters
- Improved policy details layout
- Enhanced search/filter capabilities

---

**v1 is production-ready with advanced policy management, configuration tracking, and TALM integration!**


