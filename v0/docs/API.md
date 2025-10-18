# API Documentation

This document describes the REST API endpoints provided by the RHACM Global Hub Monitor backend.

## Base URL

```
https://<route-hostname>/api
```

## Authentication

All API endpoints (except health checks) require authentication via Bearer token:

```http
Authorization: Bearer <jwt-token>
```

### Obtaining a Token

When using OpenShift SSO, the token is automatically obtained through the OAuth flow.

For manual API access:

```bash
TOKEN=$(oc whoami -t)
curl -H "Authorization: Bearer $TOKEN" https://rhacm-monitor.apps.example.com/api/hubs
```

## Response Format

All API responses follow this format:

### Success Response

```json
{
  "success": true,
  "data": { ... },
  "message": "Optional success message"
}
```

### Error Response

```json
{
  "success": false,
  "error": "Error message",
  "message": "Optional additional context"
}
```

## Endpoints

### Health Endpoints

#### GET /api/health

Health check endpoint.

**Authentication:** Not required

**Response:**

```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

**Example:**

```bash
curl https://rhacm-monitor.apps.example.com/api/health
```

---

#### GET /api/ready

Readiness check endpoint.

**Authentication:** Not required

**Response:**

```json
{
  "status": "ready",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

#### GET /api/live

Liveness check endpoint.

**Authentication:** Not required

**Response:**

```json
{
  "status": "alive",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

### Hub Endpoints

#### GET /api/hubs

List all managed hub clusters.

**Authentication:** Required

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "name": "hub1",
      "namespace": "hub1",
      "status": "Ready",
      "version": "v1.29.0",
      "conditions": [
        {
          "type": "Available",
          "status": "True",
          "lastTransitionTime": "2024-01-01T00:00:00Z",
          "reason": "ClusterAvailable",
          "message": "Cluster is available"
        }
      ],
      "clusterInfo": {
        "clusterID": "abc-123",
        "kubernetesVersion": "v1.29.0",
        "platform": "AWS",
        "region": "us-east-1",
        "openshiftVersion": "4.15.0",
        "consoleURL": "https://console-openshift-console.apps.hub1.example.com",
        "apiURL": "https://api.hub1.example.com:6443",
        "networkType": "OVNKubernetes",
        "createdAt": "2024-01-01T00:00:00Z"
      },
      "nodesInfo": [],
      "policiesInfo": [],
      "managedClusters": [],
      "labels": {
        "vendor": "OpenShift"
      },
      "annotations": {},
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Example:**

```bash
curl -H "Authorization: Bearer $TOKEN" \
  https://rhacm-monitor.apps.example.com/api/hubs
```

---

#### GET /api/hubs/{name}

Get details of a specific hub cluster.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type   | Description |
|-----------|--------|-------------|
| name      | string | Hub name    |

**Response:**

```json
{
  "success": true,
  "data": {
    "name": "hub1",
    "namespace": "hub1",
    "status": "Ready",
    "version": "v1.29.0",
    "conditions": [...],
    "clusterInfo": {...},
    "nodesInfo": [
      {
        "name": "node1",
        "status": "Ready",
        "role": "master",
        "internalIP": "10.0.1.1",
        "externalIP": "203.0.113.1",
        "kernelVersion": "5.14.0-284.el9.x86_64",
        "osImage": "Red Hat Enterprise Linux CoreOS",
        "containerRuntime": "cri-o://1.29.0",
        "kubeletVersion": "v1.29.0",
        "conditions": [...],
        "capacity": {
          "cpu": "8",
          "memory": "32Gi",
          "storage": "100Gi",
          "ephemeralStorage": "100Gi",
          "pods": "250"
        },
        "allocatable": {
          "cpu": "7500m",
          "memory": "30Gi",
          "storage": "95Gi",
          "ephemeralStorage": "95Gi",
          "pods": "250"
        },
        "labels": {},
        "annotations": {},
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "policiesInfo": [
      {
        "name": "policy-compliance",
        "namespace": "policies",
        "remediationAction": "inform",
        "complianceState": "Compliant",
        "severity": "high",
        "categories": ["Security"],
        "standards": ["NIST"],
        "controls": ["CM-2"],
        "violations": 0,
        "placementRules": ["placement-hub1"],
        "disabled": false,
        "labels": {},
        "annotations": {},
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "managedClusters": [...],
    "labels": {},
    "annotations": {},
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

**Example:**

```bash
curl -H "Authorization: Bearer $TOKEN" \
  https://rhacm-monitor.apps.example.com/api/hubs/hub1
```

**Error Responses:**

- `404 Not Found` - Hub not found
- `401 Unauthorized` - Invalid or missing authentication token
- `500 Internal Server Error` - Server error

---

#### GET /api/hubs/{name}/clusters

List all managed clusters for a specific hub.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type   | Description |
|-----------|--------|-------------|
| name      | string | Hub name    |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "name": "sno1",
      "namespace": "sno1",
      "status": "Ready",
      "version": "v1.28.0",
      "conditions": [...],
      "clusterInfo": {
        "clusterID": "def-456",
        "kubernetesVersion": "v1.28.0",
        "platform": "BareMetal",
        "region": "",
        "openshiftVersion": "4.14.0",
        "consoleURL": "https://console-openshift-console.apps.sno1.example.com",
        "apiURL": "https://api.sno1.example.com:6443",
        "networkType": "OVNKubernetes",
        "createdAt": "2024-01-01T00:00:00Z"
      },
      "nodesInfo": [
        {
          "name": "sno1-master-0",
          "status": "Ready",
          "role": "master",
          "internalIP": "10.0.2.1",
          "externalIP": "",
          "kernelVersion": "5.14.0-284.el9.x86_64",
          "osImage": "Red Hat Enterprise Linux CoreOS",
          "containerRuntime": "cri-o://1.28.0",
          "kubeletVersion": "v1.28.0",
          "conditions": [...],
          "capacity": {...},
          "allocatable": {...},
          "labels": {},
          "annotations": {},
          "createdAt": "2024-01-01T00:00:00Z"
        }
      ],
      "policiesInfo": [...],
      "labels": {
        "cluster.open-cluster-management.io/clusterset": "hub1"
      },
      "annotations": {},
      "hubName": "hub1",
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Example:**

```bash
curl -H "Authorization: Bearer $TOKEN" \
  https://rhacm-monitor.apps.example.com/api/hubs/hub1/clusters
```

---

## Data Models

### ManagedHub

| Field            | Type              | Description                           |
|------------------|-------------------|---------------------------------------|
| name             | string            | Hub cluster name                      |
| namespace        | string            | Hub namespace                         |
| status           | string            | Hub status (Ready, NotReady, etc.)    |
| version          | string            | Kubernetes version                    |
| conditions       | Condition[]       | Hub conditions                        |
| clusterInfo      | ClusterInfo       | Cluster information                   |
| nodesInfo        | NodeInfo[]        | Node information                      |
| policiesInfo     | PolicyInfo[]      | Policy information                    |
| managedClusters  | ManagedCluster[]  | Managed spoke clusters                |
| labels           | map[string]string | Kubernetes labels                     |
| annotations      | map[string]string | Kubernetes annotations                |
| createdAt        | string            | Creation timestamp                    |

### ManagedCluster

| Field        | Type              | Description                           |
|--------------|-------------------|---------------------------------------|
| name         | string            | Cluster name                          |
| namespace    | string            | Cluster namespace                     |
| status       | string            | Cluster status                        |
| version      | string            | Kubernetes version                    |
| conditions   | Condition[]       | Cluster conditions                    |
| clusterInfo  | ClusterInfo       | Cluster information                   |
| nodesInfo    | NodeInfo[]        | Node information                      |
| policiesInfo | PolicyInfo[]      | Policy information                    |
| labels       | map[string]string | Kubernetes labels                     |
| annotations  | map[string]string | Kubernetes annotations                |
| hubName      | string            | Managing hub name                     |
| createdAt    | string            | Creation timestamp                    |

### ClusterInfo

| Field             | Type   | Description                |
|-------------------|--------|----------------------------|
| clusterID         | string | Unique cluster identifier  |
| kubernetesVersion | string | Kubernetes version         |
| platform          | string | Platform (AWS, Azure, etc.)|
| region            | string | Cloud region               |
| openshiftVersion  | string | OpenShift version          |
| consoleURL        | string | Console URL                |
| apiURL            | string | API server URL             |
| networkType       | string | Network plugin type        |
| createdAt         | string | Creation timestamp         |

### NodeInfo

| Field            | Type              | Description              |
|------------------|-------------------|--------------------------|
| name             | string            | Node name                |
| status           | string            | Node status              |
| role             | string            | Node role                |
| internalIP       | string            | Internal IP address      |
| externalIP       | string            | External IP address      |
| kernelVersion    | string            | Kernel version           |
| osImage          | string            | OS image                 |
| containerRuntime | string            | Container runtime        |
| kubeletVersion   | string            | Kubelet version          |
| conditions       | Condition[]       | Node conditions          |
| capacity         | ResourceList      | Node capacity            |
| allocatable      | ResourceList      | Allocatable resources    |
| labels           | map[string]string | Node labels              |
| annotations      | map[string]string | Node annotations         |
| createdAt        | string            | Creation timestamp       |

### PolicyInfo

| Field             | Type              | Description                      |
|-------------------|-------------------|----------------------------------|
| name              | string            | Policy name                      |
| namespace         | string            | Policy namespace                 |
| remediationAction | string            | Remediation action               |
| complianceState   | string            | Compliance state                 |
| severity          | string            | Policy severity                  |
| categories        | string[]          | Policy categories                |
| standards         | string[]          | Compliance standards             |
| controls          | string[]          | Control identifiers              |
| violations        | int               | Number of violations             |
| placementRules    | string[]          | Placement rules                  |
| disabled          | bool              | Whether policy is disabled       |
| labels            | map[string]string | Policy labels                    |
| annotations       | map[string]string | Policy annotations               |
| createdAt         | string            | Creation timestamp               |

### Condition

| Field               | Type   | Description                    |
|---------------------|--------|--------------------------------|
| type                | string | Condition type                 |
| status              | string | Condition status               |
| lastTransitionTime  | string | Last transition time           |
| reason              | string | Reason for condition           |
| message             | string | Human-readable message         |

### ResourceList

| Field            | Type   | Description                |
|------------------|--------|----------------------------|
| cpu              | string | CPU resources              |
| memory           | string | Memory resources           |
| storage          | string | Storage resources          |
| ephemeralStorage | string | Ephemeral storage          |
| pods             | string | Maximum pods               |

## Error Codes

| Status Code | Description                                    |
|-------------|------------------------------------------------|
| 200         | Success                                        |
| 400         | Bad Request - Invalid input                    |
| 401         | Unauthorized - Invalid or missing token        |
| 403         | Forbidden - Insufficient permissions           |
| 404         | Not Found - Resource not found                 |
| 500         | Internal Server Error - Server error           |

## Rate Limiting

Currently, no rate limiting is implemented. Future versions may include rate limiting for API protection.

## Versioning

The API is currently at version 1.0. Future breaking changes will be introduced under new API versions (e.g., `/api/v2/`).

