package main

import (
	"github.com/skygangsta/go-logger"
	"path"
	"time"
)

func main() {
	err := logger.InitLogger(path.Join("logger.xml"))
	if err == nil {
		go func() {
			for {
				logger.Trace("Trace message")
				logger.Debug("Debug message")
				logger.Info("Info message")
				logger.Warn("Warn message")
				logger.Error("Error message")
				time.Sleep(30 * time.Second)
			}
		}()

		select {
		case <-time.After(24 * 7 * time.Hour):
			return
		}
	}
}
