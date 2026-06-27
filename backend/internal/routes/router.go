package routes

import (
	"database/sql"
	"net/http"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB, cfg config.Config) http.Handler {
	router := mux.NewRouter()
	router.Use(middleware.CORS)
	router.Use(middleware.JSON)
	router.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(db, cfg)
	categoryHandler := handlers.NewCategoryHandler(db)
	equipmentHandler := handlers.NewEquipmentHandler(db)
	bookingHandler := handlers.NewBookingHandler(db)
	userHandler := handlers.NewUserHandler(db, cfg)
	chatHandler := handlers.NewChatHandler(db, cfg)

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", healthHandler.Check).Methods(http.MethodGet)
	api.HandleFunc("/auth/register", authHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/categories", categoryHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/summary", equipmentHandler.Summary).Methods(http.MethodGet)
	api.HandleFunc("/equipment", equipmentHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/equipment/{id:[0-9]+}", equipmentHandler.Get).Methods(http.MethodGet)
	api.HandleFunc("/profiles/{id:[0-9]+}", userHandler.GetPublicProfile).Methods(http.MethodGet)

	protected := api.NewRoute().Subrouter()
	protected.Use(middleware.Auth(cfg, db))
	protected.HandleFunc("/auth/me", authHandler.Me).Methods(http.MethodGet)
	protected.HandleFunc("/profile", userHandler.UpdateProfile).Methods(http.MethodPut)
	protected.HandleFunc("/bookings", bookingHandler.List).Methods(http.MethodGet)
	protected.HandleFunc("/bookings", bookingHandler.Create).Methods(http.MethodPost)
	protected.HandleFunc("/bookings/owner", bookingHandler.OwnerList).Methods(http.MethodGet)
	protected.HandleFunc("/bookings/{id:[0-9]+}/cancel", bookingHandler.Cancel).Methods(http.MethodPatch)
	protected.HandleFunc("/bookings/{id:[0-9]+}/status", bookingHandler.UpdateStatus).Methods(http.MethodPatch)
	protected.HandleFunc("/equipment", equipmentHandler.Create).Methods(http.MethodPost)
	protected.HandleFunc("/equipment/mine", equipmentHandler.Mine).Methods(http.MethodGet)
	protected.HandleFunc("/equipment/{id:[0-9]+}", equipmentHandler.Update).Methods(http.MethodPut)
	protected.HandleFunc("/chats", chatHandler.List).Methods(http.MethodGet)
	protected.HandleFunc("/chats", chatHandler.Start).Methods(http.MethodPost)
	protected.HandleFunc("/chats/{id:[0-9]+}/messages", chatHandler.Messages).Methods(http.MethodGet)
	protected.HandleFunc("/chats/{id:[0-9]+}/messages", chatHandler.Send).Methods(http.MethodPost)

	moderation := protected.NewRoute().Subrouter()
	moderation.Use(middleware.RequireAdminOrModerator)
	moderation.HandleFunc("/moderation/equipment", equipmentHandler.ModerationList).Methods(http.MethodGet)
	moderation.HandleFunc("/moderation/equipment/{id:[0-9]+}", equipmentHandler.Moderate).Methods(http.MethodPatch)

	admin := protected.NewRoute().Subrouter()
	admin.Use(middleware.RequireAdmin)
	admin.HandleFunc("/categories", categoryHandler.Create).Methods(http.MethodPost)
	admin.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.Update).Methods(http.MethodPut)
	admin.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.Delete).Methods(http.MethodDelete)
	admin.HandleFunc("/equipment/{id:[0-9]+}", equipmentHandler.Delete).Methods(http.MethodDelete)
	admin.HandleFunc("/equipment/{id:[0-9]+}/visibility", equipmentHandler.SetHidden).Methods(http.MethodPatch)
	admin.HandleFunc("/users", userHandler.List).Methods(http.MethodGet)
	admin.HandleFunc("/users/{id:[0-9]+}/role", userHandler.UpdateRole).Methods(http.MethodPatch)
	admin.HandleFunc("/users/{id:[0-9]+}/blocked", userHandler.SetBlocked).Methods(http.MethodPatch)

	return router
}
