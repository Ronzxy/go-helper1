/* Copyright 2018 sky<skygangsta@hotmail.com>. All rights reserved.
 *
 * Licensed under the Apache License, version 2.0 (the "License").
 * You may not use this work except in compliance with the License, which is
 * available at www.apache.org/licenses/LICENSE-2.0
 *
 * This software is distributed on an "AS IS" basis, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied, as more fully set forth in the License.
 *
 * See the NOTICE file distributed with this work for information regarding copyright ownership.
 */

package logger

import (
	"fmt"
	"github.com/skygangsta/go-utils"
	"io"
	"log"
	"os"
	"time"
)

const (
	ALL   = 0
	TRACE = 1
	DEBUG = 2
	INFO  = 3
	WARN  = 4
	ERROR = 5
	FATAL = 6
	OFF   = 7

	LogTimeFormat = "2006/01/02 15:04:05.000000"
)

var (
	DefaultWriter     io.Writer = os.Stdout
	DefaultLogger               = NewLogger(ALL)
	DefaultLoggerInfo           = NewLogger(INFO)
)

type Logger struct {
	level        int         // 日志级别
	writer       io.Writer   // 写入io
	color        *Color      // 颜色
	consoleColor bool        // 是否使用终端颜色
	name         string      // 日志名称
	timeFormat   string      // 时间格式
	logger       *log.Logger // 日志对象
}

func NewLogger(level int) *Logger {
	return NewLoggerWithWriter(DefaultWriter, level)
}

func Default() *Logger {
	return DefaultLogger
}

func DefaultInfo() *Logger {
	return DefaultLoggerInfo
}

func NewLoggerWithWriter(w io.Writer, level int) *Logger {
	this := &Logger{
		level:        level,
		writer:       w,
		color:        NewColor(),
		consoleColor: false,
		name:         util.NewPath().WorkName(),
		timeFormat:   LogTimeFormat,
		logger:       log.New(w, "", log.LUTC),
	}

	return this
}

func (this *Logger) SetName(name string) {
	this.name = name
}

func (this *Logger) SetWriter(w io.Writer) {
	this.writer = w
	this.logger.SetOutput(w)
}

func (this *Logger) UseConsoleColor(b bool) {
	this.consoleColor = b
}

func (this *Logger) SetTimeFormat(timeFormat string) {
	this.timeFormat = timeFormat
}

func (this *Logger) write(color, level, format string, v ...interface{}) {
	if this.writer == os.Stdout {
		if !this.consoleColor {
			color = ""
		}

		format = fmt.Sprintf("%s%s - %s - %-5s - %s%s%s\n",
			this.color.Clear(),
			this.name, time.Now().Format(this.timeFormat),
			level, color, format,
			this.color.Clear())
	} else {
		format = fmt.Sprintf("%s - %s - %-5s - %s\n",
			this.name, time.Now().Format(this.timeFormat), level, format)
	}

	this.logger.Printf(format, v...)
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *Logger) Tracef(format string, v ...interface{}) {
	if this.level <= TRACE {
		this.write(this.color.Blue(), "TRACE", format, v...)
	}
}

func (this *Logger) Debugf(format string, v ...interface{}) {
	if this.level <= DEBUG {
		this.write(this.color.Green(), "DEBUG", format, v...)
	}
}

func (this *Logger) Infof(format string, v ...interface{}) {
	if this.level <= INFO {
		this.write(this.color.Cyan(), "INFO", format, v...)
	}
}

func (this *Logger) Warnf(format string, v ...interface{}) {
	if this.level <= WARN {
		this.write(this.color.Magenta(), "WARN", format, v...)
	}
}

func (this *Logger) Errorf(format string, v ...interface{}) {
	if this.level <= ERROR {
		this.write(this.color.Yello(), "ERROR", format, v...)
	}
}

func (this *Logger) Fatalf(format string, v ...interface{}) {
	if this.level <= FATAL {
		this.write(this.color.Red(), "FATAL", format, v...)
		os.Exit(1)
	}
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *Logger) Trace(message string) {
	this.Tracef("%s", message)
}

func (this *Logger) Debug(message string) {
	this.Debugf("%s", message)
}

func (this *Logger) Info(message string) {
	this.Infof("%s", message)
}

func (this *Logger) Warn(message string) {
	this.Warnf("%s", message)
}

func (this *Logger) Error(message string) {
	this.Errorf("%s", message)
}

func (this *Logger) Fatal(message string) {
	this.Fatalf("%s", message)
}
