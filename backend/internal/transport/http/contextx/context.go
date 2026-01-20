package contextx

import "github.com/gin-gonic/gin"

const RequestIDContextKey = "request_id"

func RequestIDFrom(c *gin.Context) string {
	if v, ok := c.Get(RequestIDContextKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
