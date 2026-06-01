FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o sams ./cmd/server/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/sams .
EXPOSE 8080
CMD ["./sams"]
