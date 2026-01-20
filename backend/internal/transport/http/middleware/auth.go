package middleware

import (
	"mutabaahapi/internal/platform/apperror"
	appjwt "mutabaahapi/internal/platform/jwt"
	"mutabaahapi/internal/transport/http/contextx"
	resp "mutabaahapi/internal/transport/http/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func AuthRequired(l zerolog.Logger, t *appjwt.Tokenizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := strings.TrimSpace(c.GetHeader("Authorization"))
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			resp.Error(c, apperror.New(apperror.ErrCodeUnauthorized, "missing bearer token"))
			c.Abort()
			return
		}

		token := strings.TrimSpace(h[len("Bearer "):])
		if token == "" {
			resp.Error(c, apperror.New(apperror.ErrCodeUnauthorized, "invalid token"))
			c.Abort()
			return
		}

		claims, err := t.Verify(token)
		if err != nil {
			l.Warn().Err(err).Msg("token verify failed")
			resp.Error(c, apperror.New(apperror.ErrCodeUnauthorized, "invalid token"))
			c.Abort()
			return
		}

		// Set User ID ke Context (ini sudah benar)
		c.Set(contextx.UserIDContextKey, claims.UserID)

		// REVISI: Hapus set Email, ganti dengan Role
		// c.Set(contextx.UserEmailContextKey, claims.Email) <--- INI PENYEBAB ERROR

		// Set Role agar bisa diakses di handler selanjutnya
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
