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

// logger writer interface
type Writer interface {
	Tracef(format string, v ...interface{})

	Debugf(format string, v ...interface{})

	Infof(format string, v ...interface{})

	Warnf(format string, v ...interface{})

	Errorf(format string, v ...interface{})

	Fatalf(format string, v ...interface{})

	Trace(message string)

	Debug(message string)

	Info(message string)

	Warn(message string)

	Error(message string)

	Fatal(message string)

	rolling()
}
