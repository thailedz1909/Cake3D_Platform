# Cake3D_Platform

Clean Architecture sample API (Go) with:
- JWT authentication
- MySQL
- Redis (token blacklist for logout)
- Docker / Docker Compose

## Run with Docker

```bash
cp .env.example .env
docker compose up --build
```

API runs at: `http://localhost:8080`

## Main Endpoints

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout` (Bearer token required)
- `GET /api/v1/users/me` (Bearer token required)

## Local Run (without Docker)

```bash
cp .env.example .env
go mod tidy
go run ./cmd/api
```
