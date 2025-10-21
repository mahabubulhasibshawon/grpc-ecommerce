// Dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o grpc-ecommerce cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/grpc-ecommerce .
CMD ["./grpc-ecommerce"]