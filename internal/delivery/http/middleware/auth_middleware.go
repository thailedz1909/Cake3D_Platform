package middleware

import (
	"errors"
	"net/http"
	"strings"

	"cake3d_platform/internal/domain"
	"cake3d_platform/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager *jwt.Manager
	tokenRepo  domain.TokenRepository
}

func NewAuthMiddleware(jwtManager *jwt.Manager, tokenRepo domain.TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		tokenRepo:  tokenRepo,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if strings.TrimSpace(token) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid authorization header"})
			return
		}

		claims, err := m.jwtManager.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		blacklisted, err := m.tokenRepo.IsBlacklisted(c.Request.Context(), claims.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to validate token"})
			return
		}
		if blacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token revoked"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func ClaimsFromContext(c *gin.Context) (*jwt.Claims, error) {
	val, ok := c.Get("claims")
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	claims, ok := val.(*jwt.Claims)
	if !ok {
		return nil, errors.New("invalid claims in context")
	}
	return claims, nil
}
