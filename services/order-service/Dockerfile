FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.work ./
COPY pkg/ ./pkg/
COPY proto/ ./proto/
COPY services/order-service/ ./services/order-service/
RUN go build -o /order-service ./services/order-service/cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /order-service /order-service
EXPOSE 50052 8081
CMD ["/order-service"]
