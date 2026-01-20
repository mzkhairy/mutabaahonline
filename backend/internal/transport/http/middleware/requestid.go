package middleware

import (
	"mutabaahapi/internal/transport/http/contextx"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(contextx.RequestIDContextKey, id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}
