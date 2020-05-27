/* Copyright 2018 Ron Zhang <ronzxy@mx.aketi.cn>. All rights reserved.
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
	"github.com/ronzxy/go-helper"
	"io"
	"log"
	"os"
	"runtime"
	xormlog "xorm.io/xorm/log"
)

type LoggerWriter struct {
	allowLevel      LogLevel
	denyLevel       LogLevel
	prefix          string // 工作名
	name            string // 日志名
	formatter       Formatter
	skipCallerDepth int
	closeFilter     bool
	showSQL         bool

	*log.Logger
}

func NewLoggerWriter(w io.Writer, level LogLevel) *LoggerWriter {
	this := &LoggerWriter{
		allowLevel:      level,
		denyLevel:       OFF,
		prefix:          helper.Path.WorkName(),
		skipCallerDepth: defaultSkipCallerDepth,
		Logger:          log.New(w, "", log.LUTC),
	}

	this.SetFormatter(NewTextFormatter())

	return this
}

func (this *LoggerWriter) NewLogger(w io.Writer) {
	this.Logger = log.New(w, "", log.LUTC)
}

func (this *LoggerWriter) SetDenyLevel(level LogLevel) {
	if level > this.denyLevel {
		this.allowLevel = OFF
	} else {
		this.denyLevel = level
	}
}

func (this *LoggerWriter) SetSkipCallerDepth(skipCallerDepth int) {
	this.skipCallerDepth = skipCallerDepth
}

func (this *LoggerWriter) SetWriter(w io.Writer) {
	this.SetOutput(w)
}

func (this *LoggerWriter) SetFormatter(formatter Formatter) {
	this.formatter = formatter
}

func (this *LoggerWriter) filter(frame *runtime.Frame) bool {
	if this.closeFilter {
		return true
	}

	if config.PackageFilters == nil {
		return false
	}

	for _, filter := range config.PackageFilters {
		packageName := GetPackageName(frame.Function)
		if packageName == filter.Name {
			if len(filter.Loggers) == 0 {
				// No Logger define
				return false
			}

			for _, name := range filter.Loggers {
				if name == this.name {
					return true
				}
			}
		}
	}

	if len(config.DefaultFilter.Loggers) == 0 {
		// No Logger define
		return false
	}

	// TODO: if define by package name whether the default filer can be used?
	for _, name := range config.DefaultFilter.Loggers {
		if name == this.name {
			return true
		}
	}

	return false
}

func (this *LoggerWriter) Write(level LogLevel, args ...interface{}) error {
	// Reject logs that are less than the allowed level
	if level < this.allowLevel {
		return nil
	}

	// Reject logs greater than or equal to the rejection level
	if level >= this.denyLevel {
		return nil
	}

	frame := GetCaller(this.skipCallerDepth)

	if !this.filter(frame) {
		return nil
	}

	var (
		data = map[string]interface{}{}
	)

	data["Prefix"] = this.prefix
	data["Level"] = ConvertLevel2String(level)
	data["Line"] = frame.Line
	data["PackageName"] = GetPackageName(frame.Function)
	data["File"] = GetFileName(frame)

	return this.Logger.Output(0, this.formatter.Message(data, args...))
}

/*
Rewrite log.Logger function
*/
func (this *LoggerWriter) Output(calldepth int, s string) error {

	return this.Write(ALL, s)
}

// Printf calls this.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (this *LoggerWriter) Printf(format string, v ...interface{}) {
	this.Write(ALL, fmt.Sprintf(format, v...))
}

// Print calls this.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (this *LoggerWriter) Print(v ...interface{}) { this.Write(ALL, fmt.Sprint(v...)) }

// Println calls this.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (this *LoggerWriter) Println(v ...interface{}) { this.Write(ALL, fmt.Sprintln(v...)) }

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (this *LoggerWriter) Fatal(v ...interface{}) {
	this.Write(ALL, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (this *LoggerWriter) Fatalf(format string, v ...interface{}) {
	this.Write(ALL, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (this *LoggerWriter) Fatalln(v ...interface{}) {
	this.Write(ALL, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (this *LoggerWriter) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	this.Write(ALL, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (this *LoggerWriter) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	this.Write(ALL, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (this *LoggerWriter) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	this.Write(ALL, s)
	panic(s)
}

// Prefix returns the output prefix for the logger.
func (this *LoggerWriter) Prefix() string {

	return this.prefix
}

// SetPrefix sets the output prefix for the logger.
func (this *LoggerWriter) SetPrefix(prefix string) {
	this.prefix = prefix
}

/*
Implement Writer
*/
func (this *LoggerWriter) Tracef(format string, args ...interface{}) {
	this.Write(TRACE, fmt.Sprintf(format, args...))
}

func (this *LoggerWriter) Debugf(format string, args ...interface{}) {
	this.Write(DEBUG, fmt.Sprintf(format, args...))
}

func (this *LoggerWriter) Infof(format string, args ...interface{}) {
	this.Write(INFO, fmt.Sprintf(format, args...))
}

func (this *LoggerWriter) Warnf(format string, args ...interface{}) {
	this.Write(WARN, fmt.Sprintf(format, args...))
}

func (this *LoggerWriter) Errorf(format string, args ...interface{}) {
	this.Write(ERROR, fmt.Sprintf(format, args...))
}

func (this *LoggerWriter) FatalfWithExit(exit bool, format string, args ...interface{}) {
	this.Write(FATAL, fmt.Sprintf(format, args...))

	if exit {
		os.Exit(1)
	}
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *LoggerWriter) Trace(args ...interface{}) {
	this.Write(TRACE, args...)
}

func (this *LoggerWriter) Debug(args ...interface{}) {
	this.Write(DEBUG, args...)
}

func (this *LoggerWriter) Info(args ...interface{}) {
	this.Write(INFO, args...)
}

func (this *LoggerWriter) Warn(args ...interface{}) {
	this.Write(WARN, args...)
}

func (this *LoggerWriter) Error(args ...interface{}) {
	this.Write(ERROR, args...)
}

func (this *LoggerWriter) FatalWithExit(exit bool, args ...interface{}) {
	this.Write(FATAL, args...)

	if exit {
		os.Exit(1)
	}
}

/*
Implement xorm logger
*/

func (this *LoggerWriter) Level() xormlog.LogLevel {
	// Set to xorm log all
	return xormlog.LOG_DEBUG
}

func (this *LoggerWriter) SetLevel(l xormlog.LogLevel) {

}

func (this *LoggerWriter) ShowSQL(show ...bool) {
	if len(show) == 0 {
		this.showSQL = true
		return
	}
	this.showSQL = show[0]
}

func (this *LoggerWriter) IsShowSQL() bool {
	return this.showSQL
}

func (this *LoggerWriter) BeforeSQL(context xormlog.LogContext) {
	fmt.Println(context)
}

func (this *LoggerWriter) AfterSQL(context xormlog.LogContext) {
	fmt.Println(context)
}
