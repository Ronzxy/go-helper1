package logger

import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	*LogWriter

	color        *Color // 颜色
	consoleColor bool   // 是否使用终端颜色
}

var (
	defaultLogger = NewConsoleLogger(ALL)
)

func DefaultConsoleLogger() *ConsoleLogger {
	return defaultLogger
}

func NewConsoleLogger(level int) *ConsoleLogger {
	this := &ConsoleLogger{
		LogWriter:    NewLogWriter(DefaultWriter, level),
		color:        NewColor(),
		consoleColor: true,
	}

	return this
}

func (this *ConsoleLogger) addColor(color, format string) string {
	if this.consoleColor {
		format = fmt.Sprintf("%s%s%s%s", this.color.Clear(), color, format, this.color.Clear())
	}

	return format
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *ConsoleLogger) Tracef(format string, v ...interface{}) {
	format = this.addColor(this.color.Blue(), format)
	this.write(TRACE, "TRACE", format, v...)
}

func (this *ConsoleLogger) Debugf(format string, v ...interface{}) {
	format = this.addColor(this.color.Green(), format)
	this.write(DEBUG, "DEBUG", format, v...)
}

func (this *ConsoleLogger) Infof(format string, v ...interface{}) {
	format = this.addColor(this.color.Cyan(), format)
	this.write(INFO, "INFO", format, v...)
}

func (this *ConsoleLogger) Warnf(format string, v ...interface{}) {
	format = this.addColor(this.color.Magenta(), format)
	this.write(WARN, "WARN", format, v...)
}

func (this *ConsoleLogger) Errorf(format string, v ...interface{}) {
	format = this.addColor(this.color.Yello(), format)
	this.write(ERROR, "ERROR", format, v...)
}

func (this *ConsoleLogger) Fatalf(format string, v ...interface{}) {
	format = this.addColor(this.color.Red(), format)
	this.write(FATAL, "FATAL", format, v...)
	os.Exit(1)
}

func (this *ConsoleLogger) Trace(message string) {
	this.Tracef("%s", message)
}

func (this *ConsoleLogger) Debug(message string) {
	this.Debugf("%s", message)
}

func (this *ConsoleLogger) Info(message string) {
	this.Infof("%s", message)
}

func (this *ConsoleLogger) Warn(message string) {
	this.Warnf("%s", message)
}

func (this *ConsoleLogger) Error(message string) {
	this.Errorf("%s", message)
}

func (this *ConsoleLogger) Fatal(message string) {
	this.Fatalf("%s", message)
}
