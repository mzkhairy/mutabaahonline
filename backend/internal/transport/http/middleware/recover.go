package middleware

import (
	"mutabaahapi/internal/platform/apperror"
	resp "mutabaahapi/internal/transport/http/response"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Recover is a custom recovery middleware that logs the panic and returns a JSON INTERNAL_ERROR.
func Recover(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				l.Error().Any("panic", r).Bytes("stack", debug.Stack()).Msg("panic recovered")
				resp.Error(c, apperror.ErrInternalError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
