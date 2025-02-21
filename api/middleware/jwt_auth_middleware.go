package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UseTokenVerification(router *gin.RouterGroup) {
	router.Use(func(ctx *gin.Context) {
		authenCookie, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		// make a verify token request to identity server
		req, err := http.NewRequest("GET", "http://localhost:9091/token-verification", nil)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		req.Header.Set("Authorization", authenCookie)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			ctx.AbortWithStatus(resp.StatusCode)
			return
		}
		defer resp.Body.Close()
		ctx.Next()
	})
}
