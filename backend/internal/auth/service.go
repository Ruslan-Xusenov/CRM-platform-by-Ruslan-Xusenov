package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/crm-platform/backend/pkg/cache"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db    *pgxpool.Pool
	redis *cache.RedisClient
}

func NewService(db *pgxpool.Pool, redis *cache.RedisClient) *Service {
	return &Service{db: db, redis: redis}
}

// Register creates a new tenant and admin user.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*TokenResponse, error) {
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	tenantID := uuid.New()
	userID := uuid.New()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create tenant
	_, err = tx.Exec(ctx,
		`INSERT INTO tenants (id, name, plan, is_active, created_at) VALUES ($1, $2, 'free', true, NOW())`,
		tenantID, req.TenantName)
	if err != nil {
		return nil, fmt.Errorf("create tenant: %w", err)
	}

	// Create admin user
	_, err = tx.Exec(ctx,
		`INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, role, is_active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, true, NOW())`,
		userID, tenantID, req.Email, string(hash), req.FirstName, req.LastName, RoleAdmin)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	user := User{
		ID: userID, TenantID: tenantID, Email: req.Email,
		FirstName: req.FirstName, LastName: req.LastName,
		Role: RoleAdmin, IsActive: true, CreatedAt: time.Now(),
	}

	return s.generateTokens(ctx, user)
}

// Login authenticates a user and returns tokens.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	var user User
	var passwordHash string

	err := s.db.QueryRow(ctx,
		`SELECT id, tenant_id, email, password_hash, first_name, last_name, role, is_active, created_at
		 FROM users WHERE email = $1 AND is_active = true`, req.Email).
		Scan(&user.ID, &user.TenantID, &user.Email, &passwordHash,
			&user.FirstName, &user.LastName, &user.Role, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return s.generateTokens(ctx, user)
}

// RefreshToken validates a refresh token and issues new tokens.
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	tokenHash := hashToken(refreshToken)

	// Validate refresh token exists and is not revoked
	var userID uuid.UUID
	var expiresAt time.Time
	err := s.db.QueryRow(ctx,
		`SELECT user_id, expires_at FROM refresh_tokens
		 WHERE token_hash = $1 AND revoked = false AND expires_at > NOW()`,
		tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Revoke old token
	s.db.Exec(ctx, `UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`, tokenHash)

	// Fetch user
	var user User
	err = s.db.QueryRow(ctx,
		`SELECT id, tenant_id, email, first_name, last_name, role, is_active, created_at
		 FROM users WHERE id = $1 AND is_active = true`, userID).
		Scan(&user.ID, &user.TenantID, &user.Email, &user.FirstName, &user.LastName,
			&user.Role, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.generateTokens(ctx, user)
}

// Logout revokes a refresh token.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	_, err := s.db.Exec(ctx, `UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`, tokenHash)
	return err
}

// ValidateAccessToken parses and validates a JWT access token.
func (s *Service) ValidateAccessToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(getEnvOrDefault("JWT_ACCESS_SECRET", "dev-access-secret")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return &claims, nil
}

// generateTokens creates access + refresh token pair.
func (s *Service) generateTokens(ctx context.Context, user User) (*TokenResponse, error) {
	accessSecret := getEnvOrDefault("JWT_ACCESS_SECRET", "dev-access-secret")
	now := time.Now()
	accessExp := now.Add(15 * time.Minute)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID.String(),
		"tenant_id": user.TenantID.String(),
		"email":     user.Email,
		"role":      user.Role,
		"iat":       now.Unix(),
		"exp":       accessExp.Unix(),
	})
	accessString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// Generate refresh token (opaque)
	refreshString := uuid.NewString() + "-" + uuid.NewString()
	refreshHash := hashToken(refreshString)
	refreshExp := now.Add(7 * 24 * time.Hour)

	// Store refresh token hash in DB
	_, err = s.db.Exec(ctx,
		`INSERT INTO refresh_tokens (token_hash, user_id, expires_at, revoked) VALUES ($1, $2, $3, false)`,
		refreshHash, user.ID, refreshExp)
	if err != nil {
		slog.Error("Failed to store refresh token", "error", err)
	}

	return &TokenResponse{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		ExpiresIn:    int64(15 * 60),
		User:         user,
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
