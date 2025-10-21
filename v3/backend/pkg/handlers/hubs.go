package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rhacm-global-hub-monitor/backend/pkg/cache"
	"github.com/rhacm-global-hub-monitor/backend/pkg/client"
	"github.com/rhacm-global-hub-monitor/backend/pkg/models"
)

// HubHandler handles hub-related requests
type HubHandler struct {
	rhacmClient *client.RHACMClient
	cache       *cache.Cache
}

// NewHubHandler creates a new hub handler with caching
func NewHubHandler(rhacmClient *client.RHACMClient) *HubHandler {
	return &HubHandler{
		rhacmClient: rhacmClient,
		cache:       cache.NewCache(30 * time.Minute), // Cache for 30 minutes
	}
}

// ListHubs godoc
// @Summary List all managed hubs
// @Description Get a list of all managed RHACM hub clusters
// @Tags hubs
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/hubs [get]
func (h *HubHandler) ListHubs(c *gin.Context) {
	ctx := c.Request.Context()

	// Try cache first
	cacheKey := "hubs:list"
	if cached, found := h.cache.Get(cacheKey); found {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Data:    cached,
		})
		return
	}

	// Cache miss - fetch from clusters
	hubs, err := h.rhacmClient.GetManagedHubs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Store in cache
	h.cache.Set(cacheKey, hubs)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    hubs,
	})
}

// GetHub godoc
// @Summary Get a specific hub
// @Description Get details of a specific managed hub
// @Tags hubs
// @Accept json
// @Produce json
// @Param name path string true "Hub name"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/hubs/{name} [get]
func (h *HubHandler) GetHub(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Try cache first
	cacheKey := "hub:" + name
	if cached, found := h.cache.Get(cacheKey); found {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Data:    cached,
		})
		return
	}

	// Cache miss - fetch from cluster
	hub, err := h.rhacmClient.GetManagedHub(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Store in cache
	h.cache.Set(cacheKey, hub)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    hub,
	})
}

// ListHubClusters godoc
// @Summary List clusters for a hub
// @Description Get a list of all managed clusters for a specific hub
// @Tags hubs
// @Accept json
// @Produce json
// @Param name path string true "Hub name"
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/hubs/{name}/clusters [get]
func (h *HubHandler) ListHubClusters(c *gin.Context) {
	ctx := c.Request.Context()
	hubName := c.Param("name")

	clusters, err := h.rhacmClient.GetManagedClustersForHub(ctx, hubName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    clusters,
	})
}

// RefreshHubsCache clears the hubs list cache
func (h *HubHandler) RefreshHubsCache(c *gin.Context) {
	h.cache.Delete("hubs:list")
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "Cache cleared",
	})
}
