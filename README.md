# Distributed E-commerce Microservices Platform

This project is a production-grade distributed microservices platform built with Golang, following an event-driven architecture.

## Architecture

The system consists of several independent microservices communicating via gRPC (synchronous) and NATS JetStream (asynchronous events).

![Architecture Diagram](https://raw.githubusercontent.com/username/dist-ecommerce-go/main/docs/architecture.png) *Note: Placeholder for actual diagram*

### Services

*   **API Gateway**: The entry point for all clients. Exposes REST APIs and routes them to internal gRPC services.
*   **User Service**: Manages user profiles and authentication.
*   **Order Service**: Handles the lifecycle of an order.
*   **Payment Service**: Processes payments and reacts to new orders.
*   **Inventory Service**: Manages product stock and updates inventory based on confirmed payments.
*   **Notification Service**: Sends user notifications (emails/SMS) based on system events.

### Tech Stack

*   **Language**: Golang
*   **Communication**: gRPC, REST (via Gin)
*   **Event Bus**: NATS JetStream
*   **Databases**: Postgres (GORM)
*   **Observability**: OpenTelemetry (Tracing), Prometheus (Metrics)
*   **Infrastructure**: Docker, Kubernetes, Terraform

## Getting Started

### Prerequisites

*   Go 1.25+
*   Docker & Docker Compose
*   `protoc` and Go plugins (if modifying protos)

### Running Locally

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/username/dist-ecommerce-go.git
    cd dist-ecommerce-go
    ```

2.  **Start all services**:
    ```bash
    docker-compose up --build
    ```

3.  **Access the services**:
    *   API Gateway: `http://localhost:8080`
    *   Prometheus: `http://localhost:9090`
    *   Jaeger UI: `http://localhost:16686`

### Example Workflow

1.  **Create a User**:
    ```bash
    curl -X POST http://localhost:8080/users -d '{
      "email": "john@example.com",
      "password": "securepassword",
      "full_name": "John Doe"
    }'
    ```

2.  **Place an Order**:
    ```bash
    curl -X POST http://localhost:8080/orders -d '{
      "user_id": "USER_ID_FROM_PREVIOUS_STEP",
      "items": [
        {"product_id": "PROD-123", "quantity": 2, "unit_price": 49.99}
      ]
    }'
    ```

3.  **Watch the Event Flow**:
    Check the logs of `payment-service`, `inventory-service`, and `notification-service` to see them reacting to the events.

## Deployment

### Kubernetes

Manifests are located in `deployments/k8s/`.

```bash
kubectl apply -f deployments/k8s/
```

### Infrastructure (Terraform)

Infrastructure templates are in `deployments/terraform/`.

```bash
cd deployments/terraform
terraform init
terraform apply
```

## Observability

*   **Tracing**: Each service is instrumented with OpenTelemetry. Traces are sent to Jaeger.
*   **Metrics**: Prometheus scrapes `/metrics` endpoints on port 8081 for each service.
