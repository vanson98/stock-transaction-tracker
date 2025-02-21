package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func UseCors(router *gin.RouterGroup) {

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:              []string{"http://localhost:8081"},
		AllowCredentials:          true,
		AllowMethods:              []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:              []string{"Content-Type", "Authorization"},
		MaxAge:                    12 * time.Hour,
		OptionsResponseStatusCode: http.StatusOK,
	})

	router.Use(corsMiddleware)

	// Handle OPTIONS requests
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}
