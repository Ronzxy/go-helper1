package main

import (
	"github.com/skygangsta/go-logger"
	"path"
	"time"
)

type Message struct {
	Dir    string
	Perms  int
	Parent string
}

func main() {
	message1 := &Message{
		Dir:    "/home/sky",
		Perms:  755,
		Parent: "/home",
	}
	message2 := &Message{
		Dir:    "/home/sky",
		Perms:  755,
		Parent: "/home",
	}

	err := logger.Init(path.Join("logger.xml"))
	if err == nil {
		go func() {
			for {
				logger.Trace("Trace message")
				logger.Debug("Debug message")
				logger.Infof("Info message %d", 1024)
				logger.Warn(message1, message2)
				logger.Error("Error message", message1)
				time.Sleep(30 * time.Second)
			}
		}()

		select {
		case <-time.After(24 * 7 * time.Hour):
			return
		}
	}
}
