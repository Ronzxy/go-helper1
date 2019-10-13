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
	"github.com/skygangsta/go-helper"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
)

// Logger level
type LogLevel int

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
const (
	ALL   = 0
	TRACE = 1
	DEBUG = 2
	INFO  = 3
	WARN  = 4
	ERROR = 5
	FATAL = 6
	OFF   = 7

	DefaultLogTimeFormat = "2006/01/02 15:04:05.000000"
)

var (
	DefaultWriter  = os.Stdout
	FileWithoutPKG = false
)

const (
	defaultSkipCallerDepth = 5
)

// Convert string level name to level
func ConvertString2Level(str string) LogLevel {
	var level LogLevel
	switch strings.ToUpper(str) {
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

func ConvertLevel2String(level LogLevel) string {
	str := ""

	switch level {
	case ALL:
		str = "ALL"
	case TRACE:
		str = "TRACE"
	case DEBUG:
		str = "DEBUG"
	case INFO:
		str = "INFO"
	case WARN:
		str = "WARN"
	case ERROR:
		str = "ERROR"
	case FATAL:
		str = "FATAL"
	default:
		str = "OFF"
	}

	return str
}

// Find the variable define in string
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

// Replace the variable define in properties
func VariableReplaceByConfig(str string) string {
	for {
		varPattern, varName := Variable("$", "([a-zA-Z_][0-9a-zA-Z_]+)", str)
		if varName == "" {
			// 没有发现变量，退出循环
			break
		}

		varName = RemoveEnterAndSpace(propertyMap[varName])

		str = strings.Replace(str, varPattern, varName, -1)
	}

	str = RemoveEnterAndSpace(str)

	return str
}

func RemoveEnterAndSpace(str string) string {
	str = strings.Replace(str, "\r\n", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Trim(str, " ")

	return str
}

func GetCaller(skip int) *runtime.Frame {
	pcs := make([]uintptr, 16)
	depth := runtime.Callers(skip, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	frame, _ := frames.Next()

	return &frame
}

// Get the package name by runtime.Frame.Function
func GetPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func GetFileName(frame *runtime.Frame) string {
	var (
		err error
		i   = 0
		s   = ""
	)
	if FileWithoutPKG {
		s, err = helper.Path.FileName(frame.File)
		if err == nil {
			return s
		}
	}

	s = GetPackageName(frame.Function)

	if s != "" {
		if s != "main" {
			if strings.Count(frame.File, s) > 0 {
				i = strings.LastIndex(frame.File, s)
				return frame.File[i:]
			}
		}
	}

	s = os.Getenv("GOPATH")

	if s != "" {
		if strings.Count(s, ":") > 0 {
			paths := strings.Split(s, ":")
			for _, s = range paths {
				s = path.Join(s, "src")

				if strings.Contains(frame.File, s) {
					i = len(s) + 1
				}
			}
		} else {
			s = path.Join(s, "src")

			if strings.Contains(frame.File, s) {
				i = len(s) + 1
			}
		}
	}

	if i <= 0 {
		var err error

		s, err = helper.Path.FileName(frame.File)
		if err == nil {
			return s
		}
	}

	return frame.File[i:]
}
