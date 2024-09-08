package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func AuthAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != os.Getenv("X_API_KEY") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
