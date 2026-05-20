package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/dekaiju/go-skeleton/types"
)

func WithGinContext(c *gin.Context) *logrus.Entry {
	if c == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	traceID := c.GetString(types.ContextTraceID)
	if traceID == "" {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return logrus.WithField("traceId", traceID)
}
