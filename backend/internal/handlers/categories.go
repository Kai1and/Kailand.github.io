package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	db *sql.DB
}

type categoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(), `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list categories")
		return
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to scan category")
			return
		}
		categories = append(categories, c)
	}

	middleware.WriteJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" {
		middleware.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	category := models.Category{}
	err := h.db.QueryRowContext(
		r.Context(),
		`INSERT INTO categories (name, description)
		 VALUES ($1, $2)
		 RETURNING id, name, description, created_at, updated_at`,
		req.Name,
		req.Description,
	).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		middleware.WriteError(w, http.StatusConflict, "category already exists")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" {
		middleware.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	category := models.Category{}
	err := h.db.QueryRowContext(
		r.Context(),
		`UPDATE categories
		 SET name = $1, description = $2, updated_at = now()
		 WHERE id = $3
		 RETURNING id, name, description, created_at, updated_at`,
		req.Name,
		req.Description,
		id,
	).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "category not found")
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusConflict, "failed to update category")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	result, err := h.db.ExecContext(r.Context(), `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		middleware.WriteError(w, http.StatusConflict, "category is in use")
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		middleware.WriteError(w, http.StatusNotFound, "category not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || id <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}
