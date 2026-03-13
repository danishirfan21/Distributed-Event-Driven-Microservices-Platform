package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/username/dist-ecommerce-go/pkg/events"
	"github.com/username/dist-ecommerce-go/services/payment-service/internal/models"
	"github.com/username/dist-ecommerce-go/services/payment-service/internal/repository"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, orderID string, amount float64) (*models.Payment, error)
}

type PaymentServiceWithItems struct {
	repo     repository.PaymentRepository
	eventBus *events.EventBus
}

func NewPaymentService(repo repository.PaymentRepository, eventBus *events.EventBus) *PaymentServiceWithItems {
	return &PaymentServiceWithItems{repo: repo, eventBus: eventBus}
}

func (s *PaymentServiceWithItems) ProcessPayment(ctx context.Context, orderID string, amount float64) (*models.Payment, error) {
	return s.ProcessPaymentWithItems(ctx, orderID, amount, nil)
}

func (s *PaymentServiceWithItems) ProcessPaymentWithItems(ctx context.Context, orderID string, amount float64, items interface{}) (*models.Payment, error) {
	// Simulate payment processing
	time.Sleep(100 * time.Millisecond)

	payment := &models.Payment{
		OrderID:               orderID,
		Amount:                amount,
		Status:                "CONFIRMED",
		ProviderTransactionID: fmt.Sprintf("txn_%s", uuid.New().String()),
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Publish event
	event := map[string]interface{}{
		"order_id":       payment.OrderID,
		"payment_id":     payment.ID,
		"transaction_id": payment.ProviderTransactionID,
		"status":         payment.Status,
		"items":          items,
	}
	_ = s.eventBus.Publish(events.PaymentConfirmed, event)

	return payment, nil
}
