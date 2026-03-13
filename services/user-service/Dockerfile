FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.work ./
COPY pkg/ ./pkg/
COPY proto/ ./proto/
COPY services/user-service/ ./services/user-service/
RUN go build -o /user-service ./services/user-service/cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /user-service /user-service
EXPOSE 50051 8081
CMD ["/user-service"]
