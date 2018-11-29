package logger

import (
	"os"
)

type FileLogger struct {
	*LogWriter
	File *os.File
}

func NewFileLogger(level int, configFile string, perm uint32) *FileLogger {
	var (
		this = &FileLogger{}
		err  error
	)

	this.File, err = os.OpenFile(configFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		DefaultConsoleLogger().Errorf("error: Open config file %v", err)
		return nil
	}

	this.LogWriter = NewLogWriter(this.File, level)

	return this
}
