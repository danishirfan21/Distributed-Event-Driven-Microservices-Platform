package repository

import (
	"context"

	"github.com/username/dist-ecommerce-go/services/inventory-service/internal/models"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetByID(ctx context.Context, id string) (*models.Product, error)
	UpdateStock(ctx context.Context, id string, quantity int32) error
	RecordHistory(ctx context.Context, history *models.StockHistory) error
	IsEventProcessed(ctx context.Context, eventID, serviceName string) (bool, error)
	MarkEventProcessed(ctx context.Context, eventID, serviceName string) error
}

type postgresInventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &postgresInventoryRepository{db: db}
}

func (r *postgresInventoryRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *postgresInventoryRepository) UpdateStock(ctx context.Context, id string, quantity int32) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).UpdateColumn("stock_quantity", gorm.Expr("stock_quantity + ?", quantity)).Error
}

func (r *postgresInventoryRepository) RecordHistory(ctx context.Context, history *models.StockHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *postgresInventoryRepository) IsEventProcessed(ctx context.Context, eventID, serviceName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ProcessedEvent{}).Where("id = ? AND service = ?", eventID, serviceName).Count(&count).Error
	return count > 0, err
}

func (r *postgresInventoryRepository) MarkEventProcessed(ctx context.Context, eventID, serviceName string) error {
	return r.db.WithContext(ctx).Create(&models.ProcessedEvent{ID: eventID, Service: serviceName}).Error
}
