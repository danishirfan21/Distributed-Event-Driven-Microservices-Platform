package service

import (
	"context"

	"github.com/username/dist-ecommerce-go/pkg/events"
	"github.com/username/dist-ecommerce-go/services/order-service/internal/models"
	"github.com/username/dist-ecommerce-go/services/order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items []models.OrderItem) (*models.Order, error)
	GetOrder(ctx context.Context, id string) (*models.Order, error)
}

type orderService struct {
	repo     repository.OrderRepository
	eventBus *events.EventBus
}

func NewOrderService(repo repository.OrderRepository, eventBus *events.EventBus) OrderService {
	return &orderService{repo: repo, eventBus: eventBus}
}

func (s *orderService) CreateOrder(ctx context.Context, userID string, items []models.OrderItem) (*models.Order, error) {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}

	order := &models.Order{
		UserID:     userID,
		TotalPrice: total,
		Status:     models.StatusPending,
		Items:      items,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Publish event
	var eventItems []map[string]interface{}
	for _, item := range order.Items {
		eventItems = append(eventItems, map[string]interface{}{
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
			"unit_price": item.UnitPrice,
		})
	}

	event := map[string]interface{}{
		"order_id":    order.ID,
		"user_id":     order.UserID,
		"total_price": order.TotalPrice,
		"items":       eventItems,
	}
	_ = s.eventBus.Publish(events.OrderCreated, event)

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.GetByID(ctx, id)
}
