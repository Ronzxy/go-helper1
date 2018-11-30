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
	"errors"
	"fmt"
	"github.com/skygangsta/go-utils"
	"os"
	"time"
)

type FileLogger struct {
	*LogWriter
	writer *os.File
}

func NewFileLogger(level int, configFile string, perm uint32) (*FileLogger, error) {
	var (
		this = &FileLogger{}
		err  error
		p    = util.NewPath()
	)

	filePath, err := p.Dir(configFile)
	if err != nil {
		return nil, err
	}

	err = p.Create(filePath, 0755)
	if err != nil {
		return nil, err
	}

	this.writer, err = os.OpenFile(configFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error: Open config file %v", err))
	}

	this.LogWriter = NewLogWriter(this.writer, level)

	return this, nil
}

func (this *FileLogger) rolling() {
	for {
		this.Infof("ddd %v", time.Now())
		time.Sleep(10 * time.Second)
	}
}
