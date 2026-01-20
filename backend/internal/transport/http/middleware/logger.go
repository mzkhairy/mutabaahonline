package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"mutabaahapi/internal/platform/logx"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const maxLoggedBody = 8192 // 8KB cap for req/resp bodies

type bodyCaptureWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *bodyCaptureWriter) Write(b []byte) (int, error) {
	if w.buf != nil && w.buf.Len() < maxLoggedBody {
		// copy up to cap
		remaining := maxLoggedBody - w.buf.Len()
		if len(b) > remaining {
			w.buf.Write(b[:remaining])
		} else {
			w.buf.Write(b)
		}
	}
	return w.ResponseWriter.Write(b)
}

func Logger(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		reqID := c.GetString("request_id")
		// request body capture (only JSON and capped)
		var reqBody []byte
		if c.Request.Body != nil && c.ContentType() == "application/json" {
			b, _ := io.ReadAll(io.LimitReader(c.Request.Body, maxLoggedBody))
			reqBody = b
			c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
		}
		// response capture
		bw := &bodyCaptureWriter{ResponseWriter: c.Writer, buf: &bytes.Buffer{}}
		c.Writer = bw

		// attach request-scoped logger to context
		reqLogger := l.With().Str("req_id", reqID).Str("method", method).Str("path", path).Logger()
		c.Set("logger", reqLogger)
		c.Request = c.Request.WithContext(logx.NewContext(c.Request.Context(), reqLogger))

		// start event (debug-level; will be shown only when level is debug)
		reqLogger.Info().Msg("http_request")

		c.Next()

		dur := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		evt := reqLogger.Info()
		if status >= 500 {
			evt = reqLogger.Error()
		} else if status >= 400 {
			evt = reqLogger.Warn()
		}
		evt.Int("status", status).
			Int64("duration_ms", dur.Milliseconds()).
			Int("bytes_out", size).
			Msg("http_response")

		// For error responses, add request/response bodies to logs (capped)
		if status >= 400 {
			e := reqLogger.WithLevel(zerolog.WarnLevel)
			if len(reqBody) > 0 {
				if json.Valid(reqBody) {
					e = e.RawJSON("req_body", reqBody)
				} else {
					e = e.Str("req_body_text", string(reqBody))
				}
			}
			respBody := bw.buf.Bytes()
			if len(respBody) > 0 {
				if json.Valid(respBody) {
					e = e.RawJSON("resp_body", respBody)
				} else {
					e = e.Str("resp_body_text", string(respBody))
				}
			}
			e.Int("status", status).Msg("http_error")
		}
	}
}
