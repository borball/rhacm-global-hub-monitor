# RHACM Global Hub Monitor - v3.0.0

**Release Date:** October 21, 2025  
**Status:** ✅ In Development - First Feature Complete

## Overview

v3 builds upon the solid foundation of v2 with a major new feature: comprehensive dark/light mode support with professional color schemes.

## Features Delivered

### 1. Dark/Light Mode Toggle (COMPLETE) ✅

**Implementation:**
- CSS custom properties (variables) for all colors
- Two professionally designed themes
- Smooth 0.3s transitions
- localStorage persistence
- Toggle button in header

**Light Mode - Modern Professional:**
- Background: #f5f7fa (soft gray-blue)
- Cards: #ffffff (pure white)
- Text: #1f2937 (warm black)
- Status colors: Fresh green/amber

**Dark Mode - GitHub-Inspired:**
- Background: #0f1419 (deep dark)
- Cards: #1c2128 (charcoal)
- Text: #e6edf3 (bright white-blue)
- Status colors: Neon green/yellow for visibility

**Components Themed:**
- ✅ Homepage statistics cards (redesigned with left border accent)
- ✅ Hub cards (all fields)
- ✅ Node cards (K8s + Hardware sections)
- ✅ Spoke detail cards (all 6 stat cards)
- ✅ Policy pages (compliance cards, detail rows)
- ✅ Search/filter inputs
- ✅ Tables and rows
- ✅ Forms and buttons
- ✅ All UI elements

**User Experience:**
- Click toggle button (top-right corner)
- Instant theme switch
- Preference saved automatically
- Applies on page load

### Baseline Features (from v2)

v3 includes all v2 features:
- ✅ Performance caching (90s TTL, 350x faster)
- ✅ Console and GitOps URLs
- ✅ Refactored codebase (~200 lines eliminated)
- ✅ Aligned UI layout
- ✅ Context-aware display

## Technical Details

**Files Modified:**
- `frontend-static/styles.css` - CSS variables + dark theme (~60 new rules)
- `frontend-static/index.html` - Theme toggle button + persistence JS
- `frontend-static/app.js` - Replaced hardcoded colors with CSS classes
- Version: v=20251103

**CSS Classes Created:**
- `.k8s-section` - Kubernetes info sections
- `.hardware-section` - Hardware info sections
- `.config-badge` - Configuration version badges
- `.stat-card` - Homepage statistics
- `.compliance-card` - Policy compliance cards
- `.policy-summary-card` - Policy detail summary cards
- `.policy-message-*` - Policy status messages
- `.spoke-stat-card` - Spoke stat cards
- `.code-block` - Code/message blocks

## Testing

**All Pages Verified:**
- ✅ Homepage with statistics
- ✅ Hub list cards
- ✅ Hub details (all tabs)
- ✅ Node cards (merged K8s + Hardware)
- ✅ Spoke details
- ✅ Policy pages
- ✅ Add hub form

**Both Themes Working:**
- ✅ Light mode: Professional and clean
- ✅ Dark mode: Easy on eyes
- ✅ Toggle: Instant switching
- ✅ Persistence: Working

## Development Status

**Completed:**
- ✅ Dark/Light mode feature

**In Progress:**
- 🔜 Next v3 features (TBD)

**Version:** v3.0.0-dev  
**Frontend Version:** v=20251103

---

*v3 is actively being developed with dark mode as the first major enhancement*
