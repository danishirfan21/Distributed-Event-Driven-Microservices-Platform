package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/username/dist-ecommerce-go/services/payment-service/internal/service"
)

type EventHandler struct {
	svc *service.PaymentServiceWithItems
}

func NewEventHandler(svc *service.PaymentServiceWithItems) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) HandleOrderCreated(data []byte) error {
	var event struct {
		OrderID    string  `json:"order_id"`
		TotalPrice float64 `json:"total_price"`
		Items      interface{} `json:"items"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order created event: %w", err)
	}

	log.Printf("Received OrderCreated event for Order ID: %s", event.OrderID)

	// Process payment asynchronously
	// In a real system, we'd pass the items through to the payment result event
	_, err := h.svc.ProcessPaymentWithItems(nil, event.OrderID, event.TotalPrice, event.Items)
	if err != nil {
		log.Printf("Failed to process payment for Order ID %s: %v", event.OrderID, err)
		return err
	}

	return nil
}
