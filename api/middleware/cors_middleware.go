package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func UseCors(router *gin.RouterGroup) {
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
	})
	router.Use(corsMiddleware)
}
