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
	DefaultWriter io.Writer = os.Stdout
)

type LogWriter struct {
	allowLevel int // 日志级别
	denyLevel  int
	name       string      // 日志名称
	timeFormat string      // 时间格式
	logger     *log.Logger // 日志对象
}

func NewLogWriter(w io.Writer, level int) *LogWriter {
	this := &LogWriter{
		allowLevel: level,
		denyLevel:  OFF,
		name:       util.NewPath().WorkName(),
		timeFormat: LogTimeFormat,
		logger:     log.New(w, "", log.LUTC),
	}

	return this
}

func (this *LogWriter) SetDenyLevel(level int) {
	if level > this.denyLevel {
		this.allowLevel = OFF
	} else {
		this.denyLevel = level
	}
}

func (this *LogWriter) SetName(name string) {
	this.name = name
}

func (this *LogWriter) SetWriter(w io.Writer) {
	this.logger.SetOutput(w)
}

func (this *LogWriter) SetTimeFormat(timeFormat string) {
	this.timeFormat = timeFormat
}

func (this *LogWriter) write(level int, levelName, format string, v ...interface{}) {
	if this.allowLevel <= level {
		if this.denyLevel > level {
			format = fmt.Sprintf("%s - %s - %-5s - %s\n",
				this.name, time.Now().Format(this.timeFormat), levelName, format)

			this.logger.Printf(format, v...)
		}
	}
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *LogWriter) Tracef(format string, v ...interface{}) {
	this.write(TRACE, "TRACE", format, v...)
}

func (this *LogWriter) Debugf(format string, v ...interface{}) {
	this.write(DEBUG, "DEBUG", format, v...)
}

func (this *LogWriter) Infof(format string, v ...interface{}) {
	this.write(INFO, "INFO", format, v...)
}

func (this *LogWriter) Warnf(format string, v ...interface{}) {
	this.write(WARN, "WARN", format, v...)
}

func (this *LogWriter) Errorf(format string, v ...interface{}) {
	this.write(ERROR, "ERROR", format, v...)
}

func (this *LogWriter) Fatalf(format string, v ...interface{}) {
	this.write(FATAL, "FATAL", format, v...)
	os.Exit(1)
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *LogWriter) Trace(message string) {
	this.Tracef("%s", message)
}

func (this *LogWriter) Debug(message string) {
	this.Debugf("%s", message)
}

func (this *LogWriter) Info(message string) {
	this.Infof("%s", message)
}

func (this *LogWriter) Warn(message string) {
	this.Warnf("%s", message)
}

func (this *LogWriter) Error(message string) {
	this.Errorf("%s", message)
}

func (this *LogWriter) Fatal(message string) {
	this.Fatalf("%s", message)
}
