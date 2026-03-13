package handlers

import (
	"context"

	pb "github.com/username/dist-ecommerce-go/proto/order"
	"github.com/username/dist-ecommerce-go/services/order-service/internal/models"
	"github.com/username/dist-ecommerce-go/services/order-service/internal/service"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	var items []models.OrderItem
	for _, item := range req.Items {
		items = append(items, models.OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	order, err := h.svc.CreateOrder(ctx, req.UserId, items)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Order: mapOrderToPb(order),
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	order, err := h.svc.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Order: mapOrderToPb(order),
	}, nil
}

func mapOrderToPb(order *models.Order) *pb.Order {
	var items []*pb.OrderItem
	for _, item := range order.Items {
		items = append(items, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	return &pb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     string(order.Status),
		Items:      items,
		CreatedAt:  order.CreatedAt.String(),
	}
}
