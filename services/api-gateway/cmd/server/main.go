package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/username/dist-ecommerce-go/pkg/metrics"
	"github.com/username/dist-ecommerce-go/pkg/tracing"
	orderpb "github.com/username/dist-ecommerce-go/proto/order"
	userpb "github.com/username/dist-ecommerce-go/proto/user"
	"github.com/username/dist-ecommerce-go/services/api-gateway/internal/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serviceName := "api-gateway"

	// Tracing
	collectorAddr := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if collectorAddr != "" {
		tp, err := tracing.InitTracer(serviceName, collectorAddr)
		if err != nil {
			log.Printf("failed to initialize tracer: %v", err)
		} else {
			defer tp.Shutdown(nil)
		}
	}

	// Metrics
	go func() {
		if err := metrics.StartMetricsServer(":8081"); err != nil {
			log.Printf("failed to start metrics server: %v", err)
		}
	}()

	// gRPC Clients
	userSvcAddr := os.Getenv("USER_SERVICE_ADDR")
	if userSvcAddr == "" {
		userSvcAddr = "localhost:50051"
	}
	userConn, err := grpc.Dial(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	orderSvcAddr := os.Getenv("ORDER_SERVICE_ADDR")
	if orderSvcAddr == "" {
		orderSvcAddr = "localhost:50052"
	}
	orderConn, err := grpc.Dial(orderSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to order service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	// HTTP Server
	handler := handlers.NewGatewayHandler(userClient, orderClient)
	r := gin.Default()

	r.POST("/users", handler.CreateUser)
	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders/:id", handler.GetOrder)

	log.Printf("API Gateway listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}
