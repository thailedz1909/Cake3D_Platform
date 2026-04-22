package main

import (
	"log"

	"cake3d_platform/internal/config"
	delivery "cake3d_platform/internal/delivery/http"
	"cake3d_platform/internal/domain"
	mysqlrepo "cake3d_platform/internal/repository/mysql"
	redisrepo "cake3d_platform/internal/repository/redis"
	"cake3d_platform/internal/usecase"
	"cake3d_platform/pkg/jwt"
)

func main() {
	cfg := config.Load()

	db, err := config.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("mysql connect failed: %v", err)
	}

	redisClient, err := config.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("redis connect failed: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("db migrate failed: %v", err)
	}

	userRepo := mysqlrepo.NewUserRepository(db)
	tokenRepo := redisrepo.NewTokenRepository(redisClient)
	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.AccessTokenTTL)
	authUC := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtManager)

	router := delivery.NewRouter(cfg, authUC, jwtManager, tokenRepo)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
