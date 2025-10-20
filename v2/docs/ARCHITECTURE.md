# Architecture Documentation

## System Architecture

The RHACM Global Hub Monitor follows a modern three-tier architecture optimized for cloud-native deployments.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Browser                            │
│                    (React + TypeScript)                         │
└───────────────────────────┬─────────────────────────────────────┘
                            │ HTTPS
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                    OpenShift Route                              │
│                    (TLS Termination)                            │
└───────────────────────────┬─────────────────────────────────────┘
                            │
        ┌───────────────────┴───────────────────┐
        │                                       │
        │ Static Files                          │ API Requests
        ▼                                       ▼
┌─────────────────┐                   ┌─────────────────────┐
│  Frontend Pod   │                   │   Backend Pod       │
│  (nginx + SPA)  │◄──────Proxy───────│   (Go API Server)  │
└─────────────────┘                   └──────────┬──────────┘
                                                 │
                            ┌────────────────────┼────────────────────┐
                            │                    │                    │
                            ▼                    ▼                    ▼
                    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
                    │ Kubernetes   │    │   RHACM      │    │  OpenShift   │
                    │   API        │    │   API        │    │   OAuth      │
                    └──────────────┘    └──────────────┘    └──────────────┘
                            │                    │
                            └────────────────────┘
                                       │
                        ┌──────────────┴──────────────┐
                        │                             │
                        ▼                             ▼
                ┌───────────────┐           ┌───────────────┐
                │ Managed Hubs  │           │    Policies   │
                │   (Clusters)  │           │   (Resources) │
                └───────────────┘           └───────────────┘
                        │
                        └─── Manages ───►  SNO Spoke Clusters
```

## Component Architecture

### Frontend Layer

```
┌────────────────────────────────────────────────────────────┐
│                      React Application                      │
├────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │   Pages      │  │  Components  │  │    Routing   │    │
│  │ - Dashboard  │  │ - Layout     │  │ React Router │    │
│  │ - HubsList   │  │ - ClusterCard│  │              │    │
│  │ - HubDetail  │  │ - StatusLabel│  │              │    │
│  │ - ClusterDtl │  │ - ...        │  │              │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
│         │                  │                  │            │
│         └──────────────────┴──────────────────┘            │
│                            │                               │
│                            ▼                               │
│              ┌─────────────────────────┐                   │
│              │    State Management     │                   │
│              │    (React Query)        │                   │
│              └────────────┬────────────┘                   │
│                           │                                │
│                           ▼                                │
│              ┌─────────────────────────┐                   │
│              │     API Services        │                   │
│              │   (Axios Client)        │                   │
│              └─────────────────────────┘                   │
│                                                             │
└──────────────────────────┬──────────────────────────────────┘
                           │ REST API
                           ▼
```

### Backend Layer

```
┌────────────────────────────────────────────────────────────┐
│                    Go Backend Server                        │
├────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ HTTP Server  │──│  Middleware  │──│    Router    │    │
│  │    (Gin)     │  │ - Auth       │  │   (Routes)   │    │
│  │              │  │ - CORS       │  │              │    │
│  │              │  │ - Logging    │  │              │    │
│  └──────────────┘  └──────────────┘  └──────┬───────┘    │
│                                              │             │
│                           ┌──────────────────┴────────┐    │
│                           │                           │    │
│                  ┌────────▼────────┐       ┌──────────▼──┐│
│                  │    Handlers     │       │   Auth      ││
│                  │ - HubHandler    │       │ - JWT       ││
│                  │ - HealthHandler │       │ - OAuth     ││
│                  └────────┬────────┘       └─────────────┘│
│                           │                                │
│                  ┌────────▼────────┐                       │
│                  │  RHACM Client   │                       │
│                  │ - GetHubs       │                       │
│                  │ - GetClusters   │                       │
│                  └────────┬────────┘                       │
│                           │                                │
│                  ┌────────▼────────┐                       │
│                  │ Kubernetes      │                       │
│                  │ Client          │                       │
│                  │ (client-go)     │                       │
│                  └─────────────────┘                       │
│                                                             │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
                  Kubernetes API Server
```

## Data Flow

### Listing Hubs Flow

```
1. User requests hub list
   │
   ▼
2. Frontend calls GET /api/hubs
   │
   ▼
3. Backend authenticates request (JWT validation)
   │
   ▼
4. HubHandler.ListHubs() called
   │
   ▼
5. RHACMClient.GetManagedHubs() called
   │
   ▼
6. Kubernetes client queries ManagedCluster resources
   │
   ▼
7. Filter clusters with hub labels
   │
   ▼
