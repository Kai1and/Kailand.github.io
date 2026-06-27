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

type ChatHandler struct {
	db  *sql.DB
	cfg config.Config
}

type startConversationRequest struct {
	EquipmentID int64  `json:"equipment_id"`
	Message     string `json:"message"`
}

type sendMessageRequest struct {
	Body          string `json:"body"`
	AttachmentURL string `json:"attachment_url"`
}

func NewChatHandler(db *sql.DB, cfg config.Config) *ChatHandler {
	return &ChatHandler{db: db, cfg: cfg}
}

func (h *ChatHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rows, err := h.db.QueryContext(
		r.Context(),
		`SELECT c.id, c.equipment_id, c.customer_id, c.owner_id, c.last_message, c.created_at, c.updated_at,
		        e.id, e.owner_id, e.category_id, e.name, e.description, e.serial, e.image_url, e.location, e.price_per_day, e.available, e.hidden, e.created_at, e.updated_at,
		        customer.id, customer.name, customer.role, customer.phone, customer.city, customer.avatar_url, customer.bio, customer.created_at,
		        owner.id, owner.name, owner.role, owner.phone, owner.city, owner.avatar_url, owner.bio, owner.created_at,
		        (SELECT count(*) FROM messages m WHERE m.conversation_id = c.id AND m.sender_id <> $1 AND m.read_at IS NULL)
		 FROM conversations c
		 JOIN equipment e ON e.id = c.equipment_id
		 JOIN users customer ON customer.id = c.customer_id
		 JOIN users owner ON owner.id = c.owner_id
		 WHERE c.customer_id = $1 OR c.owner_id = $1
		 ORDER BY c.updated_at DESC`,
		user.ID,
	)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list conversations")
		return
	}
	defer rows.Close()

	conversations := []models.Conversation{}
	for rows.Next() {
		var c models.Conversation
		var e models.Equipment
		var customer models.PublicUser
		var owner models.PublicUser
		if err := rows.Scan(
			&c.ID, &c.EquipmentID, &c.CustomerID, &c.OwnerID, &c.LastMessage, &c.CreatedAt, &c.UpdatedAt,
			&e.ID, &e.OwnerID, &e.CategoryID, &e.Name, &e.Description, &e.Serial, &e.ImageURL, &e.Location, &e.PricePerDay, &e.Available, &e.Hidden, &e.CreatedAt, &e.UpdatedAt,
			&customer.ID, &customer.Name, &customer.Role, &customer.Phone, &customer.City, &customer.AvatarURL, &customer.Bio, &customer.CreatedAt,
			&owner.ID, &owner.Name, &owner.Role, &owner.Phone, &owner.City, &owner.AvatarURL, &owner.Bio, &owner.CreatedAt,
			&c.UnreadCount,
		); err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to scan conversation")
			return
		}
		c.LastMessage, _ = security.DecryptString(h.cfg.DataEncryptionKey, c.LastMessage)
		c.Equipment = &e
		c.Customer = &customer
		c.Owner = &owner
		conversations = append(conversations, c)
	}

	middleware.WriteJSON(w, http.StatusOK, conversations)
}

func (h *ChatHandler) Start(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req startConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.EquipmentID <= 0 {
		middleware.WriteError(w, http.StatusBadRequest, "equipment_id is required")
		return
	}

	var ownerID int64
	err := h.db.QueryRowContext(r.Context(), `SELECT owner_id FROM equipment WHERE id = $1 AND moderation_status = 'approved'`, req.EquipmentID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		middleware.WriteError(w, http.StatusNotFound, "equipment not found")
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to get equipment owner")
		return
	}
	if ownerID == user.ID {
		middleware.WriteError(w, http.StatusBadRequest, "you cannot start chat with yourself")
		return
	}

	var conversationID int64
	storedMessage := ""
	if req.Message != "" {
		storedMessage, err = security.EncryptString(h.cfg.DataEncryptionKey, req.Message)
		if err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to protect message")
			return
		}
	}
	err = h.db.QueryRowContext(
		r.Context(),
		`INSERT INTO conversations (equipment_id, customer_id, owner_id, last_message)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (equipment_id, customer_id, owner_id)
		 DO UPDATE SET updated_at = now()
		 RETURNING id`,
		req.EquipmentID,
		user.ID,
		ownerID,
		storedMessage,
	).Scan(&conversationID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to start conversation")
		return
	}

	if req.Message != "" {
		if _, err := h.insertMessage(r, conversationID, user.ID, req.Message); err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to send message")
			return
		}
	}

	middleware.WriteJSON(w, http.StatusCreated, map[string]int64{"id": conversationID})
}

