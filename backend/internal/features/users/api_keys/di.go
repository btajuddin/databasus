package api_keys

import (
	"sync"
	"sync/atomic"

	users_services "databasus-backend/internal/features/users/services"
	"databasus-backend/internal/util/logger"
)

var (
	setupOnce sync.Once
	isSetup   atomic.Bool
)

var apiKeyRepository = &ApiKeyRepository{}

var apiKeyService = &ApiKeyService{
	apiKeyRepository,
	nil, // auditLogWriter - set via setter in SetupDependencies
}

var apiKeyController = &ApiKeyController{
	apiKeyService,
}

func GetApiKeyService() *ApiKeyService {
	return apiKeyService
}

func GetApiKeyController() *ApiKeyController {
	return apiKeyController
}

func SetupDependencies() {
	wasAlreadySetup := isSetup.Load()

	setupOnce.Do(func() {
		users_services.GetUserService().SetApiKeyValidator(apiKeyService)

		isSetup.Store(true)
	})

	if wasAlreadySetup {
		logger.GetLogger().Warn("SetupDependencies called multiple times, ignoring subsequent call")
	}
}
