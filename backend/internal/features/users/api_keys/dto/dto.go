package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateApiKeyResponseDTO struct {
	ApiKey    string    `json:"apiKey"`
	KeyPrefix string    `json:"keyPrefix"`
	CreatedAt time.Time `json:"createdAt"`
}

type ApiKeyInfoDTO struct {
	ID         uuid.UUID  `json:"id"`
	KeyPrefix  string     `json:"keyPrefix"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
}

type RegenerateApiKeyResponseDTO struct {
	ApiKey    string    `json:"apiKey"`
	KeyPrefix string    `json:"keyPrefix"`
	CreatedAt time.Time `json:"createdAt"`
}
