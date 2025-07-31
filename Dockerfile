# 1. Aşama: Build
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go

# 2. Aşama: Çalışma
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server .
COPY .env .

CMD ["./server"]
