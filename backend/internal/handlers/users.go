package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/security"
)

type UserHandler struct {
	db  *sql.DB
	cfg config.Config
}

type userRoleRequest struct {
	Role models.UserRole `json:"role"`
}

type blockedRequest struct {
	Blocked  bool   `json:"blocked"`
	Reason   string `json:"reason"`
	Evidence string `json:"evidence"`
}

type profileRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	City      string `json:"city"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

func NewUserHandler(db *sql.DB, cfg config.Config) *UserHandler {
	return &UserHandler{db: db, cfg: cfg}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(), `SELECT id, name, email, password_hash, role, phone, city, avatar_url, bio, blocked, ban_reason, ban_evidence, created_at, updated_at FROM users ORDER BY id`)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.Phone, &u.City, &u.AvatarURL, &u.Bio, &u.Blocked, &u.BanReason, &u.BanEvidence, &u.CreatedAt, &u.UpdatedAt); err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to scan user")
			return
		}
		u.BanEvidence, _ = security.DecryptString(h.cfg.DataEncryptionKey, u.BanEvidence)
		users = append(users, u)
	}

	middleware.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req userRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Role != models.RoleUser && req.Role != models.RoleAdmin && req.Role != models.RoleModerator {
		middleware.WriteError(w, http.StatusBadRequest, "invalid role")
		return
	}

	user := models.User{}
	err := h.db.QueryRowContext(
		r.Context(),
		`UPDATE users
		 SET role = $1, updated_at = now()
		 WHERE id = $2
		 RETURNING id, name, email, password_hash, role, phone, city, avatar_url, bio, blocked, ban_reason, ban_evidence, created_at, updated_at`,
		req.Role,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Phone, &user.City, &user.AvatarURL, &user.Bio, &user.Blocked, &user.BanReason, &user.BanEvidence, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	user := models.PublicUser{}
	err := h.db.QueryRowContext(
		r.Context(),
		`SELECT id, name, role, phone, city, avatar_url, bio, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Role, &user.Phone, &user.City, &user.AvatarURL, &user.Bio, &user.CreatedAt)
	if err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "profile not found")
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to get profile")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	current, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req profileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" {
		middleware.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	user := models.User{}
	err := h.db.QueryRowContext(
		r.Context(),
		`UPDATE users
		 SET name = $1, phone = $2, city = $3, avatar_url = $4, bio = $5, updated_at = now()
		 WHERE id = $6
		 RETURNING id, name, email, password_hash, role, phone, city, avatar_url, bio, blocked, ban_reason, ban_evidence, created_at, updated_at`,
		req.Name,
		req.Phone,
		req.City,
		req.AvatarURL,
		req.Bio,
		current.ID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Phone, &user.City, &user.AvatarURL, &user.Bio, &user.Blocked, &user.BanReason, &user.BanEvidence, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) SetBlocked(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	current, ok := middleware.CurrentUser(r)
	if !ok || current.ID == id {
		middleware.WriteError(w, http.StatusBadRequest, "cannot change your own access")
		return
	}

	var req blockedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Blocked && req.Reason == "" {
		middleware.WriteError(w, http.StatusBadRequest, "ban reason is required")
		return
	}

	evidence := req.Evidence
	if req.Blocked {
		encrypted, err := security.EncryptString(h.cfg.DataEncryptionKey, req.Evidence)
		if err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to protect ban evidence")
			return
		}
		evidence = encrypted
	} else {
		req.Reason = ""
		evidence = ""
	}

	result, err := h.db.ExecContext(
		r.Context(),
		`UPDATE users SET blocked = $1, ban_reason = $2, ban_evidence = $3, updated_at = now() WHERE id = $4`,
		req.Blocked,
		req.Reason,
		evidence,
		id,
	)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update user access")
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		middleware.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
