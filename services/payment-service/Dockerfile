FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.work ./
COPY pkg/ ./pkg/
COPY proto/ ./proto/
COPY services/payment-service/ ./services/payment-service/
RUN go build -o /payment-service ./services/payment-service/cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /payment-service /payment-service
EXPOSE 8081
CMD ["/payment-service"]
