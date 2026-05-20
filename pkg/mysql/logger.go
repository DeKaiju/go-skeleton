package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/dekaiju/go-skeleton/types"
)

type traceLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
	Colorful      bool
}

func NewTraceLogger(logLevel logger.LogLevel, slowThreshold time.Duration) *traceLogger {
	l := &traceLogger{}
	l.LogLevel = logLevel
	l.SlowThreshold = slowThreshold
	l.Colorful = true
	return l
}

func (l *traceLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l traceLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		output := append([]interface{}{utils.FileWithLineNum()}, data...)
		output = append(output, ctx.Value(types.ContextTraceID))
		logrus.Printf(msg, output...)
	}
}

func (l traceLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		output := append([]interface{}{utils.FileWithLineNum()}, data...)
		output = append(output, ctx.Value(types.ContextTraceID))
		logrus.Printf(msg, output...)
	}
}

func (l traceLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		output := append([]interface{}{utils.FileWithLineNum()}, data...)
		output = append(output, ctx.Value(types.ContextTraceID))
		logrus.Printf(msg, output...)
	}
}

func (l traceLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel == 0 {
		return
	}

	elapsed := time.Since(begin)
	traceID, _ := ctx.Value(types.ContextTraceID).(string)

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		logSQL(logrus.ErrorLevel, traceID, elapsed, err, fc)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		logSQL(logrus.WarnLevel, traceID, elapsed, fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold), fc)
	case l.LogLevel >= logger.Info:
		logSQL(logrus.InfoLevel, traceID, elapsed, nil, fc)
	}
}

func logSQL(level logrus.Level, traceID string, elapsed time.Duration, sqlErr interface{}, fc func() (string, int64)) {
	sql, rows := fc()
	rowsValue := fmt.Sprintf("%d", rows)
	if rows == -1 {
		rowsValue = "-"
	}

	entry := logrus.WithFields(logrus.Fields{
		"file": utils.FileWithLineNum(),
		"rows": rowsValue,
		"sql":  sql,
	})
	if traceID != "" {
		entry = entry.WithField("traceId", traceID)
	}
	if sqlErr != nil {
		entry = entry.WithField("error", sqlErr)
	}

	elapsedMS := fmt.Sprintf("%.2fms", float64(elapsed.Nanoseconds())/1e6)
	switch level {
	case logrus.ErrorLevel:
		entry.Errorf("[SQL Error] %s", elapsedMS)
	case logrus.WarnLevel:
		entry.Warnf("[SQL Slow] %s", elapsedMS)
	default:
		entry.Infof("[SQL] %s", elapsedMS)
	}
}
