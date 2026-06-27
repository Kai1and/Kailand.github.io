package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/middleware"
	"backend/internal/models"
)

type BookingHandler struct {
	db *sql.DB
}

type bookingRequest struct {
	EquipmentID int64  `json:"equipment_id"`
	StartAt     string `json:"start_at"`
	EndAt       string `json:"end_at"`
	Comment     string `json:"comment"`
}

type bookingStatusRequest struct {
	Status models.BookingStatus `json:"status"`
}

func NewBookingHandler(db *sql.DB) *BookingHandler {
	return &BookingHandler{db: db}
}

func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	query := bookingListQuery()
	args := []any{}
	if user.Role != models.RoleAdmin {
		query += ` WHERE b.user_id = $1`
		args = append(args, user.ID)
	}
	query += ` ORDER BY b.created_at DESC`
	h.writeBookingRows(w, r, query, args...)
}

func (h *BookingHandler) OwnerList(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	query := bookingListQuery() + ` WHERE e.owner_id = $1 ORDER BY b.created_at DESC`
	h.writeBookingRows(w, r, query, user.ID)
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req bookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	startAt, endAt, ok := parseBookingTime(w, req.StartAt, req.EndAt)
	if !ok {
		return
	}
	if req.EquipmentID <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "equipment_id is required")
		return
	}

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to start booking")
		return
	}
	defer tx.Rollback()

	var ownerID int64
	var available bool
	var hidden bool
	var moderation models.ModerationStatus
	err = tx.QueryRowContext(
		r.Context(),
		`SELECT owner_id, available, hidden, moderation_status FROM equipment WHERE id = $1 FOR UPDATE`,
		req.EquipmentID,
	).Scan(&ownerID, &available, &hidden, &moderation)
	if err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to check equipment")
		return
	}
	if hidden || moderation != models.ModerationApproved {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	}
	if ownerID == user.ID {
		middleware.WriteError(w, http.StatusBadRequest, "you cannot book your own equipment")
		return
	}
	if !available {
		middleware.WriteError(w, http.StatusConflict, "equipment is already in stop list")
		return
	}

	var conflicts int
	err = tx.QueryRowContext(
		r.Context(),
		`SELECT count(*)
		 FROM bookings
		 WHERE equipment_id = $1
		   AND status IN ('pending', 'approved')
		   AND tstzrange(start_at, end_at, '[)') && tstzrange($2, $3, '[)')`,
		req.EquipmentID,
		startAt,
		endAt,
	).Scan(&conflicts)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to check booking conflicts")
		return
	}
	if conflicts > 0 {
		middleware.WriteError(w, http.StatusConflict, "equipment is already booked for this period")
		return
	}

	booking := models.Booking{}
	err = tx.QueryRowContext(
		r.Context(),
		`INSERT INTO bookings (user_id, equipment_id, start_at, end_at, status, comment)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, equipment_id, start_at, end_at, status, comment, created_at, updated_at`,
		user.ID,
		req.EquipmentID,
		startAt,
		endAt,
		models.BookingPending,
		req.Comment,
	).Scan(&booking.ID, &booking.UserID, &booking.EquipmentID, &booking.StartAt, &booking.EndAt, &booking.Status, &booking.Comment, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to create booking")
		return
	}

	if _, err := tx.ExecContext(r.Context(), `UPDATE equipment SET available = false, updated_at = now() WHERE id = $1`, req.EquipmentID); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to move equipment to stop list")
		return
	}
	if err := tx.Commit(); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to save booking")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, booking)
}

