FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.work ./
COPY pkg/ ./pkg/
COPY proto/ ./proto/
COPY services/notification-service/ ./services/notification-service/
RUN go build -o /notification-service ./services/notification-service/cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /notification-service /notification-service
EXPOSE 8081
CMD ["/notification-service"]
