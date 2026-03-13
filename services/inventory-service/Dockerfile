FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.work ./
COPY pkg/ ./pkg/
COPY proto/ ./proto/
COPY services/inventory-service/ ./services/inventory-service/
RUN go build -o /inventory-service ./services/inventory-service/cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /inventory-service /inventory-service
EXPOSE 8081
CMD ["/inventory-service"]
