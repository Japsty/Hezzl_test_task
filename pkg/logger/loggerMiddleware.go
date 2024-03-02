package logger

import (
	"github.com/gin-gonic/gin"
	"log"
)

// LogMiddleware - мидлварь для отслеживания логов по запросам
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request: %s %s %s", c.Request.Method, c.Request.URL, c.Params)
		c.Next()
	}
}
