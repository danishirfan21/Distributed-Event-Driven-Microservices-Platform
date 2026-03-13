package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/username/dist-ecommerce-go/pkg/database"
	"github.com/username/dist-ecommerce-go/pkg/events"
	"github.com/username/dist-ecommerce-go/pkg/metrics"
	"github.com/username/dist-ecommerce-go/pkg/tracing"
	"github.com/username/dist-ecommerce-go/services/payment-service/internal/handlers"
	"github.com/username/dist-ecommerce-go/services/payment-service/internal/models"
	"github.com/username/dist-ecommerce-go/services/payment-service/internal/repository"
	"github.com/username/dist-ecommerce-go/services/payment-service/internal/service"
)

func main() {
	serviceName := "payment-service"

	// Database
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=payments port=5432 sslmode=disable"
	}
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Migrate schema
	db.AutoMigrate(&models.Payment{})

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
	repo := repository.NewPaymentRepository(db)
	svc := service.NewPaymentService(repo, eb)
	eventHandler := handlers.NewEventHandler(svc)

	// Subscribe to events
	_, err = eb.Subscribe(events.OrderCreated, "payment-service", eventHandler.HandleOrderCreated)
	if err != nil {
		log.Fatalf("failed to subscribe to OrderCreated: %v", err)
	}

	log.Printf("Payment Service started and subscribing to events")

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Printf("Shutting down Payment Service")
}
