package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend/internal/middleware"
	"backend/internal/models"
)

type EquipmentHandler struct {
	db *sql.DB
}

type equipmentRequest struct {
	CategoryID  int64  `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Serial      string `json:"serial"`
	ImageURL    string `json:"image_url"`
	Location    string `json:"location"`
	PricePerDay int64  `json:"price_per_day"`
	Available   *bool  `json:"available"`
}

type moderationRequest struct {
	Status models.ModerationStatus `json:"status"`
	Reason string                  `json:"reason"`
}

func NewEquipmentHandler(db *sql.DB) *EquipmentHandler {
	return &EquipmentHandler{db: db}
}

func (h *EquipmentHandler) Summary(w http.ResponseWriter, r *http.Request) {
	var summary struct {
		EquipmentTotal  int `json:"equipment_total"`
		AvailableTotal  int `json:"available_total"`
		BusyTotal       int `json:"busy_total"`
		ActiveBookings  int `json:"active_bookings"`
		PendingBookings int `json:"pending_bookings"`
	}

	err := h.db.QueryRowContext(
		r.Context(),
		`SELECT
			(SELECT count(*) FROM equipment WHERE NOT hidden AND moderation_status = 'approved'),
			(SELECT count(*) FROM equipment WHERE NOT hidden AND moderation_status = 'approved' AND available),
			(SELECT count(*) FROM equipment WHERE NOT hidden AND moderation_status = 'approved' AND NOT available),
			(SELECT count(*) FROM bookings WHERE status = 'approved'),
			(SELECT count(*) FROM bookings WHERE status = 'pending')`,
	).Scan(&summary.EquipmentTotal, &summary.AvailableTotal, &summary.BusyTotal, &summary.ActiveBookings, &summary.PendingBookings)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to load summary")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, summary)
}

func (h *EquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(
		r.Context(),
		`SELECT e.id, e.owner_id, e.category_id, e.name, e.description, e.serial, e.image_url, e.location, e.price_per_day, e.available, e.hidden, e.moderation_status, e.reject_reason, e.created_at, e.updated_at,
		        c.id, c.name, c.description, c.created_at, c.updated_at,
		        u.id, u.name, u.role, u.phone, u.city, u.avatar_url, u.bio, u.created_at
		 FROM equipment e
		 JOIN categories c ON c.id = e.category_id
		 JOIN users u ON u.id = e.owner_id
		 WHERE NOT e.hidden AND e.moderation_status = 'approved'
		 ORDER BY e.created_at DESC`,
	)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list equipment")
		return
	}
	defer rows.Close()

	items := []models.Equipment{}
	for rows.Next() {
		var e models.Equipment
		var c models.Category
		var u models.PublicUser
		if err := rows.Scan(
			&e.ID, &e.OwnerID, &e.CategoryID, &e.Name, &e.Description, &e.Serial, &e.ImageURL, &e.Location, &e.PricePerDay, &e.Available, &e.Hidden, &e.Moderation, &e.RejectReason, &e.CreatedAt, &e.UpdatedAt,
			&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
			&u.ID, &u.Name, &u.Role, &u.Phone, &u.City, &u.AvatarURL, &u.Bio, &u.CreatedAt,
		); err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to scan equipment")
			return
		}
		e.Category = &c
		e.Owner = &u
		items = append(items, e)
	}

	middleware.WriteJSON(w, http.StatusOK, items)
}

func (h *EquipmentHandler) Mine(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	h.writeEquipmentList(w, r, `WHERE e.owner_id = $1`, user.ID)
}

func (h *EquipmentHandler) ModerationList(w http.ResponseWriter, r *http.Request) {
	h.writeEquipmentList(w, r, `WHERE e.moderation_status = 'pending'`)
}

func (h *EquipmentHandler) writeEquipmentList(w http.ResponseWriter, r *http.Request, where string, args ...any) {
	rows, err := h.db.QueryContext(
		r.Context(),
		`SELECT e.id, e.owner_id, e.category_id, e.name, e.description, e.serial, e.image_url, e.location, e.price_per_day, e.available, e.hidden, e.moderation_status, e.reject_reason, e.created_at, e.updated_at,
		        c.id, c.name, c.description, c.created_at, c.updated_at,
		        u.id, u.name, u.role, u.phone, u.city, u.avatar_url, u.bio, u.created_at
		 FROM equipment e
		 JOIN categories c ON c.id = e.category_id
		 JOIN users u ON u.id = e.owner_id
		 `+where+`
		 ORDER BY e.created_at DESC`,
		args...,
	)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list equipment")
		return
	}
	defer rows.Close()

	items := []models.Equipment{}
	for rows.Next() {
		var e models.Equipment
		var c models.Category
		var u models.PublicUser
		if err := rows.Scan(
			&e.ID, &e.OwnerID, &e.CategoryID, &e.Name, &e.Description, &e.Serial, &e.ImageURL, &e.Location, &e.PricePerDay, &e.Available, &e.Hidden, &e.Moderation, &e.RejectReason, &e.CreatedAt, &e.UpdatedAt,
			&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
			&u.ID, &u.Name, &u.Role, &u.Phone, &u.City, &u.AvatarURL, &u.Bio, &u.CreatedAt,
		); err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to scan equipment")
			return
		}
		e.Category = &c
		e.Owner = &u
		items = append(items, e)
	}

	middleware.WriteJSON(w, http.StatusOK, items)
}

func (h *EquipmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	equipment, err := h.getByID(r, id)
	if err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to get equipment")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, equipment)
}

func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req equipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" || req.CategoryID <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "name and category_id are required")
		return
	}
	if req.PricePerDay <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "price_per_day is required")
		return
	}
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	available := true
	if req.Available != nil {
		available = *req.Available
	}

	var id int64
	err := h.db.QueryRowContext(
		r.Context(),
		`INSERT INTO equipment (owner_id, category_id, name, description, serial, image_url, location, price_per_day, available, moderation_status, reject_reason)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'pending', '')
		 RETURNING id`,
		user.ID,
		req.CategoryID,
		req.Name,
		req.Description,
		req.Serial,
		req.ImageURL,
		req.Location,
		req.PricePerDay,
		available,
	).Scan(&id)
	if err != nil {
		middleware.WriteError(w, http.StatusConflict, "failed to create equipment")
		return
	}

	equipment, err := h.getByIDForUser(r, id, user.ID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to get equipment")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, equipment)
}

func (h *EquipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req equipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" || req.CategoryID <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "name and category_id are required")
		return
	}
	if req.PricePerDay <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "price_per_day is required")
		return
	}

	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var ownerID int64
	if err := h.db.QueryRowContext(r.Context(), `SELECT owner_id FROM equipment WHERE id = $1`, id).Scan(&ownerID); err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	} else if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to check listing owner")
		return
	}
	if user.Role != models.RoleAdmin && user.Role != models.RoleModerator && ownerID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "equipment access denied")
		return
	}

	available := true
	if req.Available != nil {
		available = *req.Available
	}
	moderationStatus := models.ModerationPending
	if user.Role == models.RoleAdmin || user.Role == models.RoleModerator {
		moderationStatus = models.ModerationApproved
	}

	result, err := h.db.ExecContext(
		r.Context(),
		`UPDATE equipment
		 SET category_id = $1, name = $2, description = $3, serial = $4, image_url = $5, location = $6, price_per_day = $7, available = $8, moderation_status = $9, reject_reason = '', updated_at = now()
		 WHERE id = $10`,
		req.CategoryID,
		req.Name,
		req.Description,
		req.Serial,
		req.ImageURL,
		req.Location,
		req.PricePerDay,
		available,
		moderationStatus,
		id,
	)
	if err != nil {
		middleware.WriteError(w, http.StatusConflict, "failed to update equipment")
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	}

	equipment, err := h.getByIDForUser(r, id, user.ID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to get equipment")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, equipment)
}

func (h *EquipmentHandler) Moderate(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req moderationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Status != models.ModerationApproved && req.Status != models.ModerationRejected {
		middleware.WriteError(w, http.StatusBadRequest, "invalid moderation status")
		return
	}
	if req.Status == models.ModerationRejected && req.Reason == "" {
		middleware.WriteError(w, http.StatusBadRequest, "reject reason is required")
		return
	}
	result, err := h.db.ExecContext(
		r.Context(),
		`UPDATE equipment SET moderation_status = $1, reject_reason = $2, updated_at = now() WHERE id = $3`,
		req.Status,
		req.Reason,
		id,
	)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to moderate listing")
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EquipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	result, err := h.db.ExecContext(r.Context(), `DELETE FROM equipment WHERE id = $1`, id)
	if err != nil {
		middleware.WriteError(w, http.StatusConflict, "equipment has bookings")
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EquipmentHandler) SetHidden(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var req struct {
		Hidden bool `json:"hidden"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	result, err := h.db.ExecContext(r.Context(), `UPDATE equipment SET hidden = $1, updated_at = now() WHERE id = $2`, req.Hidden, id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to change listing visibility")
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EquipmentHandler) getByID(r *http.Request, id int64) (models.Equipment, error) {
	return h.getByIDWithFilter(r, id, `AND e.moderation_status = 'approved'`)
}

