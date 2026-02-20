package users_interfaces

import "github.com/google/uuid"

type AuditLogWriter interface {
	WriteAuditLog(message string, userID *uuid.UUID, workspaceID *uuid.UUID)
}

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type ApiKeyValidator interface {
	ValidateApiKey(fullKey string) (uuid.UUID, error)
}