func (h *BookingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req bookingStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if !validBookingStatus(req.Status) {
		middleware.WriteError(w, http.StatusBadRequest, "invalid status")
		return
	}

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update booking")
		return
	}
	defer tx.Rollback()

	booking := models.Booking{}
	var ownerID int64
	err = tx.QueryRowContext(
		r.Context(),
		`SELECT b.id, b.user_id, b.equipment_id, b.start_at, b.end_at, b.status, b.comment, b.created_at, b.updated_at, e.owner_id
		 FROM bookings b
		 JOIN equipment e ON e.id = b.equipment_id
		 WHERE b.id = $1
		 FOR UPDATE`,
		id,
	).Scan(&booking.ID, &booking.UserID, &booking.EquipmentID, &booking.StartAt, &booking.EndAt, &booking.Status, &booking.Comment, &booking.CreatedAt, &booking.UpdatedAt, &ownerID)
	if err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "booking not found")
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to get booking")
		return
	}
	if user.Role != models.RoleAdmin && ownerID != user.ID {
		middleware.WriteError(w, http.StatusForbidden, "booking access denied")
		return
	}

	err = tx.QueryRowContext(
		r.Context(),
		`UPDATE bookings
		 SET status = $1, updated_at = now()
		 WHERE id = $2
		 RETURNING id, user_id, equipment_id, start_at, end_at, status, comment, created_at, updated_at`,
		req.Status,
		id,
	).Scan(&booking.ID, &booking.UserID, &booking.EquipmentID, &booking.StartAt, &booking.EndAt, &booking.Status, &booking.Comment, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update booking")
		return
	}
	if err := h.refreshEquipmentAvailability(r, tx, booking.EquipmentID); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to refresh equipment availability")
		return
	}
	if err := tx.Commit(); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to save booking")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, booking)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to cancel booking")
		return
	}
	defer tx.Rollback()

	var equipmentID int64
	query := `UPDATE bookings SET status = $1, updated_at = now() WHERE id = $2`
	args := []any{models.BookingCancelled, id}
	if user.Role != models.RoleAdmin {
		query += ` AND user_id = $3`
		args = append(args, user.ID)
	}
	query += ` RETURNING equipment_id`
	err = tx.QueryRowContext(r.Context(), query, args...).Scan(&equipmentID)
	if err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "booking not found")
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to cancel booking")
		return
	}
	if err := h.refreshEquipmentAvailability(r, tx, equipmentID); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to refresh equipment availability")
		return
	}
	if err := tx.Commit(); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to save booking")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookingHandler) writeBookingRows(w http.ResponseWriter, r *http.Request, query string, args ...any) {
	rows, err := h.db.QueryContext(r.Context(), query, args...)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list bookings")
		return
	}
	defer rows.Close()

	bookings := []models.Booking{}
	for rows.Next() {
		var b models.Booking
		var e models.Equipment
		var customer models.User
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.EquipmentID, &b.StartAt, &b.EndAt, &b.Status, &b.Comment, &b.CreatedAt, &b.UpdatedAt,
			&e.ID, &e.OwnerID, &e.CategoryID, &e.Name, &e.Description, &e.Serial, &e.ImageURL, &e.Location, &e.PricePerDay, &e.Available, &e.Hidden, &e.Moderation, &e.RejectReason, &e.CreatedAt, &e.UpdatedAt,
			&customer.ID, &customer.Name, &customer.Email, &customer.Role, &customer.Phone, &customer.City, &customer.AvatarURL, &customer.Bio, &customer.Blocked, &customer.CreatedAt, &customer.UpdatedAt,
		); err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to scan booking")
			return
		}
		b.Equipment = &e
		b.User = &customer
		bookings = append(bookings, b)
	}

	middleware.WriteJSON(w, http.StatusOK, bookings)
}

func bookingListQuery() string {
	return `SELECT b.id, b.user_id, b.equipment_id, b.start_at, b.end_at, b.status, b.comment, b.created_at, b.updated_at,
	        e.id, e.owner_id, e.category_id, e.name, e.description, e.serial, e.image_url, e.location, e.price_per_day, e.available, e.hidden, e.moderation_status, e.reject_reason, e.created_at, e.updated_at,
	        u.id, u.name, u.email, u.role, u.phone, u.city, u.avatar_url, u.bio, u.blocked, u.created_at, u.updated_at
	 FROM bookings b
	 JOIN equipment e ON e.id = b.equipment_id
	 JOIN users u ON u.id = b.user_id`
}

func (h *BookingHandler) refreshEquipmentAvailability(r *http.Request, tx *sql.Tx, equipmentID int64) error {
	var active int
	if err := tx.QueryRowContext(
		r.Context(),
		`SELECT count(*) FROM bookings WHERE equipment_id = $1 AND status IN ('pending', 'approved')`,
		equipmentID,
	).Scan(&active); err != nil {
		return err
	}
	_, err := tx.ExecContext(r.Context(), `UPDATE equipment SET available = $1, updated_at = now() WHERE id = $2`, active == 0, equipmentID)
	return err
}

func parseBookingTime(w http.ResponseWriter, startRaw, endRaw string) (time.Time, time.Time, bool) {
	startAt, err := time.Parse(time.RFC3339, startRaw)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "start_at must be RFC3339")
		return time.Time{}, time.Time{}, false
	}
	endAt, err := time.Parse(time.RFC3339, endRaw)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "end_at must be RFC3339")
		return time.Time{}, time.Time{}, false
	}
	if !endAt.After(startAt) {
		middleware.WriteError(w, http.StatusBadRequest, "end_at must be after start_at")
		return time.Time{}, time.Time{}, false
	}
	return startAt, endAt, true
}

func validBookingStatus(status models.BookingStatus) bool {
	switch status {
	case models.BookingPending, models.BookingApproved, models.BookingRejected, models.BookingCancelled, models.BookingReturned:
		return true
	default:
		return false
	}
}
