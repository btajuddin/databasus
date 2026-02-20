package api_keys

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"databasus-backend/internal/features/users/api_keys/dto"
	"databasus-backend/internal/features/users/api_keys/models"
	users_interfaces "databasus-backend/internal/features/users/interfaces"
	users_models "databasus-backend/internal/features/users/models"
	users_repositories "databasus-backend/internal/features/users/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	apiKeyPrefix       = "dbs_live_"
	apiKeyPrefixLength = 17
	apiKeySecretLength = 44
)

type ApiKeyService struct {
	repository     *ApiKeyRepository
	userRepository *users_repositories.UserRepository
	auditLogWriter users_interfaces.AuditLogWriter
}

func (s *ApiKeyService) SetAuditLogWriter(writer users_interfaces.AuditLogWriter) {
	s.auditLogWriter = writer
}

func (s *ApiKeyService) GenerateApiKey(userID uuid.UUID) (*dto.CreateApiKeyResponseDTO, error) {
	existingKey, err := s.repository.GetByUserID(userID)
	if err == nil && existingKey != nil {
		return nil, errors.New("user already has an API key")
	}

	secret, err := s.generateRandomSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	fullKey := apiKeyPrefix + secret

	hashedKey, err := bcrypt.GenerateFromPassword([]byte(fullKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash API key: %w", err)
	}

	uniquePrefix := fullKey[:apiKeyPrefixLength]

	apiKey := &models.UserApiKey{
		UserID:    userID,
		HashedKey: string(hashedKey),
		KeyPrefix: uniquePrefix,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repository.Create(apiKey); err != nil {
		return nil, fmt.Errorf("failed to store API key: %w", err)
	}

	s.writeAuditLog("API key generated", &userID)

	return &dto.CreateApiKeyResponseDTO{
		ApiKey:    fullKey,
		KeyPrefix: uniquePrefix,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}

func (s *ApiKeyService) GetApiKeyInfo(userID uuid.UUID) (*dto.ApiKeyInfoDTO, error) {
	apiKey, err := s.repository.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("API key not found")
	}

	return &dto.ApiKeyInfoDTO{
		ID:         apiKey.ID,
		KeyPrefix:  apiKey.KeyPrefix,
		CreatedAt:  apiKey.CreatedAt,
		LastUsedAt: apiKey.LastUsedAt,
	}, nil
}

func (s *ApiKeyService) RegenerateApiKey(userID uuid.UUID) (*dto.RegenerateApiKeyResponseDTO, error) {
	if err := s.repository.DeleteByUserID(userID); err != nil {
		return nil, fmt.Errorf("failed to delete existing API key: %w", err)
	}

	response, err := s.GenerateApiKey(userID)
	if err != nil {
		return nil, err
	}

	return &dto.RegenerateApiKeyResponseDTO{
		ApiKey:    response.ApiKey,
		KeyPrefix: response.KeyPrefix,
		CreatedAt: response.CreatedAt,
	}, nil
}

func (s *ApiKeyService) ValidateApiKey(fullKey string) (*users_models.User, error) {
	if len(fullKey) <= apiKeyPrefixLength {
		return nil, errors.New("invalid API key format")
	}

	prefix := fullKey[:len(apiKeyPrefix)]
	if prefix != apiKeyPrefix {
		return nil, errors.New("invalid API key format")
	}

	uniquePrefix := fullKey[:apiKeyPrefixLength]

	apiKey, err := s.repository.GetByKeyPrefix(uniquePrefix)
	if err != nil {
		return nil, errors.New("API key not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(apiKey.HashedKey), []byte(fullKey)); err != nil {
		return nil, errors.New("invalid API key")
	}

	if err := s.repository.UpdateLastUsedAt(apiKey.ID); err != nil {
		return nil, fmt.Errorf("failed to update last used timestamp: %w", err)
	}

	user, err := s.userRepository.GetUserByID(apiKey.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *ApiKeyService) generateRandomSecret() (string, error) {
	bytes := make([]byte, apiKeySecretLength/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (s *ApiKeyService) writeAuditLog(message string, userID *uuid.UUID) {
	if s.auditLogWriter != nil {
		s.auditLogWriter.WriteAuditLog(message, userID, nil)
	}
}
