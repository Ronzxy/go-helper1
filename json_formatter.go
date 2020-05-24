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
	"encoding/json"
	"time"
)

type JSONFormatter struct {
	Indent bool
}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		Indent: false,
	}
}

func (this *JSONFormatter) Message(data map[string]interface{}, args ...interface{}) string {
	var (
		buf []byte
		err error
	)

	data["Time"] = time.Now().Format(DefaultLogTimeFormat)
	data["Message"] = args

	if this.Indent {
		buf, err = json.MarshalIndent(data, "", "")
	} else {
		buf, err = json.Marshal(data)
	}

	if err != nil {
		Errorf("")
	}

	return string(buf)
}