8. For each hub, fetch:
   - Cluster info (from cluster claims)
   - Nodes info (from Node resources)
   - Policies info (from Policy resources)
   - Managed clusters (from ManagedCluster with hub labels)
   │
   ▼
9. Transform K8s objects to API models
   │
   ▼
10. Return JSON response
    │
    ▼
11. Frontend updates UI with React Query cache
```

### Authentication Flow

```
1. User accesses application
   │
   ▼
2. Frontend redirects to OpenShift OAuth
   │
   ▼
3. User logs in with OpenShift credentials
   │
   ▼
4. OAuth returns JWT token
   │
   ▼
5. Frontend stores token in localStorage
   │
   ▼
6. All API requests include: Authorization: Bearer <token>
   │
   ▼
7. Backend middleware validates token:
   - Extract token from header
   - Verify signature with JWKS
   - Check expiration
   - Extract user info
   │
   ▼
8. If valid: proceed with request
   If invalid: return 401 Unauthorized
```

## Deployment Architecture

### Kubernetes Resources

```
┌─────────────────────────────────────────────────────────────┐
│                    Namespace: rhacm-monitor                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────┐      │
│  │              Backend Deployment                  │      │
│  │  Replicas: 2                                     │      │
│  │  ┌────────────┐        ┌────────────┐           │      │
│  │  │  Pod 1     │        │  Pod 2     │           │      │
│  │  │ Container: │        │ Container: │           │      │
│  │  │  backend   │        │  backend   │           │      │
│  │  └────────────┘        └────────────┘           │      │
│  └──────────────────┬───────────────────────────────┘      │
│                     │                                       │
│                     ▼                                       │
│           ┌─────────────────┐                               │
│           │ Backend Service │                               │
│           │  ClusterIP      │                               │
│           │  Port: 8080     │                               │
│           └─────────────────┘                               │
│                                                              │
│  ┌──────────────────────────────────────────────────┐      │
│  │             Frontend Deployment                  │      │
│  │  Replicas: 2                                     │      │
│  │  ┌────────────┐        ┌────────────┐           │      │
│  │  │  Pod 1     │        │  Pod 2     │           │      │
│  │  │ Container: │        │ Container: │           │      │
│  │  │  frontend  │        │  frontend  │           │      │
│  │  └────────────┘        └────────────┘           │      │
│  └──────────────────┬───────────────────────────────┘      │
│                     │                                       │
│                     ▼                                       │
│           ┌─────────────────┐                               │
│           │Frontend Service │                               │
│           │  ClusterIP      │                               │
│           │  Port: 80       │                               │
│           └────────┬────────┘                               │
│                    │                                        │
└────────────────────┼────────────────────────────────────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │   OpenShift Route     │
         │   TLS: Edge           │
         │   Host: rhacm-monitor │
         └───────────────────────┘
```

### RBAC Architecture

```
┌────────────────────────────────────────────────────┐
│              ServiceAccount                        │
│           rhacm-monitor                            │
└──────────────────┬─────────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────────┐
│           ClusterRoleBinding                       │
│         rhacm-monitor                              │
└──────────────────┬─────────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────────┐
│            ClusterRole                             │
│          rhacm-monitor                             │
├────────────────────────────────────────────────────┤
│  Permissions:                                      │
│  - Read ManagedClusters                            │
│  - Read Policies                                   │
│  - Read Nodes                                      │
│  - Read Secrets (for hub kubeconfigs)              │
└────────────────────────────────────────────────────┘
```

## Security Architecture

### Network Security

```
┌────────────────────────────────────────────────────┐
│              External Traffic                      │
│          (HTTPS - Port 443)                        │
└──────────────────┬─────────────────────────────────┘
                   │ TLS Termination
                   ▼
┌────────────────────────────────────────────────────┐
│           OpenShift Router                         │
│         (Edge TLS Termination)                     │
└──────────────────┬─────────────────────────────────┘
                   │ HTTP (internal)
                   ▼
┌────────────────────────────────────────────────────┐
│          Frontend Service                          │
│            (ClusterIP)                             │
└──────────────────┬─────────────────────────────────┘
                   │
     ┌─────────────┴─────────────┐
     │                           │
     ▼                           ▼
┌─────────┐              ┌──────────────┐
│Frontend │              │   Backend    │
│  Pods   │──API Proxy──►│    Pods      │
└─────────┘              └──────┬───────┘
                                │
                                │ Service Account Token
                                ▼
                    ┌───────────────────────┐
                    │   Kubernetes API      │
                    │   (TLS - Port 6443)   │
                    └───────────────────────┘
```

### Authentication & Authorization

```
┌──────────────┐
│    User      │
└──────┬───────┘
       │ 1. Login
       ▼
