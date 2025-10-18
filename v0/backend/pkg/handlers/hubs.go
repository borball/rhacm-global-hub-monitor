package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhacm-global-hub-monitor/backend/pkg/client"
	"github.com/rhacm-global-hub-monitor/backend/pkg/models"
)

// HubHandler handles hub-related requests
type HubHandler struct {
	rhacmClient *client.RHACMClient
}

// NewHubHandler creates a new hub handler
func NewHubHandler(rhacmClient *client.RHACMClient) *HubHandler {
	return &HubHandler{
		rhacmClient: rhacmClient,
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

	hubs, err := h.rhacmClient.GetManagedHubs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

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

	hub, err := h.rhacmClient.GetManagedHub(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

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
