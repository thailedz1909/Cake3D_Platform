package http

import (
	"cake3d_platform/internal/config"
	"cake3d_platform/internal/delivery/http/handler"
	"cake3d_platform/internal/delivery/http/middleware"
	"cake3d_platform/internal/domain"
	"cake3d_platform/internal/usecase"
	"cake3d_platform/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, authUC usecase.AuthUseCase, jwtManager *jwt.Manager, tokenRepo domain.TokenRepository) *gin.Engine {
	if cfg.AppEnv == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	authHandler := handler.NewAuthHandler(authUC)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager, tokenRepo)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authMiddleware.RequireAuth(), authHandler.Logout)
		}

		users := v1.Group("/users")
		users.Use(authMiddleware.RequireAuth())
		{
			users.GET("/me", authHandler.Me)
		}
	}

	return r
}
