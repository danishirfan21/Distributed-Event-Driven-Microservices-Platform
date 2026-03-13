package main

import (
	"log"
	"net"
	"os"

	"github.com/username/dist-ecommerce-go/pkg/common"
	"github.com/username/dist-ecommerce-go/pkg/database"
	"github.com/username/dist-ecommerce-go/pkg/metrics"
	"github.com/username/dist-ecommerce-go/pkg/tracing"
	pb "github.com/username/dist-ecommerce-go/proto/user"
	"github.com/username/dist-ecommerce-go/services/user-service/internal/handlers"
	"github.com/username/dist-ecommerce-go/services/user-service/internal/models"
	"github.com/username/dist-ecommerce-go/services/user-service/internal/repository"
	"github.com/username/dist-ecommerce-go/services/user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	serviceName := "user-service"

	// Database
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=users port=5432 sslmode=disable"
	}
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Migrate schema
	db.AutoMigrate(&models.User{})

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

	// Cache
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	cache := common.NewCache(redisAddr)

	// Repository & Service
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, cache)
	handler := handlers.NewUserHandler(svc)

	// gRPC Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, handler)
	reflection.Register(s)

	log.Printf("User Service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
