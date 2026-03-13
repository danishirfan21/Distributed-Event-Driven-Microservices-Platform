package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	userpb "github.com/username/dist-ecommerce-go/proto/user"
	orderpb "github.com/username/dist-ecommerce-go/proto/order"
)

type GatewayHandler struct {
	userClient  userpb.UserServiceClient
	orderClient orderpb.OrderServiceClient
}

func NewGatewayHandler(userClient userpb.UserServiceClient, orderClient orderpb.OrderServiceClient) *GatewayHandler {
	return &GatewayHandler{
		userClient:  userClient,
		orderClient: orderClient,
	}
}

func (h *GatewayHandler) CreateUser(c *gin.Context) {
	var req userpb.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.userClient.CreateUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *GatewayHandler) CreateOrder(c *gin.Context) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.orderClient.CreateOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *GatewayHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	res, err := h.orderClient.GetOrder(context.Background(), &orderpb.GetOrderRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
