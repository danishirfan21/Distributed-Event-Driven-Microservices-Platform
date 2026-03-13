# System Design: Distributed E-commerce Microservices

## Architecture Overview

The system follows an event-driven microservices architecture designed for scalability, loose coupling, and maintainability.

### Services

1.  **API Gateway**:
    *   Entry point for all client requests.
    *   Handles RESTful API requests.
    *   Proxies requests to internal services using gRPC.
    *   Handles authentication and rate limiting.

2.  **User Service**:
    *   Manages user accounts, profiles, and authentication.
    *   Database: Postgres.

3.  **Order Service**:
    *   Handles order creation and management.
    *   Orchestrates the initial step of the order workflow.
    *   Publishes `OrderCreated` events.
    *   Database: Postgres.

4.  **Payment Service**:
    *   Processes payments for orders.
    *   Subscribes to `OrderCreated` events.
    *   Publishes `PaymentConfirmed` or `PaymentFailed` events.
    *   Database: Postgres.

5.  **Inventory Service**:
    *   Manages product catalog and stock levels.
    *   Subscribes to `PaymentConfirmed` events to reserve/deduct stock.
    *   Database: Postgres.

6.  **Notification Service**:
    *   Sends notifications (email/SMS - mocked) to users.
    *   Subscribes to various events (`OrderCreated`, `PaymentConfirmed`).

## Communication Patterns

### Synchronous (gRPC)
Used for request-response interactions where immediate feedback is required (e.g., API Gateway to User Service for authentication, API Gateway to Order Service for order submission).

### Asynchronous (NATS JetStream)
Used for event-driven workflows. Services emit events when state changes occur, and other services react accordingly. This ensures loose coupling and high availability.

## Event Workflow: Order Placement

1.  **Client** sends `POST /orders` to **API Gateway**.
2.  **API Gateway** calls **Order Service** via gRPC.
3.  **Order Service** creates a new order in `PENDING` state and persists it to Postgres.
4.  **Order Service** publishes an `order.created` event to **NATS**.
5.  **Payment Service** receives `order.created`:
    *   Processes payment (simulated).
    *   Persists payment record.
    *   Publishes `payment.confirmed` event.
6.  **Inventory Service** receives `payment.confirmed`:
    *   Deducts stock for the ordered items.
    *   Persists stock update.
7.  **Notification Service** receives `payment.confirmed`:
    *   Sends a confirmation notification to the user.

## Data Schema (High Level)

### User Service
*   `users`: `id`, `email`, `password_hash`, `full_name`, `created_at`

### Order Service
*   `orders`: `id`, `user_id`, `total_price`, `status`, `created_at`
*   `order_items`: `id`, `order_id`, `product_id`, `quantity`, `unit_price`

### Payment Service
*   `payments`: `id`, `order_id`, `amount`, `status`, `provider_transaction_id`, `created_at`

### Inventory Service
*   `products`: `id`, `name`, `description`, `price`, `stock_quantity`

## Observability & Reliability

*   **Tracing**: OpenTelemetry integrated into all services, exporting to Jaeger/Tempo.
*   **Metrics**: Prometheus metrics exported from each service.
*   **Retries**: NATS JetStream provides built-in retry mechanisms for event delivery. gRPC clients implement exponential backoff.
*   **Idempotency**: All event consumers are designed to be idempotent to handle duplicate event deliveries.
