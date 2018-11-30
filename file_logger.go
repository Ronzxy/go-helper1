package logger

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type FileLogger struct {
	*LogWriter
	File *os.File
}

func NewFileLogger(level int, configFile string, perm uint32) (*FileLogger, error) {
	var (
		this = &FileLogger{}
		err  error
	)

	this.File, err = os.OpenFile(configFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error: Open config file %v", err))
	}

	this.LogWriter = NewLogWriter(this.File, level)

	return this, nil
}

func (this *FileLogger) rolling() {
	for {
		this.Infof("ddd %v", time.Now())
		time.Sleep(10 * time.Second)
	}
}
