package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/username/dist-ecommerce-go/services/notification-service/internal/service"
)

type EventHandler struct {
	svc service.NotificationService
}

func NewEventHandler(svc service.NotificationService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) HandleOrderCreated(data []byte) error {
	var event struct {
		OrderID string `json:"order_id"`
		UserID  string `json:"user_id"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order created event: %w", err)
	}

	log.Printf("Received OrderCreated event for Order ID: %s", event.OrderID)
	return h.svc.SendEmail(nil, "user@example.com", "Order Received", fmt.Sprintf("Your order %s has been received.", event.OrderID))
}

func (h *EventHandler) HandlePaymentConfirmed(data []byte) error {
	var event struct {
		OrderID string `json:"order_id"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payment confirmed event: %w", err)
	}

	log.Printf("Received PaymentConfirmed event for Order ID: %s", event.OrderID)
	return h.svc.SendEmail(nil, "user@example.com", "Order Confirmed", fmt.Sprintf("Your payment for order %s has been confirmed.", event.OrderID))
}
