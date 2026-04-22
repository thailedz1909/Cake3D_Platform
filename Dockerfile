FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/api ./api
COPY .env.example ./.env

EXPOSE 8080

CMD ["./api"]
