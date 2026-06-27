package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db  *sql.DB
	cfg config.Config
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var namePattern = regexp.MustCompile(`^[\p{L}][\p{L}\s'-]{1,58}[\p{L}]$`)

type authResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func NewAuthHandler(db *sql.DB, cfg config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if !namePattern.MatchString(req.Name) {
		middleware.WriteError(w, http.StatusBadRequest, "name must contain letters, not numbers")
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "email must be valid")
		return
	}
	if len(req.Password) < 6 {
		middleware.WriteError(w, http.StatusBadRequest, "password must contain at least 6 chars")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := models.User{}
	err = h.db.QueryRowContext(
		r.Context(),
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, email, password_hash, role, phone, city, avatar_url, bio, blocked, created_at, updated_at`,
		req.Name,
		req.Email,
		string(hash),
		models.RoleUser,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Phone, &user.City, &user.AvatarURL, &user.Bio, &user.Blocked, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		middleware.WriteError(w, http.StatusConflict, "user already exists")
		return
	}

	token, err := h.issueToken(user)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, authResponse{Token: token, User: user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	user, err := h.findUserByEmail(r, req.Email)
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if user.Blocked {
		middleware.WriteError(w, http.StatusForbidden, "account is blocked")
		return
	}

	token, err := h.issueToken(user)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, authResponse{Token: token, User: user})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	current, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user := models.User{}
	err := h.db.QueryRowContext(
		r.Context(),
		`SELECT id, name, email, password_hash, role, phone, city, avatar_url, bio, blocked, created_at, updated_at FROM users WHERE id = $1`,
		current.ID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Phone, &user.City, &user.AvatarURL, &user.Bio, &user.Blocked, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) findUserByEmail(r *http.Request, email string) (models.User, error) {
	if email == "" {
		return models.User{}, errors.New("empty email")
	}

	user := models.User{}
	err := h.db.QueryRowContext(
		r.Context(),
		`SELECT id, name, email, password_hash, role, phone, city, avatar_url, bio, blocked, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Phone, &user.City, &user.AvatarURL, &user.Bio, &user.Blocked, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (h *AuthHandler) issueToken(user models.User) (string, error) {
	now := time.Now()
	claims := middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(h.cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.Email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}
