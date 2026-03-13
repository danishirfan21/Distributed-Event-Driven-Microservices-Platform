FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.work and all go.mod files
COPY go.work ./
COPY pkg/ ./pkg/
COPY proto/ ./proto/
COPY services/api-gateway/ ./services/api-gateway/

# Build
RUN go build -o /api-gateway ./services/api-gateway/cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /api-gateway /api-gateway
EXPOSE 8080 8081
CMD ["/api-gateway"]
