package main

import (
	"log"
	"net"
	"os"

	"github.com/username/dist-ecommerce-go/pkg/database"
	"github.com/username/dist-ecommerce-go/pkg/events"
	"github.com/username/dist-ecommerce-go/pkg/metrics"
	"github.com/username/dist-ecommerce-go/pkg/tracing"
	pb "github.com/username/dist-ecommerce-go/proto/order"
	"github.com/username/dist-ecommerce-go/services/order-service/internal/handlers"
	"github.com/username/dist-ecommerce-go/services/order-service/internal/models"
	"github.com/username/dist-ecommerce-go/services/order-service/internal/repository"
	"github.com/username/dist-ecommerce-go/services/order-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	serviceName := "order-service"

	// Database
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=orders port=5432 sslmode=disable"
	}
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Migrate schema
	db.AutoMigrate(&models.Order{}, &models.OrderItem{})

	// NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	eb, err := events.NewEventBus(natsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer eb.Close()

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

	// Repository & Service
	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo, eb)
	handler := handlers.NewOrderHandler(svc)

	// gRPC Server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, handler)
	reflection.Register(s)

	log.Printf("Order Service listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
