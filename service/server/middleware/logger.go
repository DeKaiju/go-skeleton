package middleware

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/dekaiju/go-skeleton/types"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &responseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		c.Next()

		statusCode := c.Writer.Status()
		if statusCode < 400 {
			return
		}

		fields := logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  statusCode,
			"traceId": c.GetString(types.ContextTraceID),
		}

		var responseBody map[string]interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &responseBody); err == nil {
			if message, ok := responseBody["message"].(string); ok && message != "" {
				fields["message"] = message
			}
		}

		logrus.WithFields(fields).Errorf("request failed with status %d", statusCode)
	}
}