func (h *EquipmentHandler) getByIDForUser(r *http.Request, id int64, userID int64) (models.Equipment, error) {
	return h.getByIDWithFilter(r, id, `AND (e.moderation_status = 'approved' OR e.owner_id = $2)`, userID)
}

func (h *EquipmentHandler) getByIDWithFilter(r *http.Request, id int64, filter string, args ...any) (models.Equipment, error) {
	var e models.Equipment
	var c models.Category
	var u models.PublicUser
	params := append([]any{id}, args...)
	err := h.db.QueryRowContext(
		r.Context(),
		`SELECT e.id, e.owner_id, e.category_id, e.name, e.description, e.serial, e.image_url, e.location, e.price_per_day, e.available, e.hidden, e.moderation_status, e.reject_reason, e.created_at, e.updated_at,
		        c.id, c.name, c.description, c.created_at, c.updated_at,
		        u.id, u.name, u.role, u.phone, u.city, u.avatar_url, u.bio, u.created_at
		 FROM equipment e
		 JOIN categories c ON c.id = e.category_id
		 JOIN users u ON u.id = e.owner_id
		 WHERE e.id = $1 AND NOT e.hidden `+filter,
		params...,
	).Scan(
		&e.ID, &e.OwnerID, &e.CategoryID, &e.Name, &e.Description, &e.Serial, &e.ImageURL, &e.Location, &e.PricePerDay, &e.Available, &e.Hidden, &e.Moderation, &e.RejectReason, &e.CreatedAt, &e.UpdatedAt,
		&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
		&u.ID, &u.Name, &u.Role, &u.Phone, &u.City, &u.AvatarURL, &u.Bio, &u.CreatedAt,
	)
	e.Category = &c
	e.Owner = &u
	return e, err
}
