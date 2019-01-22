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

import "github.com/go-xorm/core"

// logger writer interface
type Writer interface {
	Tracef(format string, args ...interface{})

	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(exit bool, format string, args ...interface{})

	Trace(args ...interface{})

	Debug(args ...interface{})

	Info(args ...interface{})

	Warn(args ...interface{})

	Error(args ...interface{})

	Fatal(exit bool, args ...interface{})

	CheckRollingSize()

	/*
	   Include xorm logger
	*/

	Level() core.LogLevel

	SetLevel(l core.LogLevel)

	ShowSQL(show ...bool)

	IsShowSQL() bool
}
