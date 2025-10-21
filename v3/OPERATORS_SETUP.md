# Operators Feature Setup Guide

## Overview

The RHACM Global Hub Monitor v3 includes comprehensive operator monitoring for both hubs and spoke clusters.

## Current Status

### Working ✅
- **Hub Operators:** Fully functional
  - Fetches ClusterServiceVersion (CSV) resources
  - Shows all installed operators with versions
  - Example: acm1 shows ~45 unique operators (304 total installations)
  - Grouped by operator name
  - Multiple namespace installations combined

### Limited ⚠️
- **Spoke Cluster Operators:** Requires setup
  - Code is ready to fetch operators
  - Needs kubeconfig secrets for spoke clusters
  - Currently shows 0 without setup

## Architecture

```
Global Hub (vhub)
  └─ Hub (acm1)
      ├─ Hub Operators: ✅ Direct connection via kubeconfig
      └─ Spoke Clusters
          └─ sno146
              ├─ Operators installed ON sno146
              └─ Kubeconfig stored in sno146 namespace on acm1
```

## Setup for Spoke Operators

To enable operator monitoring for spoke clusters, create kubeconfig secrets:

### Option 1: Using Existing Kubeconfig

```bash
# On the hub cluster (e.g., acm1)
# For each spoke cluster (e.g., sno146)

# Create namespace if it doesn't exist
oc create namespace sno146

# Create secret with spoke kubeconfig
oc create secret generic sno146-admin-kubeconfig \
  --from-file=kubeconfig=/path/to/sno146/kubeconfig \
  -n sno146

# Label it for RHACM monitor
oc label secret sno146-admin-kubeconfig \
  created-by=rhacm-monitor \
  -n sno146
```

### Option 2: Extract from RHACM

If RHACM already manages the spoke, it may have the kubeconfig:

```bash
# Check if RHACM has the kubeconfig
oc get secret -n sno146 | grep kubeconfig

# If it exists, just label it
oc label secret <secret-name> created-by=rhacm-monitor -n sno146
```

## Verification

After creating the secrets:

1. Wait ~90 seconds for cache to expire
2. Refresh the RHACM Monitor UI
3. Navigate to Hub → Spoke Clusters tab
4. Check the Operators column - should show count > 0
5. Expand spoke details - should show operators list

## Expected Results

With secrets in place:

```
Spoke: sno146
  Operators: 7 (actual count)
  
  Expanded View:
    - cluster-logging.v6.2.4
    - lifecycle-agent.v4.18.0
    - local-storage-operator.v4.18.0
    - oadp-operator.v1.4.5
    - packageserver
    - sriov-fec.v2.11.1
    - sriov-network-operator.v4.18.0
```

## Known Limitations

- Spoke operators require kubeconfig secrets
- Without secrets, spoke operators show 0
- Hub operators work without additional setup
- This is by design (operators are on spoke, not hub)

## Troubleshooting

**Spoke shows 0 operators:**
1. Check if kubeconfig secret exists:
   ```bash
   oc get secret -n sno146 sno146-admin-kubeconfig
   ```
2. Check backend logs:
   ```bash
   oc logs -l component=backend -n rhacm-monitor | grep sno146
   ```
3. Look for: "Fetched X operators from spoke sno146"
4. Or: "No kubeconfig secret for spoke sno146"

**Secret exists but still shows 0:**
1. Verify kubeconfig is valid
2. Check RBAC permissions
3. Ensure spoke is accessible from hub

## Summary

- **Hub operators:** ✅ Working out of the box
- **Spoke operators:** ⚠️ Requires kubeconfig setup
- **Code:** ✅ Ready to fetch when secrets exist
- **Setup:** Create {spoke}-admin-kubeconfig secrets on hub

---

*For questions or issues, check the backend logs for detailed error messages.*

