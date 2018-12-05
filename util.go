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
	"io"
	"os"
	"regexp"
	"strings"
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

// Convert string level name to level
func ConvertLevelName(levelName string) int {
	level := 0
	switch strings.ToUpper(levelName) {
	case "ALL":
		level = ALL
	case "TRACE":
		level = TRACE
	case "DEBUG":
		level = DEBUG
	case "INFO":
		level = INFO
	case "WARN":
		level = WARN
	case "ERROR":
		level = ERROR
	case "FATAL":
		level = FATAL
	default:
		level = OFF
	}

	return level
}

// Find the variable define in properties
func Variable(varchar, pattern, str string) (string, string) {
	var (
		varPattern string
		varName    string
	)
	r := regexp.MustCompile(fmt.Sprintf(`\%s{%s}`, varchar, pattern))
	varPattern = r.FindString(str)

	r = regexp.MustCompile(pattern)
	varName = r.FindString(varPattern)

	return varPattern, varName
}
