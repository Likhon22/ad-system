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
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         Role
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
