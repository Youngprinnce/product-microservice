package logger

import (
"os"

log "github.com/sirupsen/logrus"
)

func Initialize() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func Info(msg string) {
	log.Info(msg)
}

func Error(msg string) {
	log.Error(msg)
}

func Fatal(msg string) {
	log.Fatal(msg)
}

func Debug(msg string) {
	log.Debug(msg)
}

func Warn(msg string) {
	log.Warn(msg)
}
