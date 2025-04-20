FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/task_service cmd/main.go

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app
COPY --from=builder /app/bin/task_service .
COPY .env .

EXPOSE 8080

CMD ["./task_service"]