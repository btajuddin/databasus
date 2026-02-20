package models

import (
	"time"

	"github.com/google/uuid"
)

type UserApiKey struct {
	ID         uuid.UUID  `json:"id"         gorm:"column:id"`
	UserID     uuid.UUID  `json:"userId"     gorm:"column:user_id"`
	HashedKey  string     `json:"-"          gorm:"column:hashed_key"`
	KeyPrefix  string     `json:"keyPrefix"  gorm:"column:key_prefix"`
	CreatedAt  time.Time  `json:"createdAt"  gorm:"column:created_at"`
	LastUsedAt *time.Time `json:"lastUsedAt" gorm:"column:last_used_at"`
}

func (UserApiKey) TableName() string {
	return "user_api_keys"
}
