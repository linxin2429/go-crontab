package common

import (
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
)

var Logger *logrus.Logger

func InitLogger(filename string) {
	logfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("create file %s failed: %v", filename, err)
	}
	logger := logrus.New()
	logger.SetOutput(io.MultiWriter(logfile, os.Stdout))
	Logger = logger
}
