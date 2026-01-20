package response

import (
	"mutabaahapi/internal/transport/http/contextx"
	"mutabaahapi/internal/transport/http/httperr"

	"github.com/gin-gonic/gin"
)

type successBody struct {
	RequestID string  `json:"request_id,omitempty"`
	Message   string  `json:"message"`
	Data      any     `json:"data,omitempty"`
	Paging    *Paging `json:"paging,omitempty"`
}

// Paging info for list responses.
type Paging struct {
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
	Total  *int `json:"total,omitempty"`
}

type Option func(*successBody)

func WithPaging(p Paging) Option { return func(b *successBody) { b.Paging = &p } }

func Success(c *gin.Context, status int, message string, data any, opts ...Option) {
	reqID := contextx.RequestIDFrom(c)
	body := successBody{RequestID: reqID, Message: message, Data: data}
	for _, opt := range opts {
		opt(&body)
	}
	c.JSON(status, body)
}

func Error(c *gin.Context, err error) {
	reqID := contextx.RequestIDFrom(c)
	c.JSON(httperr.StatusFor(err), httperr.Body(err, reqID))
}
