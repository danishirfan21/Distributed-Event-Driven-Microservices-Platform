package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/username/dist-ecommerce-go/pkg/events"
	"github.com/username/dist-ecommerce-go/pkg/metrics"
	"github.com/username/dist-ecommerce-go/pkg/tracing"
	"github.com/username/dist-ecommerce-go/services/notification-service/internal/handlers"
	"github.com/username/dist-ecommerce-go/services/notification-service/internal/service"
)

func main() {
	serviceName := "notification-service"

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

	// Service
	svc := service.NewNotificationService()
	eventHandler := handlers.NewEventHandler(svc)

	// Subscribe to events
	_, err = eb.Subscribe(events.OrderCreated, "notification-service", eventHandler.HandleOrderCreated)
	if err != nil {
		log.Fatalf("failed to subscribe to OrderCreated: %v", err)
	}

	_, err = eb.Subscribe(events.PaymentConfirmed, "notification-service", eventHandler.HandlePaymentConfirmed)
	if err != nil {
		log.Fatalf("failed to subscribe to PaymentConfirmed: %v", err)
	}

	log.Printf("Notification Service started and subscribing to events")

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Printf("Shutting down Notification Service")
}
