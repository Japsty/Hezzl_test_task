package logger

import (
	"github.com/gin-gonic/gin"
	"log"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL, c.Request.Body, c.Request.Response)
		c.Next()
	}
}
