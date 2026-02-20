package api_keys

import (
	"net/http"
	"testing"

	"databasus-backend/internal/features/users/api_keys/dto"
	users_controllers "databasus-backend/internal/features/users/controllers"
	users_dto "databasus-backend/internal/features/users/dto"
	users_enums "databasus-backend/internal/features/users/enums"
	users_middleware "databasus-backend/internal/features/users/middleware"
	users_services "databasus-backend/internal/features/users/services"
	users_testing "databasus-backend/internal/features/users/testing"
	test_utils "databasus-backend/internal/util/testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UpsertApiKey_WhenUserHasNoKey_ApiKeyCreated(t *testing.T) {
	router := createApiKeyTestRouter()
	user := users_testing.CreateTestUser(users_enums.UserRoleMember)

	var response dto.CreateApiKeyResponseDTO
	test_utils.MakePostRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		nil,
		http.StatusOK,
		&response,
	)

	assert.NotEmpty(t, response.ApiKey)
	assert.NotEmpty(t, response.KeyPrefix)
	assert.Contains(t, response.ApiKey, "databasus_api_")
	assert.Equal(t, "databasus_api_", response.KeyPrefix[:14])
}

func Test_UpsertApiKey_WhenUserAlreadyHasKey_ReplacesWithNewKey(t *testing.T) {
	router := createApiKeyTestRouter()
	user := users_testing.CreateTestUser(users_enums.UserRoleMember)

	var firstResponse dto.CreateApiKeyResponseDTO
	test_utils.MakePostRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		nil,
		http.StatusOK,
		&firstResponse,
	)

	var secondResponse dto.CreateApiKeyResponseDTO
	test_utils.MakePostRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		nil,
		http.StatusOK,
		&secondResponse,
	)

	assert.NotEmpty(t, secondResponse.ApiKey)
	assert.NotEqual(t, firstResponse.ApiKey, secondResponse.ApiKey)
	assert.NotEqual(t, firstResponse.KeyPrefix, secondResponse.KeyPrefix)
}

func Test_GetApiKey_WhenKeyExists_ReturnsKeyInfo(t *testing.T) {
	router := createApiKeyTestRouter()
	user := users_testing.CreateTestUser(users_enums.UserRoleMember)

	var createResponse dto.CreateApiKeyResponseDTO
	test_utils.MakePostRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		nil,
		http.StatusOK,
		&createResponse,
	)

	var infoResponse dto.ApiKeyInfoDTO
	test_utils.MakeGetRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		http.StatusOK,
		&infoResponse,
	)

	assert.Equal(t, createResponse.KeyPrefix, infoResponse.KeyPrefix)
	assert.NotEmpty(t, infoResponse.ID)
}

func Test_GetApiKey_WhenNoKeyExists_ReturnsNotFound(t *testing.T) {
	router := createApiKeyTestRouter()
	user := users_testing.CreateTestUser(users_enums.UserRoleMember)

	resp := test_utils.MakeGetRequest(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		http.StatusNotFound,
	)

	assert.Contains(t, string(resp.Body), "API key not found")
}

func Test_AuthWithApiKey_ValidKey_AuthenticatesUser(t *testing.T) {
	router := createApiKeyTestRouter()
	user := users_testing.CreateTestUser(users_enums.UserRoleMember)

	var createResponse dto.CreateApiKeyResponseDTO
	test_utils.MakePostRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		nil,
		http.StatusOK,
		&createResponse,
	)

	var profile users_dto.UserProfileResponseDTO
	test_utils.MakeGetRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me",
		"Bearer "+createResponse.ApiKey,
		http.StatusOK,
		&profile,
	)

	assert.Equal(t, user.UserID, profile.ID)
	assert.Equal(t, user.Email, profile.Email)
}

func Test_AuthWithApiKey_InvalidKey_ReturnsUnauthorized(t *testing.T) {
	router := createApiKeyTestRouter()

	resp := test_utils.MakeGetRequest(
		t,
		router,
		"/api/v1/users/me",
		"Bearer invalid_api_key_format",
		http.StatusUnauthorized,
	)

	assert.Contains(t, string(resp.Body), "User not authenticated")
}

func Test_AuthWithApiKey_AfterUpsert_OldKeyStopsWorking(t *testing.T) {
	router := createApiKeyTestRouter()
	user := users_testing.CreateTestUser(users_enums.UserRoleMember)

	var createResponse dto.CreateApiKeyResponseDTO
	test_utils.MakePostRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		nil,
		http.StatusOK,
		&createResponse,
	)

	oldApiKey := createResponse.ApiKey

	test_utils.MakePostRequestAndUnmarshal(
		t,
		router,
		"/api/v1/users/me/api-key",
		"Bearer "+user.Token,
		nil,
		http.StatusOK,
		&createResponse,
	)

	resp := test_utils.MakeGetRequest(
		t,
		router,
		"/api/v1/users/me",
		"Bearer "+oldApiKey,
		http.StatusUnauthorized,
	)

	assert.Contains(t, string(resp.Body), "User not authenticated")
}

func createApiKeyTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/api/v1")

	protected := v1.Group("").Use(users_middleware.AuthMiddleware(users_services.GetUserService()))

	GetApiKeyController().RegisterRoutes(protected.(*gin.RouterGroup))
	users_controllers.GetUserController().RegisterProtectedRoutes(protected.(*gin.RouterGroup))

	users_services.GetUserService().SetAuditLogWriter(&auditLogWriterStub{})

	SetupDependencies()

	return router
}

type auditLogWriterStub struct{}

func (a *auditLogWriterStub) WriteAuditLog(message string, userID *uuid.UUID, workspaceID *uuid.UUID) {
}
