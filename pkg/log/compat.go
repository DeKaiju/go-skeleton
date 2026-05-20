package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Fatal(v ...interface{}) {
	logrus.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	logrus.Fatalf(format, v...)
}

func Print(v ...interface{}) {
	logrus.Print(v...)
}

func Printf(format string, v ...interface{}) {
	logrus.Printf(format, v...)
}

func Println(v ...interface{}) {
	logrus.Println(v...)
}

func SetOutput(w *os.File) {
	logrus.SetOutput(w)
}
