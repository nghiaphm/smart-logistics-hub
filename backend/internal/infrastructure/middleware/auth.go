package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
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
				"realm_access":       map[string]interface{}{"roles": []interface{}{"system_admin"}},
			})
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(apierrors.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Error(apierrors.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := verifier.VerifyToken(c.Request.Context(), tokenStr)
		if err != nil {
			c.Error(apierrors.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}
