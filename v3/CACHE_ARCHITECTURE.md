# Cache Architecture - v3

## Overview

v3 implements a 30-minute in-memory cache with session affinity to ensure consistent performance across multiple backend pod replicas.

## Cache Implementation

### In-Memory Cache

**Location:** `backend/pkg/cache/cache.go`

**Features:**
- Thread-safe with mutex locks
- Configurable TTL (Time To Live)
- Automatic expiration
- Simple Get/Set/Delete operations

### Shared Cache Instance

**Created in:** `backend/cmd/server/main.go`

```go
// Create shared cache instance (30 minute TTL)
sharedCache := cache.NewCache(30 * time.Minute)

// Pass to handlers
hubHandler := handlers.NewHubHandler(rhacmClient, sharedCache)
```

**Benefits:**
- Single cache instance per pod
- Shared across all handlers in that pod
- Consistent cache state within pod
- 30-minute lifetime

### Cache Keys

**Hub List:**
- Key: `"hubs:list"`
- Stores: Array of all hubs
- Used by: GET /api/hubs

**Individual Hub:**
- Key: `"hub:{name}"`
- Stores: Single hub data
- Used by: GET /api/hubs/:name

## Multi-Pod Challenge

### The Problem

With 2+ backend pod replicas:
- Each pod has its own in-memory cache
- Load balancer distributes requests randomly
- Request to Pod A: Fast (if cached in Pod A)
- Request to Pod B: Slow (not cached in Pod B)
- Inconsistent performance for users

### Solution: Session Affinity

**Configuration:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: rhacm-monitor-backend
spec:
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800  # 3 hours
```

**How It Works:**
- Requests from same client IP always go to same pod
- That pod's cache is consistent for that client
- User gets consistent fast performance
- 3-hour timeout covers typical usage sessions

**Command:**
```bash
oc patch svc rhacm-monitor-backend -n rhacm-monitor \
  -p '{"spec":{"sessionAffinity":"ClientIP","sessionAffinityConfig":{"clientIP":{"timeoutSeconds":10800}}}}'
```

### Performance with Session Affinity

**Test Results:**
```
Request 1: 11.368s (cache miss - initial)
Request 2: 0.028s (cache hit - 406x faster!)
Request 3: 1.237s (network blip or new endpoint)
Request 4: 0.029s (cache hit - 392x faster!)
Request 5: 0.036s (cache hit - 316x faster!)
```

**Average:** ~30ms for cached requests (300x+ improvement)

## Cache Refresh

### Automatic Expiration

- Cache entries expire after 30 minutes
- Next request after expiry fetches fresh data
- New data cached for another 30 minutes

### Manual Refresh

**Per-Hub Refresh:**
- Endpoint: POST /api/hubs/:name/refresh
- Clears cache for specific hub
- Next request fetches fresh data
- UI: 🔄 button on each hub card

**Client-Side Cache:**
- window.cachedHubsData for instant navigation
- Stores homepage data in browser
- Used when clicking "Back to Hubs"
- Instant display without API call

## Monitoring Cache Performance

### Backend Logs

Watch for timing in logs:
```bash
oc logs -l component=backend -n rhacm-monitor --tail=50 | grep "GET.*hubs"
```

**Fast requests (cached):**
```
[GIN] 2025/10/21 - 17:31:41 | 200 | 6.080756ms | GET "/api/hubs"
```

**Slow requests (cache miss):**
```
[GIN] 2025/10/21 - 17:31:36 | 200 | 10.688436159s | GET "/api/hubs"
```

### Testing Cache

**Test cache hit rate:**
```bash
# First request (cache miss)
time curl -k https://monitor/api/hubs

# Second request (should be cached)
time curl -k https://monitor/api/hubs
```

**Expected:**
- First: ~10s
- Second: <100ms

## Cache Strategy

### What's Cached

**Cached:**
- Hub list (GET /api/hubs)
- Individual hub data (GET /api/hubs/:name)
- Includes: nodes, policies, managed clusters
- Excludes: Spoke operators (lazy loaded separately)

**Not Cached:**
- Health endpoints (/health, /ready, /live)
- Spoke operators (/hubs/:name/spokes/:spoke/operators)
- Add/remove hub operations
- Policy YAML downloads
- CGU creation

### Cache Invalidation

**Automatic:**
- 30-minute TTL expiration
- Pod restart (cache is in-memory)

**Manual:**
- Click 🔄 refresh button (per hub)
- POST /api/hubs/:name/refresh

**Best Practices:**
- Let cache work naturally (30 min is good)
- Use refresh button when you need fresh data
- Cache balances performance with freshness

## Alternative Solutions

### If Sticky Sessions Aren't Enough

**Redis Cache (Future Enhancement):**
```go
import "github.com/go-redis/redis/v8"

// In main.go
redisClient := redis.NewClient(&redis.Options{
    Addr: "redis:6379",
})

sharedCache := cache.NewRedisCache(redisClient, 30*time.Minute)
```

**Benefits:**
- True sharing across all pods
- Survives pod restarts
- Scales horizontally

**Trade-offs:**
- Additional dependency
- More complex setup
- External service to manage

### Current Recommendation

**Stick with Session Affinity:**
- Simple configuration
- Works well with 2-3 pods
- No external dependencies
- Good performance
- Easy to maintain

## Summary

**Current Setup:**
- In-memory cache per pod (30-min TTL)
- Session affinity for consistency
- Per-hub manual refresh
- Client-side navigation cache

**Performance:**
- 300x+ improvement on cache hits
- Consistent with sticky sessions
- Scales to multiple pods
- User control with refresh buttons

**Monitoring:**
- Check backend logs for timing
- Slow requests (>1s) indicate cache miss
- Fast requests (<100ms) indicate cache hit

---

*The current cache architecture provides excellent performance while maintaining simplicity.*

