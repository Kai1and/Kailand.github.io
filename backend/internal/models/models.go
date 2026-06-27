package models

import "time"

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleAdmin     UserRole = "admin"
	RoleModerator UserRole = "moderator"
)

type ModerationStatus string

const (
	ModerationPending  ModerationStatus = "pending"
	ModerationApproved ModerationStatus = "approved"
	ModerationRejected ModerationStatus = "rejected"
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingApproved  BookingStatus = "approved"
	BookingRejected  BookingStatus = "rejected"
	BookingCancelled BookingStatus = "cancelled"
	BookingReturned  BookingStatus = "returned"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	Phone        string    `json:"phone"`
	City         string    `json:"city"`
	AvatarURL    string    `json:"avatar_url"`
	Bio          string    `json:"bio"`
	Blocked      bool      `json:"blocked"`
	BanReason    string    `json:"ban_reason,omitempty"`
	BanEvidence  string    `json:"ban_evidence,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PublicUser struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Role      UserRole  `json:"role"`
	Phone     string    `json:"phone,omitempty"`
	City      string    `json:"city"`
	AvatarURL string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Equipment struct {
	ID           int64            `json:"id"`
	OwnerID      int64            `json:"owner_id"`
	Owner        *PublicUser      `json:"owner,omitempty"`
	CategoryID   int64            `json:"category_id"`
	Category     *Category        `json:"category,omitempty"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Serial       string           `json:"serial"`
	ImageURL     string           `json:"image_url"`
	Location     string           `json:"location"`
	PricePerDay  int64            `json:"price_per_day"`
	Available    bool             `json:"available"`
	Hidden       bool             `json:"hidden"`
	Moderation   ModerationStatus `json:"moderation_status"`
	RejectReason string           `json:"reject_reason,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type Booking struct {
	ID          int64         `json:"id"`
	UserID      int64         `json:"user_id"`
	EquipmentID int64         `json:"equipment_id"`
	User        *User         `json:"user,omitempty"`
	Equipment   *Equipment    `json:"equipment,omitempty"`
	StartAt     time.Time     `json:"start_at"`
	EndAt       time.Time     `json:"end_at"`
	Status      BookingStatus `json:"status"`
	Comment     string        `json:"comment"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type Conversation struct {
	ID          int64       `json:"id"`
	EquipmentID int64       `json:"equipment_id"`
	CustomerID  int64       `json:"customer_id"`
	OwnerID     int64       `json:"owner_id"`
	Equipment   *Equipment  `json:"equipment,omitempty"`
	Customer    *PublicUser `json:"customer,omitempty"`
	Owner       *PublicUser `json:"owner,omitempty"`
	LastMessage string      `json:"last_message"`
	UnreadCount int         `json:"unread_count"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Message struct {
	ID             int64      `json:"id"`
	ConversationID int64      `json:"conversation_id"`
	SenderID       int64      `json:"sender_id"`
	Body           string     `json:"body"`
	AttachmentURL  string     `json:"attachment_url"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
