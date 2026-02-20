package users_interfaces

import (
	"github.com/google/uuid"

	users_models "databasus-backend/internal/features/users/models"
)

type AuditLogWriter interface {
	WriteAuditLog(message string, userID *uuid.UUID, workspaceID *uuid.UUID)
}

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type ApiKeyValidator interface {
	ValidateApiKey(fullKey string) (*users_models.User, error)
}
