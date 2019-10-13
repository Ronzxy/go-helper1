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
	DefaultConsoleLogger().SetSkipCallerDepth(4)
	fileLogger, err := NewFileLogger(ALL, fmt.Sprintf("logs/fileLogger-test-%s.log", helper.Time.Format("yyyy-mm-dd-HHMMSS.ns", time.Now())))
	if err != nil {
		DefaultConsoleLogger().Error(err.Error())
		return
	}

	fileLogger.closeFilter = true
	fileLogger.SetSkipCallerDepth(4)
	fileLogger.Trace("Test FileLogger trace message")
	fileLogger.Debug("Test FileLogger debug message")
	fileLogger.Info("Test FileLogger info message")
	fileLogger.Warn("Test FileLogger warn message")
	fileLogger.Error("Test FileLogger error message")

	t.Log("Test FileLogger finished.")
}

func BenchmarkFileLogger(b *testing.B) {
	DefaultConsoleLogger().SetSkipCallerDepth(4)
	fileLogger, err := NewFileLogger(ALL, fmt.Sprintf("logs/fileLogger-bench-%s.log", helper.Time.Format("yyyy-mm-dd-HHMMSS.ns", time.Now())))
	if err != nil {
		DefaultConsoleLogger().Error(err.Error())
		return
	}
	fileLogger.closeFilter = true
	fileLogger.SetSkipCallerDepth(4)

	for i := 0; i < b.N; i++ {
		fileLogger.Trace("Benchmark FileLogger trace message")
		fileLogger.Debug("Benchmark FileLogger debug message")
		fileLogger.Info("Benchmark FileLogger info message")
		fileLogger.Warn("Benchmark FileLogger warn message")
		fileLogger.Error("Benchmark FileLogger error message")

		b.Log("Benchmark FileLogger finished.")
	}
}
