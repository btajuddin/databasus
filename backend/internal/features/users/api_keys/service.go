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

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	apiKeyPrefix       = "databasus_api_"
	apiKeyPrefixLength = 22
	apiKeySecretLength = 44
)

type ApiKeyService struct {
	repository     *ApiKeyRepository
	auditLogWriter users_interfaces.AuditLogWriter
}

func (s *ApiKeyService) SetAuditLogWriter(writer users_interfaces.AuditLogWriter) {
	s.auditLogWriter = writer
}

func (s *ApiKeyService) UpsertApiKey(userID uuid.UUID) (*dto.CreateApiKeyResponseDTO, error) {
	s.repository.DeleteByUserID(userID)

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

func (s *ApiKeyService) ValidateApiKey(fullKey string) (uuid.UUID, error) {
	if len(fullKey) <= apiKeyPrefixLength {
		return uuid.Nil, errors.New("invalid API key format")
	}

	prefix := fullKey[:len(apiKeyPrefix)]
	if prefix != apiKeyPrefix {
		return uuid.Nil, errors.New("invalid API key format")
	}

	uniquePrefix := fullKey[:apiKeyPrefixLength]

	apiKey, err := s.repository.GetByKeyPrefix(uniquePrefix)
	if err != nil {
		return uuid.Nil, errors.New("API key not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(apiKey.HashedKey), []byte(fullKey)); err != nil {
		return uuid.Nil, errors.New("invalid API key")
	}

	if err := s.repository.UpdateLastUsedAt(apiKey.ID); err != nil {
		return uuid.Nil, fmt.Errorf("failed to update last used timestamp: %w", err)
	}

	return apiKey.UserID, nil
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
