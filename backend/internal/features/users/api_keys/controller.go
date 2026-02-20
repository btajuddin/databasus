package api_keys

import (
	"net/http"

	users_middleware "databasus-backend/internal/features/users/middleware"

	"github.com/gin-gonic/gin"
)

type ApiKeyController struct {
	apiKeyService *ApiKeyService
}

func (c *ApiKeyController) RegisterRoutes(router *gin.RouterGroup) {
	routes := router.Group("/users/me")
	routes.POST("/api-key", c.UpsertApiKey)
	routes.GET("/api-key", c.GetApiKey)
}

// UpsertApiKey
// @Summary Create or regenerate API key for current user
// @Description Creates a new API key or regenerates an existing one for the authenticated user. The key is only shown once.
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CreateApiKeyResponseDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/api-key [post]
func (c *ApiKeyController) UpsertApiKey(ctx *gin.Context) {
	user, isOk := users_middleware.GetUserFromContext(ctx)
	if !isOk {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	response, err := c.apiKeyService.UpsertApiKey(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetApiKey
// @Summary Get API key info for current user
// @Description Returns information about the current user's API key (without the actual key).
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiKeyInfoDTO
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/me/api-key [get]
func (c *ApiKeyController) GetApiKey(ctx *gin.Context) {
	user, isOk := users_middleware.GetUserFromContext(ctx)
	if !isOk {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	response, err := c.apiKeyService.GetApiKeyInfo(user.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
