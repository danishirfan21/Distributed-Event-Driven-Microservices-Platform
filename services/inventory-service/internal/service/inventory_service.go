package service

import (
	"context"
	"fmt"

	"github.com/username/dist-ecommerce-go/pkg/events"
	"github.com/username/dist-ecommerce-go/services/inventory-service/internal/models"
	"github.com/username/dist-ecommerce-go/services/inventory-service/internal/repository"
)

type InventoryService interface {
	DeductStock(ctx context.Context, productID string, quantity int32, orderID string) error
	GetProduct(ctx context.Context, id string) (*models.Product, error)
	IsProcessed(ctx context.Context, orderID string) (bool, error)
	MarkProcessed(ctx context.Context, orderID string) error
}

type inventoryService struct {
	repo     repository.InventoryRepository
	eventBus *events.EventBus
}

func NewInventoryService(repo repository.InventoryRepository, eventBus *events.EventBus) InventoryService {
	return &inventoryService{repo: repo, eventBus: eventBus}
}

func (s *inventoryService) DeductStock(ctx context.Context, productID string, quantity int32, orderID string) error {
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	if product.StockQuantity < quantity {
		return fmt.Errorf("insufficient stock for product %s", productID)
	}

	if err := s.repo.UpdateStock(ctx, productID, -quantity); err != nil {
		return err
	}

	history := &models.StockHistory{
		ProductID: productID,
		Change:    -quantity,
		Reason:    fmt.Sprintf("Order %s", orderID),
	}
	_ = s.repo.RecordHistory(ctx, history)

	// Publish event
	event := map[string]interface{}{
		"product_id": productID,
		"new_stock":  product.StockQuantity - quantity,
	}
	_ = s.eventBus.Publish(events.InventoryUpdated, event)

	return nil
}

func (s *inventoryService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *inventoryService) IsProcessed(ctx context.Context, orderID string) (bool, error) {
	return s.repo.IsEventProcessed(ctx, orderID, "inventory-service")
}

func (s *inventoryService) MarkProcessed(ctx context.Context, orderID string) error {
	return s.repo.MarkEventProcessed(ctx, orderID, "inventory-service")
}
