package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		claims, ok := userClaims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user claims type"})
			return
		}

		userRoles := make(map[string]bool)

		if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
			if roles, ok := realmAccess["roles"].([]interface{}); ok {
				for _, r := range roles {
					if roleStr, ok := r.(string); ok {
						userRoles[roleStr] = true
					}
				}
			}
		}

		if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
			for _, clientData := range resourceAccess {
				if clientMap, ok := clientData.(map[string]interface{}); ok {
					if roles, ok := clientMap["roles"].([]interface{}); ok {
						for _, r := range roles {
							if roleStr, ok := r.(string); ok {
								userRoles[roleStr] = true
							}
						}
					}
				}
			}
		}

		hasRole := false
		for _, role := range allowedRoles {
			if userRoles[role] {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied. Required roles missing"})
			return
		}

		c.Next()
	}
}
