# RHACM Global Hub Monitor v4.0.0 - Roadmap

**Status:** 🚧 Planning Phase  
**Target Release:** Q1 2026

## Overview

v4 focuses on resolving v3 technical debt and adding advanced monitoring capabilities. The primary goals are to fix authentication for unmanaged hubs, stabilize caching, and enhance observability.

## v3 Technical Debt to Address

### 1. Unmanaged Hub Authentication (HIGH PRIORITY)

**Current Issue:**
- Unmanaged hubs connect as "system:anonymous"
- Cannot fetch routes, nodes, cluster version
- Status shows "Unknown", limited functionality

**Root Cause:**
```go
// Current approach
config, err := clientcmd.NewClientConfigFromBytes(kubeconfigData)
// Results in anonymous authentication
```

**Planned Solutions (in priority order):**

**Option A: Direct REST Client with Auth (Recommended)**
```go
// Load kubeconfig properly
kubeConfig, err := clientcmd.Load(kubeconfigData)
if err != nil {
    return nil, err
}

// Get the first cluster and user
var cluster *clientcmdapi.Cluster
var authInfo *clientcmdapi.AuthInfo

for _, c := range kubeConfig.Clusters {
    cluster = c
    break
}

for _, u := range kubeConfig.AuthInfos {
    authInfo = u  
    break
}

// Build REST config manually with auth
config := &rest.Config{
    Host: cluster.Server,
    TLSClientConfig: rest.TLSClientConfig{
        CAData:   cluster.CertificateAuthorityData,
        CertData: authInfo.ClientCertificateData,
        KeyData:  authInfo.ClientKeyData,
    },
    BearerToken: authInfo.Token,
}

if authInfo.Username != "" {
    config.Username = authInfo.Username
    config.Password = authInfo.Password
}
```

**Option B: Exec Plugin Support**
- Handle exec-based authentication (like AWS EKS)
- Implement exec.Command execution
- Cache exec results

**Option C: Service Account Token**
- Extract service account token from kubeconfig
- Use TokenReview API for validation
- Refresh tokens as needed

**Expected Outcome:**
- ✅ Unmanaged hubs authenticate correctly
- ✅ Full cluster information fetched
- ✅ Console/GitOps URLs displayed
- ✅ Proper status (Ready/NotReady)

**Testing Plan:**
1. Test with certificate-based auth
2. Test with token-based auth
3. Test with username/password
4. Test with exec plugins (AWS, Azure, GCP)

### 2. Cache TTL Stability (MEDIUM PRIORITY)

**Current Issue:**
- Cache configured for 30 minutes
- Actually expires every ~2 minutes
- Inconsistent performance

**Evidence:**
```
Timeline from logs:
18:00:04 - 10.665s (cache miss)
18:02:22 - 10.535s (cache miss) ← Only 2min 18sec later
```

**Investigation Steps:**

1. **Add Comprehensive Cache Logging:**
```go
func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    expiresAt := time.Now().Add(c.ttl)
    c.items[key] = CacheEntry{
        Data:      value,
        ExpiresAt: expiresAt,
    }
    log.Printf("Cache SET: key=%s, expires=%s (TTL=%s)", 
        key, expiresAt.Format(time.RFC3339), c.ttl)
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    entry, found := c.items[key]
    if !found {
        log.Printf("Cache MISS: key=%s (not found)", key)
        return nil, false
    }
    
    now := time.Now()
    if now.After(entry.ExpiresAt) {
        log.Printf("Cache EXPIRED: key=%s, expired=%s ago", 
            key, now.Sub(entry.ExpiresAt))
        return nil, false
    }
    
    log.Printf("Cache HIT: key=%s, expires in=%s", 
        key, entry.ExpiresAt.Sub(now))
    return entry.Data, true
}
```

2. **Check for Cache Clearing:**
- Search for `cache.Clear()` calls
- Check if cleanup goroutine is too aggressive
- Monitor cache size over time

3. **Possible Issues:**
- Cleanup goroutine deleting too early
- TTL not being set correctly
- Pod restarts clearing cache
- Multiple cache instances despite shared code

**Planned Fixes:**

**Option A: Redis/Valkey Cache (Recommended for HA)**
```go
import "github.com/go-redis/redis/v8"

type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(addr string, ttl time.Duration) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    return &RedisCache{client: client, ttl: ttl}
}

func (c *RedisCache) Set(key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, c.ttl).Err()
}
```

**Benefits:**
- Shared across all pods
- Survives pod restarts
- Scales horizontally
- True 30-minute TTL

