package contextx

import "github.com/gin-gonic/gin"

const (
	UserIDContextKey    = "user_id"
	UserEmailContextKey = "user_email"
)

func UserIDFrom(c *gin.Context) (string, bool) {
    v, ok := c.Get(UserIDContextKey)
    if !ok { return "", false }
    id, ok := v.(string)
    return id, ok
}

func UserEmailFrom(c *gin.Context) (string, bool) {
	v, ok := c.Get(UserEmailContextKey)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
