package auth

import (
	"time"

	"github.com/google/uuid"
)

// ─── Entities ────────────────────────────────────────────────

type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	TenantID  uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password_hash"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Phone     string     `json:"phone,omitempty" db:"phone"`
	Avatar    string     `json:"avatar,omitempty" db:"avatar_url"`
	Role      string     `json:"role" db:"role"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type Tenant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Domain    string    `json:"domain,omitempty" db:"domain"`
	Plan      string    `json:"plan" db:"plan"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ─── Request/Response DTOs ───────────────────────────────────

type RegisterRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	TenantName string `json:"tenant_name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         User   `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ─── Roles & Permissions ─────────────────────────────────────

const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

// RolePermissions maps roles to their allowed API actions.
var RolePermissions = map[string][]string{
	RoleAdmin:    {"*"},
	RoleManager:  {"leads:*", "contacts:*", "companies:*", "deals:*", "pipelines:*", "calls:*", "activities:*", "files:*", "audit:read"},
	RoleOperator: {"leads:read", "leads:create", "leads:update", "contacts:read", "contacts:create", "deals:read", "calls:*", "activities:*", "files:upload"},
	RoleViewer:   {"leads:read", "contacts:read", "companies:read", "deals:read", "calls:read", "activities:read"},
}
