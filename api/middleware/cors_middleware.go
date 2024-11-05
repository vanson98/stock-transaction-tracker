package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func UseCors(engine *gin.Engine) {
	engine.Use(cors.Default())
}
