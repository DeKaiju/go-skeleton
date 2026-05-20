package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type customFormatter struct {
	logrus.TextFormatter
	baseFormatter *logrus.TextFormatter
}

func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.baseFormatter == nil {
		f.baseFormatter = &logrus.TextFormatter{
			TimestampFormat:  f.TimestampFormat,
			FullTimestamp:    f.FullTimestamp,
			ForceColors:      f.ForceColors,
			DisableColors:    f.DisableColors,
			DisableTimestamp: f.DisableTimestamp,
			DisableSorting:   f.DisableSorting,
			SortingFunc:      f.SortingFunc,
			PadLevelText:     f.PadLevelText,
			QuoteEmptyFields: true,
		}
	}

	formatted, err := f.baseFormatter.Format(entry)
	if err != nil {
		return nil, err
	}
	if len(entry.Data) == 0 {
		return formatted, nil
	}

	formattedStr := string(formatted)
	msgStart := strings.Index(formattedStr, "msg=")
	if msgStart == -1 {
		return formatted, nil
	}

	beforeMsg := formattedStr[:msgStart]
	msgValueStart := msgStart + 4
	msgEnd := msgValueStart

	if msgValueStart < len(formattedStr) && formattedStr[msgValueStart] == '"' {
		msgEnd++
		for msgEnd < len(formattedStr) {
			if formattedStr[msgEnd] == '"' && formattedStr[msgEnd-1] != '\\' {
				msgEnd++
				break
			}
			msgEnd++
		}
	} else {
		for msgEnd < len(formattedStr) && formattedStr[msgEnd] != ' ' && formattedStr[msgEnd] != '\n' {
			msgEnd++
		}
	}

	msgPart := formattedStr[msgStart:msgEnd]
	afterMsg := strings.TrimSpace(formattedStr[msgEnd:])
	if afterMsg == "" || afterMsg == "\n" {
		return formatted, nil
	}

	afterMsg = strings.TrimSuffix(afterMsg, "\n")
	result := beforeMsg + afterMsg + " " + msgPart + "\n"
	return []byte(result), nil
}

func Setup() {
	forceColors := false
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		forceColors = (fileInfo.Mode() & os.ModeCharDevice) != 0
	}

	logrus.SetFormatter(&customFormatter{
		TextFormatter: logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000000",
			FullTimestamp:   true,
			ForceColors:     forceColors,
		},
	})

	logrus.SetLevel(getLogLevel())
}

func getLogLevel() logrus.Level {
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}
