# RHACM Global Hub Monitor - v3.0.0

**Release Date:** October 21, 2025  
**Status:** ✅ In Development - Two Major Features Complete

## Features Delivered

### 1. Dark/Light Mode Toggle ✅ COMPLETE

Professional theme system with smooth transitions.

**Implementation:**
- CSS custom properties for all colors
- Two professionally designed themes  
- Smooth 0.3s transitions
- localStorage persistence
- Toggle button in header

**Themes:**

**Light Mode:**
- Background: #f5f7fa (soft gray-blue)
- Cards: #ffffff (pure white)
- Text: #1f2937 (warm black)

**Dark Mode (GitHub-Inspired):**
- Background: #0f1419 (deep dark)
- Cards: #1c2128 (charcoal)
- Text: #e6edf3 (bright white-blue)

**All Components Themed:**
- Homepage statistics, hub cards, node cards
- Spoke details, policy pages, forms
- K8s sections (blue), Hardware sections (amber)
- Tables, badges, status indicators

### 2. Operators Tab ✅ COMPLETE

Comprehensive operator monitoring with lazy loading for performance.

**Hub Operators:**
- Full operators tab with table
- ClusterServiceVersion (CSV) fetching
- Smart grouping by operator name
- Example: 304 installations → 45 unique operators
- Search/filter functionality
- Columns: Name, Version, Namespaces, Status, Provider

**Spoke Cluster Operators:**
- **Lazy Loading** for performance (scales to 1000+ spokes)
- Operators column in spoke table (shows "...")
- Operators loaded on-demand when details expanded
- Stat card updates after loading
- Compact table in expansion (top 10)
- Grouped by name

**Performance:**
- Initial page load: FAST (no spoke operator fetching)
- Spoke expansion: Operators fetched on-demand
- Endpoint: GET /api/hubs/:name/spokes/:spoke/operators
- Scales to large deployments

**Backend:**
- `backend/pkg/client/operators.go` - CSV fetching
- `backend/pkg/handlers/spokes.go` - Lazy loading endpoint
- `backend/pkg/client/hubclient.go` - NewSpokeClientFromKubeconfig
- Integrated into enrichHubWithRemoteData()

**Frontend:**
- `app.js` - Lazy loading logic, operator grouping
- Shows "Loading..." during fetch
- Updates stat card + table after load
- Grouped display with namespace counts

## Baseline Features (from v2)

v3 includes all v2 features:
- ✅ Performance caching (90s TTL, 350x faster)
- ✅ Console and GitOps URLs
- ✅ Refactored codebase
- ✅ Aligned UI layout
- ✅ Context-aware display

## Technical Details

**Files Created/Modified:**
- `backend/pkg/client/operators.go` (NEW)
- `backend/pkg/handlers/spokes.go` (NEW)
- `v3/OPERATORS_SETUP.md` (NEW - documentation)
- `backend/pkg/models/types.go` - OperatorInfo model
- `backend/pkg/client/rhacm.go` - Operator integration
- `backend/pkg/client/hubclient.go` - Spoke client
- `backend/pkg/api/router.go` - Lazy endpoint
- `frontend-static/styles.css` - Dark mode variables
- `frontend-static/index.html` - Theme toggle
- `frontend-static/app.js` - Operators tab, lazy loading

**CSS Classes Added:**
- Theme variables (--bg-*, --text-*, --badge-*)
- `.stat-card` - Homepage statistics
- `.k8s-section`, `.hardware-section`
- `.config-badge`, `.policy-summary-card`
- `.code-block`, `.spoke-stat-card`

## Testing

**All Features Verified:**
- ✅ Dark/light mode toggle working
- ✅ All pages themed correctly
- ✅ Hub operators: 45 unique displayed
- ✅ Spoke lazy loading: Working
- ✅ Operators fetched on expansion
- ✅ Performance: Fast initial load

**Performance Metrics:**
- Initial hub load: ~1.8s
- Lazy operator fetch: < 1s per spoke
- Scales to 1000+ spokes

## Deployment

**Application:** https://hubs-rhacm-monitor.apps.vhub.outbound.vz.bos2.lab  
**Images:** quay.io/bzhai/rhacm-monitor-backend:v3  
**Frontend:** v=20251104

## Known Limitations

- Spoke operators require kubeconfig secrets on hub
- Without secrets, spoke operators show empty
- Hub operators work without additional setup

See `OPERATORS_SETUP.md` for configuration details.

## Next Steps

v3 continues to be enhanced with additional features as requirements emerge.

**Current Status:** Two major features complete  
**Version:** v3.0.0-dev

---

*v3 is actively being developed with significant enhancements over v2*
