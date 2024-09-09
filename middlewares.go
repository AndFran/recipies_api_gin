package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
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

func AuthJWTAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			log.Println("err parsing token:", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// check valid token
		if token == nil || !token.Valid {
			log.Println("token is nil or not valid:", token.Valid)

			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//claims := token.Claims.(jwt.MapClaims)  // if we need the claims

		c.Next()
	}
}
