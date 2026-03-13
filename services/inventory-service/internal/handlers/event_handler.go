package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/username/dist-ecommerce-go/services/inventory-service/internal/service"
)

type EventHandler struct {
	svc service.InventoryService
}

func NewEventHandler(svc service.InventoryService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) HandlePaymentConfirmed(data []byte) error {
	var event struct {
		OrderID string `json:"order_id"`
		// In a real system, the event would include the items to deduct
		// For this simulation, we'll assume we need to fetch the order or have the items in the event
		// Let's simplify and assume the event has what we need or we mock it.
		Items []struct {
			ProductID string `json:"product_id"`
			Quantity  int32  `json:"quantity"`
		} `json:"items"`
	}

	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payment confirmed event: %w", err)
	}

	log.Printf("Received PaymentConfirmed event for Order ID: %s", event.OrderID)

	// Idempotency check
	processed, err := h.svc.IsProcessed(nil, event.OrderID)
	if err != nil {
		return err
	}
	if processed {
		log.Printf("Order %s already processed by inventory service, skipping", event.OrderID)
		return nil
	}

	if len(event.Items) == 0 {
		log.Printf("No items found in event for Order ID %s, skipping stock deduction", event.OrderID)
		return nil
	}

	for _, item := range event.Items {
		err := h.svc.DeductStock(nil, item.ProductID, item.Quantity, event.OrderID)
		if err != nil {
			log.Printf("Failed to deduct stock for product %s: %v", item.ProductID, err)
		}
	}

	return h.svc.MarkProcessed(nil, event.OrderID)
}