┌──────────────────┐
│ OpenShift OAuth  │
└──────┬───────────┘
       │ 2. JWT Token
       ▼
┌──────────────────┐
│   Frontend       │
└──────┬───────────┘
       │ 3. Bearer Token
       ▼
┌──────────────────┐
│   Backend        │──► 4. Validate JWT
│   Middleware     │
└──────┬───────────┘
       │ 5. User Context
       ▼
┌──────────────────┐
│   Handler        │
└──────┬───────────┘
       │ 6. Service Account
       ▼
┌──────────────────┐
│ Kubernetes API   │──► 7. RBAC Check
└──────────────────┘
```

## Scalability

### Horizontal Scaling

Both frontend and backend support horizontal scaling:

```
Load Balancer (OpenShift Router)
         │
         ├──► Frontend Pod 1
         ├──► Frontend Pod 2
         └──► Frontend Pod N

Backend Service (LoadBalanced)
         │
         ├──► Backend Pod 1 ──► Kubernetes API
         ├──► Backend Pod 2 ──► Kubernetes API
         └──► Backend Pod N ──► Kubernetes API
```

### Caching Strategy

```
┌──────────────┐
│   Browser    │
└──────┬───────┘
       │
       ▼
┌──────────────────┐
│  React Query     │  ← 30s stale time
│    Cache         │  ← Automatic refetch
└──────┬───────────┘
       │ Cache miss or stale
       ▼
┌──────────────────┐
│   Backend API    │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Kubernetes API   │  ← Built-in caching
│   (informers)    │  ← Watch for changes
└──────────────────┘
```

## Technology Stack Details

### Backend Stack

```
┌─────────────────────────────────────────┐
│         Application Layer               │
│  - Gin Framework (HTTP)                 │
│  - Custom Handlers & Middleware         │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│        Business Logic Layer             │
│  - RHACM Client                         │
│  - Data Transformation                  │
│  - Error Handling                       │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Data Access Layer               │
│  - Kubernetes Client-Go                 │
│  - RHACM API Client                     │
│  - Dynamic Client                       │
└─────────────────────────────────────────┘
```

### Frontend Stack

```
┌─────────────────────────────────────────┐
│         Presentation Layer              │
│  - React Components                     │
│  - PatternFly UI                        │
│  - CSS Modules                          │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│        Application Layer                │
│  - React Router (Navigation)            │
│  - React Query (State)                  │
│  - Custom Hooks                         │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Service Layer                   │
│  - Axios HTTP Client                    │
│  - API Service Functions                │
│  - Error Handling                       │
└─────────────────────────────────────────┘
```

## Performance Characteristics

### Response Times (Target)
- Health checks: < 10ms
- Hub list: < 500ms
- Hub details: < 1s
- Static assets: < 100ms

### Resource Usage (Per Pod)
- Backend: 512Mi memory, 500m CPU (normal operation)
- Frontend: 256Mi memory, 250m CPU (normal operation)

### Concurrent Users
- Designed for: 100+ concurrent users
- Backend scales horizontally
- Frontend served via CDN-like nginx

## Monitoring & Observability

### Health Checks

```
┌────────────────────────────────────────┐
│     Kubernetes Probes                  │
├────────────────────────────────────────┤
│                                        │
│  Liveness Probe:                       │
│    GET /api/live                       │
│    Every 10s                           │
│                                        │
│  Readiness Probe:                      │
│    GET /api/ready                      │
│    Every 5s                            │
│                                        │
└────────────────────────────────────────┘
```

### Logging

```
Backend Logs:
  - Structured JSON logging
  - Request/response logging
  - Error tracking
  - Performance metrics

Frontend Logs:
  - Console logging (dev)
  - Error boundaries
  - API error tracking
```

## Future Architecture Considerations

### Potential Enhancements

1. **Caching Layer**: Redis for shared cache
2. **Message Queue**: For async operations
3. **Database**: For historical data and analytics
4. **Metrics**: Prometheus integration
5. **Tracing**: OpenTelemetry for distributed tracing
6. **WebSockets**: Real-time updates
7. **GraphQL**: Alternative to REST API

### Scalability Improvements

1. **Database**: For persistence and complex queries
2. **Event-driven**: Kubernetes watch events
3. **Microservices**: Split by domain (hubs, clusters, policies)
4. **API Gateway**: Rate limiting, authentication
5. **CDN**: For static assets

## Conclusion

This architecture provides:
- ✅ Scalability through horizontal scaling
- ✅ Security through authentication and RBAC
- ✅ Maintainability through clean separation of concerns
- ✅ Observability through health checks and logging
- ✅ Performance through caching and efficient queries
- ✅ Resilience through multiple replicas and health checks

