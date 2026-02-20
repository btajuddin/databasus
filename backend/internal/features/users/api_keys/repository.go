package api_keys

import (
	"time"

	"databasus-backend/internal/features/users/api_keys/models"
	"databasus-backend/internal/storage"

	"github.com/google/uuid"
)

type ApiKeyRepository struct{}

func (r *ApiKeyRepository) Create(apiKey *models.UserApiKey) error {
	if apiKey.ID == uuid.Nil {
		apiKey.ID = uuid.New()
	}

	return storage.GetDb().Create(apiKey).Error
}

func (r *ApiKeyRepository) GetByUserID(userID uuid.UUID) (*models.UserApiKey, error) {
	var apiKey models.UserApiKey

	err := storage.GetDb().
		Where("user_id = ?", userID).
		First(&apiKey).Error

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (r *ApiKeyRepository) GetByKeyPrefix(keyPrefix string) (*models.UserApiKey, error) {
	var apiKey models.UserApiKey

	err := storage.GetDb().
		Where("key_prefix = ?", keyPrefix).
		First(&apiKey).Error

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (r *ApiKeyRepository) UpdateLastUsedAt(id uuid.UUID) error {
	now := time.Now().UTC()

	return storage.GetDb().Model(&models.UserApiKey{}).
		Where("id = ?", id).
		Update("last_used_at", now).Error
}

func (r *ApiKeyRepository) DeleteByUserID(userID uuid.UUID) error {
	return storage.GetDb().
		Where("user_id = ?", userID).
		Delete(&models.UserApiKey{}).Error
}

func (r *ApiKeyRepository) Delete(id uuid.UUID) error {
	return storage.GetDb().
		Where("id = ?", id).
		Delete(&models.UserApiKey{}).Error
}
