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
	"testing"
)

var (
	err = InitLogger("example/logger.xml")
)

func TestLogger(t *testing.T) {
	if err != nil {
		defaultLogger.Error(err.Error())
		return
	}

	Trace("Test Logger trace message")
	Debug("Test Logger debug message")
	Info("Test Logger info message")
	Warn("Test Logger warn message")
	Error("Test Logger error message")

	t.Log("Test Logger finished.")
}

func BenchmarkLogger(b *testing.B) {
	if err != nil {
		defaultLogger.Error(err.Error())
		return
	}

	for i := 0; i < b.N; i++ {
		Trace("Benchmark Logger trace message")
		Debug("Benchmark Logger debug message")
		Info("Benchmark Logger info message")
		Warn("Benchmark Logger warn message")
		Error("Benchmark Logger error message")

		b.Log("Benchmark Logger finished.")
	}
}
