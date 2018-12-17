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
	"github.com/skygangsta/go-utils"
	"io"
	"log"
)

type LogWriter struct {
	allowLevel      int
	denyLevel       int
	name            string
	logger          *log.Logger
	formatter       Formatter
	skipCallerDepth int
}

func NewLogWriter(w io.Writer, level int) *LogWriter {
	this := &LogWriter{
		allowLevel:      level,
		denyLevel:       OFF,
		name:            util.NewPath().WorkName(),
		logger:          log.New(w, "", log.LUTC),
		skipCallerDepth: defaultSkipCallerDepth,
	}

	this.SetFormatter(NewTextFormatter())

	return this
}

func (this *LogWriter) SetDenyLevel(level int) {
	if level > this.denyLevel {
		this.allowLevel = OFF
	} else {
		this.denyLevel = level
	}
}

func (this *LogWriter) SetSkipCallerDepth(skipCallerDepth int) {
	this.skipCallerDepth = skipCallerDepth
}

func (this *LogWriter) SetName(name string) {
	this.name = name
}

func (this *LogWriter) SetWriter(w io.Writer) {
	this.logger.SetOutput(w)
}

func (this *LogWriter) SetFormatter(formatter Formatter) {
	this.formatter = formatter
}

func (this *LogWriter) Println(level int, args ...interface{}) {
	if this.allowLevel <= level {
		if this.denyLevel > level {
			var (
				data = map[string]interface{}{}
			)

			frame := GetCaller(this.skipCallerDepth)

			data["Name"] = this.name
			data["Level"] = ConvertLevel2String(level)
			data["Line"] = frame.Line
			data["PackageName"] = GetPackageName(frame.Function)

			//file := frame.File
			//
			//if strings.Contains(file, "/src/") {
			//	files := strings.Split(file, "/src/")
			//	file = files[1]
			//}
			//
			data["File"] = GetFileName(frame.File)

			//fmt.Println("GOPATH", sys.DefaultGoroot)

			this.logger.Println(this.formatter.Message(data, args...))
		}
	}
}
