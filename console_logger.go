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
	"os"
)

type ConsoleLogger struct {
	*LoggerWriter

	consoleColor bool // use terminal color
}

var (
	defaultConsoleLogger = NewConsoleLogger(ALL)
)

func DefaultConsoleLogger() *ConsoleLogger {
	defaultConsoleLogger.name = "DefaultConsoleNoFilter"
	defaultConsoleLogger.closeFilter = true

	return defaultConsoleLogger
}

func NewConsoleLogger(level LogLevel) *ConsoleLogger {
	this := &ConsoleLogger{
		LoggerWriter: NewLoggerWriter(DefaultWriter, level),
		consoleColor: true,
	}

	this.SetFormatter(NewTextFormatter())

	return this
}

func (this *ConsoleLogger) addColor(color string, args ...interface{}) []interface{} {
	if this.consoleColor {
		format := append([]interface{}{}, color)
		args = append(format, args...)
		args = append(args, helper.ConsoleColor.Clear())
	}

	return args
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *ConsoleLogger) Tracef(format string, args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Blue(), fmt.Sprintf(format, args...))
	this.Write(TRACE, args...)
}

func (this *ConsoleLogger) Debugf(format string, args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Green(), fmt.Sprintf(format, args...))
	this.Write(DEBUG, args...)
}

func (this *ConsoleLogger) Infof(format string, args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Cyan(), fmt.Sprintf(format, args...))
	this.Write(INFO, args...)
}

func (this *ConsoleLogger) Warnf(format string, args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Magenta(), fmt.Sprintf(format, args...))
	this.Write(WARN, args...)
}

func (this *ConsoleLogger) Errorf(format string, args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Yello(), fmt.Sprintf(format, args...))
	this.Write(ERROR, args...)
}

func (this *ConsoleLogger) FatalfWithExit(exit bool, format string, args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Red(), fmt.Sprintf(format, args...))
	this.Write(FATAL, args...)

	if exit {
		os.Exit(1)
	}
}

func (this *ConsoleLogger) Trace(args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Blue(), args...)
	this.Write(TRACE, args...)
}

func (this *ConsoleLogger) Debug(args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Green(), args...)
	this.Write(DEBUG, args...)
}

func (this *ConsoleLogger) Info(args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Cyan(), args...)
	this.Write(INFO, args...)
}

func (this *ConsoleLogger) Warn(args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Magenta(), args...)
	this.Write(WARN, args...)
}

func (this *ConsoleLogger) Error(args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Yello(), args...)
	this.Write(ERROR, args...)
}

func (this *ConsoleLogger) FatalWithExit(exit bool, args ...interface{}) {
	args = this.addColor(helper.ConsoleColor.Red(), args...)
	this.Write(FATAL, args...)

	if exit {
		os.Exit(1)
	}
}

// Do nothing with implement interface Writer
func (this *ConsoleLogger) CheckRollingSize() {}
