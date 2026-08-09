package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/config"
)

type JWTVerifier interface {
	VerifyToken(ctx context.Context, tokenString string) (jwt.MapClaims, error)
}

func AuthMiddleware(cfg *config.Config, devSkipAuth bool, verifier JWTVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		if devSkipAuth {
			c.Set("user", jwt.MapClaims{
				"sub":                "dev-user",
				"preferred_username": "dev",
				"realm_access":       map[string]interface{}{"roles": []interface{}{"admin"}},
			})
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenStr := parts[1]
		claims, err := verifier.VerifyToken(c.Request.Context(), tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}
