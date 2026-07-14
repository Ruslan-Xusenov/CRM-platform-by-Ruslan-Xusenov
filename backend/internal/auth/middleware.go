package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/crm-platform/backend/internal/shared/middleware"
)

// JWTMiddleware validates JWT tokens and injects user claims into context.
func JWTMiddleware(svc *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				writeError(w, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			claims, err := svc.ValidateAccessToken(parts[1])
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Inject claims into context
			ctx := r.Context()
			ctx = context.WithValue(ctx, middleware.UserIDKey, (*claims)["sub"].(string))
			ctx = context.WithValue(ctx, middleware.TenantIDKey, (*claims)["tenant_id"].(string))
			ctx = context.WithValue(ctx, middleware.RoleKey, (*claims)["role"].(string))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole returns middleware that checks if the user has the required role.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(middleware.RoleKey).(string)
			if !ok {
				writeError(w, http.StatusForbidden, "Access denied")
				return
			}
			for _, allowed := range roles {
				if role == allowed || role == RoleAdmin {
					next.ServeHTTP(w, r)
					return
				}
			}
			writeError(w, http.StatusForbidden, "Insufficient permissions")
		})
	}
}
