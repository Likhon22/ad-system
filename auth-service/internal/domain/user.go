// internal/domain/user.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string
type Status string

const (
	RoleAdvertiser Role = "advertiser"
	RolePublisher  Role = "publisher"
	RoleAdmin      Role = "admin"
)

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name,omitempty"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Provider     string    `json:"provider"`
	ProviderID   string    `json:"-"`
	Role         Role      `json:"role"`
	Status       Status    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