**Option B: Debug and Fix In-Memory Cache**
- Add extensive logging
- Identify root cause
- Fix cleanup logic
- Verify TTL calculation

**Expected Outcome:**
- ✅ Reliable 30-minute cache
- ✅ Consistent performance
- ✅ Shared across pods (if Redis)

### 3. Image Deployment Issues (LOW PRIORITY)

**Current Issue:**
- imagePullPolicy: Always doesn't always pull
- Nodes cache images
- Latest code not deployed immediately

**Planned Solutions:**

**Option A: Use Specific Tags**
```yaml
image: quay.io/bzhai/rhacm-monitor-backend:v4.0.0
# Not :latest which can be cached
```

**Option B: Image Digest**
```yaml
image: quay.io/bzhai/rhacm-monitor-backend@sha256:abc123...
# Guaranteed unique
```

**Option C: Update Deployment Automation**
```bash
# In deploy.sh
oc set image deployment/rhacm-monitor-backend \
  rhacm-monitor-backend=quay.io/bzhai/rhacm-monitor-backend:v4.0.0
oc rollout restart deployment/rhacm-monitor-backend -n rhacm-monitor
```

## New Features for v4

### 1. Metrics and Monitoring

**Feature:** Prometheus metrics export

**Implementation:**
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    hubsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "rhacm_monitor_hubs_total",
        Help: "Total number of monitored hubs",
    })
    
    spokesTotal = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "rhacm_monitor_spokes_total",
        Help: "Total number of spoke clusters",
    })
    
    cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "rhacm_monitor_cache_hits_total",
        Help: "Total cache hits",
    })
)

// Endpoint: GET /metrics
```

**Benefits:**
- Monitor application health
- Track cache hit rates
- Alert on issues
- Integration with existing monitoring

### 2. Alert Rules

**Feature:** Define alert conditions

**Examples:**
- Hub disconnected for > 5 minutes
- Policy compliance < 95%
- Spoke cluster down
- Operator version mismatch

**Implementation:**
- Backend evaluates rules
- Frontend displays alerts badge
- Email/webhook notifications (optional)

### 3. Historical Data

**Feature:** Time-series data for trends

**Implementation:**
- Optional database backend (PostgreSQL/TimescaleDB)
- Track hub status over time
- Policy compliance trends
- Node inventory changes

**UI Enhancements:**
- Charts/graphs for trends
- Historical view toggle
- Export to CSV

### 4. Multi-Tenancy

**Feature:** Role-based access control

**Implementation:**
- User authentication (OAuth/OIDC)
- Hub access permissions
- Read-only vs admin roles

**Use Cases:**
- Different teams monitoring different hubs
- Audit logging
- Compliance tracking

### 5. Operator Lifecycle Management

**Feature:** Operator version tracking and recommendations

**Implementation:**
- Track operator versions across hubs
- Highlight outdated operators
- Recommend upgrades
- Link to operator documentation

**UI:**
- Version comparison table
- Upgrade recommendations
- CVE tracking (if available)

## Development Plan

### Phase 1: Foundation (Weeks 1-2)
- ✅ Fix unmanaged hub authentication
- ✅ Stabilize cache (Redis or fix in-memory)
- ✅ Add comprehensive logging
- ✅ Testing framework

### Phase 2: Monitoring (Weeks 3-4)
- Prometheus metrics
- Health check improvements
- Performance monitoring
- Cache hit rate tracking

### Phase 3: Features (Weeks 5-6)
- Alert rules system
- Historical data (optional)
- UI enhancements for new data

### Phase 4: Polish (Week 7-8)
- Documentation updates
- Migration guide from v3
- Performance testing
- Security review

## Technical Approach

### Architecture Changes

**v3 Architecture:**
```
Frontend (Nginx) → Backend (Go) → Kubernetes API
                    ↓
                In-Memory Cache (per pod)
```

**v4 Architecture:**
```
Frontend (Nginx) → Backend (Go) → Kubernetes API
                    ↓              ↓
                Redis Cache ← Metrics (Prometheus)
                    ↓
            (Optional) Database for history
