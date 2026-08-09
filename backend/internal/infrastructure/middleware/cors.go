package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(frontendURL string) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{frontendURL}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	return cors.New(corsConfig)
}
