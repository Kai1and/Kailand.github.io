package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "user"

type Claims struct {
	UserID int64           `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type AuthUser struct {
	ID    int64
	Email string
	Role  models.UserRole
}

func Auth(cfg config.Config, db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				WriteError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			tokenValue := strings.TrimPrefix(header, "Bearer ")
			if tokenValue == header || tokenValue == "" {
				WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (any, error) {
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			var blocked bool
			err = db.QueryRowContext(r.Context(), `SELECT blocked FROM users WHERE id = $1`, claims.UserID).Scan(&blocked)
			if err != nil || blocked {
				WriteError(w, http.StatusUnauthorized, "account is blocked or no longer exists")
				return
			}

			user := AuthUser{ID: claims.UserID, Email: claims.Email, Role: claims.Role}
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := CurrentUser(r)
		if !ok || user.Role != models.RoleAdmin {
			WriteError(w, http.StatusForbidden, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdminOrModerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := CurrentUser(r)
		if !ok || (user.Role != models.RoleAdmin && user.Role != models.RoleModerator) {
			WriteError(w, http.StatusForbidden, "moderator access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CurrentUser(r *http.Request) (AuthUser, bool) {
	user, ok := r.Context().Value(userContextKey).(AuthUser)
	return user, ok
}
