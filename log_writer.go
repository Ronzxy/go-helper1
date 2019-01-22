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
	"github.com/go-xorm/core"
	"github.com/skygangsta/go-helper"
	"io"
	"log"
	"os"
	"runtime"
)

type LogWriter struct {
	allowLevel      int
	denyLevel       int
	workName        string // 工作名
	loggerName      string // 日志名
	logger          *log.Logger
	formatter       Formatter
	skipCallerDepth int
	showSQL         bool
}

func NewLogWriter(w io.Writer, level int) *LogWriter {
	this := &LogWriter{
		allowLevel:      level,
		denyLevel:       OFF,
		workName:        helper.NewPathHelper().WorkName(),
		logger:          log.New(w, "", log.LUTC),
		skipCallerDepth: defaultSkipCallerDepth,
	}

	this.SetFormatter(NewTextFormatter())

	return this
}

func (this *LogWriter) SetDenyLevel(level int) {
	if level > this.denyLevel {
		this.allowLevel = OFF
	} else {
		this.denyLevel = level
	}
}

func (this *LogWriter) SetSkipCallerDepth(skipCallerDepth int) {
	this.skipCallerDepth = skipCallerDepth
}

func (this *LogWriter) SetName(name string) {
	this.workName = name
}

func (this *LogWriter) SetWriter(w io.Writer) {
	this.logger.SetOutput(w)
}

func (this *LogWriter) SetFormatter(formatter Formatter) {
	this.formatter = formatter
}

func (this *LogWriter) filter(frame *runtime.Frame) bool {

	packageName := GetPackageName(frame.Function)
	if config.Filters != nil {
		for _, filter := range config.Filters {
			if packageName == filter.Name {
				if len(filter.Loggers) > 0 {
					for _, name := range filter.Loggers {
						if name == this.loggerName {
							return true
						}
					}
				}
			}
		}
	}

	if len(config.Default.Loggers) > 0 {
		for _, name := range config.Default.Loggers {
			if name == this.loggerName {
				return true
			}
		}
	}

	return false
}

func (this *LogWriter) Println(level int, args ...interface{}) {
	frame := GetCaller(this.skipCallerDepth)

	if !this.filter(frame) {
		return
	}

	if this.allowLevel <= level {
		if this.denyLevel > level {
			var (
				data = map[string]interface{}{}
			)

			data["Name"] = this.workName
			data["Level"] = ConvertLevel2String(level)
			data["Line"] = frame.Line
			data["PackageName"] = GetPackageName(frame.Function)
			data["File"] = GetFileName(frame)

			this.logger.Println(this.formatter.Message(data, args...))
		}
	}
}

/*
Implement Writer
*/
func (this *FileLogger) Tracef(format string, args ...interface{}) {
	this.Println(TRACE, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Debugf(format string, args ...interface{}) {
	this.Println(DEBUG, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Infof(format string, args ...interface{}) {
	this.Println(INFO, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Warnf(format string, args ...interface{}) {
	this.Println(WARN, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Errorf(format string, args ...interface{}) {
	this.Println(ERROR, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Fatalf(exit bool, format string, args ...interface{}) {
	this.Println(FATAL, fmt.Sprintf(format, args...))

	if exit {
		os.Exit(1)
	}
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *FileLogger) Trace(args ...interface{}) {
	this.Println(TRACE, args...)
}

func (this *FileLogger) Debug(args ...interface{}) {
	this.Println(DEBUG, args...)
}

func (this *FileLogger) Info(args ...interface{}) {
	this.Println(INFO, args...)
}

func (this *FileLogger) Warn(args ...interface{}) {
	this.Println(WARN, args...)
}

func (this *FileLogger) Error(args ...interface{}) {
	this.Println(ERROR, args...)
}

func (this *FileLogger) Fatal(exit bool, args ...interface{}) {
	this.Println(FATAL, args...)

	if exit {
		os.Exit(1)
	}
}

/*
Implement xorm logger
*/

// Set to xorm log all
func (this *LogWriter) Level() core.LogLevel {
	return core.LOG_DEBUG
}

func (this *LogWriter) SetLevel(l core.LogLevel) {

}

func (this *LogWriter) ShowSQL(show ...bool) {
	if len(show) == 0 {
		this.showSQL = true
		return
	}
	this.showSQL = show[0]
}

func (this *LogWriter) IsShowSQL() bool {
	return this.showSQL
}
