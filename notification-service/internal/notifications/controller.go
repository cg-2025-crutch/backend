package notifications

import (
	"encoding/json"
	"net/http"

	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/notification-service/internal/models"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications/service"
	"github.com/mailru/easyjson"
)

type Controller struct {
	service *service.NotificationService
}

func NewController(service *service.NotificationService) *Controller {
	return &Controller{service: service}
}

// GetVapidKeyHandler returns the VAPID public key for web push notifications
func (c *Controller) GetVapidKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := log.FromContext(ctx)

	if r.Method != http.MethodGet {
		l.Error("method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vapidKey := c.service.GetVapidKey()

	response := map[string]string{
		"publicKey": vapidKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		l.Error("failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// SubscribeHandler handles subscription requests from clients
func (c *Controller) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := log.FromContext(ctx)

	if r.Method != http.MethodPost {
		l.Error("method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SubscriptionReq
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		l.Error("failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Extract user ID from context or headers
	// For now, using a header for user identification
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		l.Error("user ID not provided in headers")
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Endpoint == "" || req.Keys.P256dh == "" || req.Keys.Auth == "" {
		l.Error("missing required subscription fields")
		http.Error(w, "Missing required fields: endpoint, keys.p256dh, keys.auth", http.StatusBadRequest)
		return
	}

	l.Infof("Subscribing user: %s, p256dh length: %d, auth length: %d",
		userID, len(req.Keys.P256dh), len(req.Keys.Auth))

	err := c.service.SubscribeUser(ctx, userID, req.Endpoint, req.Keys.P256dh, req.Keys.Auth)
	if err != nil {
		l.Error("failed to subscribe user", err)
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Subscription successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		l.Error("failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