func (h *ChatHandler) Messages(w http.ResponseWriter, r *http.Request) {
	conversationID, ok := parseID(w, r)
	if !ok {
		return
	}
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if !h.canAccessConversation(r, conversationID, user.ID) {
		middleware.WriteError(w, http.StatusForbidden, "conversation access denied")
		return
	}
	_, _ = h.db.ExecContext(r.Context(), `UPDATE messages SET read_at = now() WHERE conversation_id = $1 AND sender_id <> $2 AND read_at IS NULL`, conversationID, user.ID)

	rows, err := h.db.QueryContext(
		r.Context(),
		`SELECT id, conversation_id, sender_id, body, attachment_url, read_at, created_at
		 FROM messages
		 WHERE conversation_id = $1
		 ORDER BY created_at ASC`,
		conversationID,
	)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list messages")
		return
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Body, &m.AttachmentURL, &m.ReadAt, &m.CreatedAt); err != nil {
			middleware.WriteError(w, http.StatusInternalServerError, "failed to scan message")
			return
		}
		m.Body, _ = security.DecryptString(h.cfg.DataEncryptionKey, m.Body)
		messages = append(messages, m)
	}

	middleware.WriteJSON(w, http.StatusOK, messages)
}

func (h *ChatHandler) Send(w http.ResponseWriter, r *http.Request) {
	conversationID, ok := parseID(w, r)
	if !ok {
		return
	}
	user, ok := middleware.CurrentUser(r)
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if !h.canAccessConversation(r, conversationID, user.ID) {
		middleware.WriteError(w, http.StatusForbidden, "conversation access denied")
		return
	}

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Body == "" && req.AttachmentURL == "" {
		middleware.WriteError(w, http.StatusBadRequest, "message body or attachment is required")
		return
	}

	message, err := h.insertMessage(r, conversationID, user.ID, req.Body, req.AttachmentURL)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to send message")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, message)
}

func (h *ChatHandler) insertMessage(r *http.Request, conversationID, senderID int64, body string, attachmentURL ...string) (models.Message, error) {
	message := models.Message{}
	storedBody, err := security.EncryptString(h.cfg.DataEncryptionKey, body)
	if err != nil {
		return message, err
	}
	attachment := ""
	if len(attachmentURL) > 0 {
		attachment = attachmentURL[0]
	}
	lastMessage := storedBody
	if body == "" && attachment != "" {
		lastMessage, _ = security.EncryptString(h.cfg.DataEncryptionKey, "Фото")
	}
	err = h.db.QueryRowContext(
		r.Context(),
		`WITH inserted AS (
			INSERT INTO messages (conversation_id, sender_id, body, attachment_url)
			VALUES ($1, $2, $3, $4)
			RETURNING id, conversation_id, sender_id, body, attachment_url, read_at, created_at
		)
		UPDATE conversations
		SET last_message = $5, updated_at = now()
		WHERE id = $1
		RETURNING (SELECT id FROM inserted), (SELECT conversation_id FROM inserted), (SELECT sender_id FROM inserted), (SELECT body FROM inserted), (SELECT attachment_url FROM inserted), (SELECT read_at FROM inserted), (SELECT created_at FROM inserted)`,
		conversationID,
		senderID,
		storedBody,
		attachment,
		lastMessage,
	).Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.Body, &message.AttachmentURL, &message.ReadAt, &message.CreatedAt)
	if err == nil {
		message.Body = body
	}
	return message, err
}

func (h *ChatHandler) canAccessConversation(r *http.Request, conversationID, userID int64) bool {
	var exists bool
	err := h.db.QueryRowContext(
		r.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM conversations
			WHERE id = $1 AND (customer_id = $2 OR owner_id = $2)
		)`,
		conversationID,
		userID,
	).Scan(&exists)
	return err == nil && exists
}
