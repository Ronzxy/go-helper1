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
	"os"
)

type ConsoleLogger struct {
	*LogWriter

	color        *Color
	consoleColor bool // use terminal color
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

	this.SetFormatter(NewTextFormatter())

	return this
}

func (this *ConsoleLogger) addColor(color string, args ...interface{}) []interface{} {
	if this.consoleColor {
		format := append([]interface{}{}, color)
		args = append(format, args...)
		args = append(args, this.color.Clear())
	}

	return args
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *ConsoleLogger) Tracef(format string, args ...interface{}) {
	args = this.addColor(this.color.Blue(), fmt.Sprintf(format, args...))
	this.Println(TRACE, args...)
}

func (this *ConsoleLogger) Debugf(format string, args ...interface{}) {
	args = this.addColor(this.color.Green(), fmt.Sprintf(format, args...))
	this.Println(DEBUG, args...)
}

func (this *ConsoleLogger) Infof(format string, args ...interface{}) {
	args = this.addColor(this.color.Cyan(), fmt.Sprintf(format, args...))
	this.Println(INFO, args...)
}

func (this *ConsoleLogger) Warnf(format string, args ...interface{}) {
	args = this.addColor(this.color.Magenta(), fmt.Sprintf(format, args...))
	this.Println(WARN, args...)
}

func (this *ConsoleLogger) Errorf(format string, args ...interface{}) {
	args = this.addColor(this.color.Yello(), fmt.Sprintf(format, args...))
	this.Println(ERROR, args...)
}

func (this *ConsoleLogger) Fatalf(exit bool, format string, args ...interface{}) {
	args = this.addColor(this.color.Red(), fmt.Sprintf(format, args...))
	this.Println(FATAL, args...)

	if exit {
		os.Exit(1)
	}
}

func (this *ConsoleLogger) Trace(args ...interface{}) {
	args = this.addColor(this.color.Blue(), args...)
	this.Println(TRACE, args...)
}

func (this *ConsoleLogger) Debug(args ...interface{}) {
	args = this.addColor(this.color.Green(), args...)
	this.Println(DEBUG, args...)
}

func (this *ConsoleLogger) Info(args ...interface{}) {
	args = this.addColor(this.color.Cyan(), args...)
	this.Println(INFO, args...)
}

func (this *ConsoleLogger) Warn(args ...interface{}) {
	args = this.addColor(this.color.Magenta(), args...)
	this.Println(WARN, args...)
}

func (this *ConsoleLogger) Error(args ...interface{}) {
	args = this.addColor(this.color.Yello(), args...)
	this.Println(ERROR, args...)
}

func (this *ConsoleLogger) Fatal(exit bool, args ...interface{}) {
	args = this.addColor(this.color.Red(), args...)
	this.Println(FATAL, args...)

	if exit {
		os.Exit(1)
	}
}

func (this *ConsoleLogger) rolling() {}
