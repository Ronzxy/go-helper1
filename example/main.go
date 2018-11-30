package main

import (
	"github.com/skygangsta/go-logger"
	"path"
)

func main() {
	err := logger.InitLogger(path.Join("logger.xml"))
	if err == nil {
		logger.Trace("Trace message")
		logger.Debug("Debug message")
		logger.Info("Info message")
		logger.Warn("Warn message")
		logger.Error("Error message")
		logger.Fatal("Fatal message")
	}
}
