FROM golang:1.18-alpine AS builder
WORKDIR /app

COPY ./ ./
RUN go mod download
RUN go build -o user_balance ./cmd/main/main.go
EXPOSE 8080
CMD [ "/app/user_balance" ]