```

### Backward Compatibility

v4 will be backward compatible with v3:
- Same API endpoints (with additions)
- Same deployment model
- Migration path from v3
- Optional features can be disabled

### Dependencies

**New Dependencies:**
- go-redis/redis (for cache) - Optional
- prometheus/client_golang (for metrics)
- (Optional) PostgreSQL driver for history

**Trade-offs:**
- More complexity vs better features
- External dependencies vs self-contained
- Performance vs functionality

## Success Criteria

### Must Have
- ✅ Unmanaged hub authentication working
- ✅ Cache TTL reliable (30 minutes)
- ✅ All v3 features maintained
- ✅ Prometheus metrics endpoint

### Should Have
- ✅ Alert rules system
- ✅ Improved error handling
- ✅ Better logging
- ✅ Migration automation

### Nice to Have
- Historical data tracking
- Multi-tenancy
- Operator lifecycle management
- Advanced analytics

## Risk Assessment

### High Risk
- **Breaking existing deployments:** Mitigated by backward compatibility
- **Performance regression:** Mitigated by testing
- **Redis dependency:** Made optional with fallback

### Medium Risk
- **Increased complexity:** Managed with documentation
- **Learning curve:** Addressed with guides

### Low Risk
- **Storage requirements:** Optional features
- **Network latency:** Redis local deployment

## Migration from v3 to v4

### Automated Migration Script

```bash
#!/bin/bash
# migrate-v3-to-v4.sh

echo "Migrating from v3 to v4..."

# 1. Backup v3 configuration
oc get secrets -n rhacm-monitor -l created-by=rhacm-monitor -o yaml > v3-backup.yaml

# 2. Move old unmanaged hub secrets to rhacm-monitor namespace
for ns in $(oc get ns -o name | cut -d'/' -f2); do
    secret="${ns}-admin-kubeconfig"
    if oc get secret -n $ns $secret 2>/dev/null | grep -q created-by=rhacm-monitor; then
        echo "Moving $secret from $ns to rhacm-monitor..."
        oc get secret -n $ns $secret -o yaml | \
            sed "s/namespace: $ns/namespace: rhacm-monitor/" | \
            oc apply -f -
        oc delete secret -n $ns $secret
    fi
done

# 3. Deploy v4
cd v4/deployment
./deploy.sh

# 4. (Optional) Deploy Redis for shared cache
# kubectl apply -f redis-deployment.yaml

echo "Migration complete!"
```

### Manual Migration Steps

1. **Backup v3 Data:**
```bash
oc get all -n rhacm-monitor -o yaml > v3-full-backup.yaml
oc get secrets -n rhacm-monitor -o yaml > v3-secrets-backup.yaml
```

2. **Update Backend Image:**
```bash
oc set image deployment/rhacm-monitor-backend \
  rhacm-monitor-backend=quay.io/bzhai/rhacm-monitor-backend:v4.0.0 \
  -n rhacm-monitor
```

3. **Deploy Redis (Optional):**
```bash
oc apply -f v4/deployment/k8s/redis.yaml
```

4. **Update Configuration:**
```bash
# If using Redis
oc set env deployment/rhacm-monitor-backend \
  CACHE_TYPE=redis \
  REDIS_ADDR=redis:6379 \
  -n rhacm-monitor
```

## Timeline

| Week | Focus | Deliverable |
|------|-------|-------------|
| 1-2 | Authentication fix | Unmanaged hubs working |
| 3-4 | Cache stability | Redis integration |
| 5-6 | Metrics & monitoring | Prometheus endpoint |
| 7-8 | Testing & polish | v4.0.0 release |

## Success Metrics

**Performance:**
- Cache hit rate > 95%
- Page load < 100ms (cached)
- No anonymous auth errors

**Functionality:**
- Unmanaged hubs: Full feature parity with managed
- Metrics: 10+ key metrics exported
- Alerts: Basic rule system working

**Quality:**
- Test coverage > 70%
- Zero critical bugs
- Documentation complete

## Open Questions

1. **Redis Deployment:**
   - Single instance or cluster?
   - Persistent storage needed?
   - How to handle Redis unavailability?

2. **Historical Data:**
   - Worth the complexity in v4?
   - Or defer to v5?
   - Storage requirements?

3. **Breaking Changes:**
   - Any API changes needed?
   - Configuration format changes?
   - Migration complexity?

## Next Steps

1. **Immediate:**
   - Prototype authentication fix
   - Test with real unmanaged hub kubeconfig
   - Verify it resolves anonymous issue

2. **Short Term:**
   - Set up Redis test environment
   - Implement cache with Redis
   - Test TTL stability

3. **Medium Term:**
   - Design metrics schema
   - Implement Prometheus endpoint
   - Create sample dashboards

---

*v4 planning is in progress. All v3 features will be maintained while addressing technical debt.*

