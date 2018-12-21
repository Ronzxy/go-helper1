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

// benchmark: go test -test.bench=".*"
package logger

import (
	"fmt"
	"github.com/skygangsta/go-helper"
	"testing"
	"time"
)

func TestFileLogger(t *testing.T) {
	filelogger, err := NewFileLogger(ALL, fmt.Sprintf("logs/fileLogger-test-%s.log", helper.NewTimeHelper().Format(time.Now(), "yyyy-mm-dd-HHMMSS.ns")))
	if err != nil {
		defaultLogger.Error(err.Error())
		return
	}

	filelogger.Trace("Test FileLogger trace message")
	filelogger.Debug("Test FileLogger debug message")
	filelogger.Info("Test FileLogger info message")
	filelogger.Warn("Test FileLogger warn message")
	filelogger.Error("Test FileLogger error message")

	t.Log("Test FileLogger finished.")
}

func BenchmarkFileLogger(b *testing.B) {
	filelogger, err := NewFileLogger(ALL, fmt.Sprintf("logs/fileLogger-bench-%s.log", helper.NewTimeHelper().Format(time.Now(), "yyyy-mm-dd-HHMMSS.ns")))
	if err != nil {
		defaultLogger.Error(err.Error())
		return
	}

	for i := 0; i < b.N; i++ {
		filelogger.Trace("Benchmark FileLogger trace message")
		filelogger.Debug("Benchmark FileLogger debug message")
		filelogger.Info("Benchmark FileLogger info message")
		filelogger.Warn("Benchmark FileLogger warn message")
		filelogger.Error("Benchmark FileLogger error message")

		b.Log("Benchmark FileLogger finished.")
	}
}
