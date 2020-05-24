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
	"errors"
	"fmt"
	"github.com/ronzxy/go-helper"
	"os"
	"path"
	"strings"
	"time"
)

type FileLogger struct {
	*LoggerWriter

	writer     *os.File
	config     Logger
	storeIndex int
	storeFirst int
}

func NewFileLogger(level LogLevel, logFile string) (*FileLogger, error) {
	var (
		fileLogger = &FileLogger{}
		err        error
	)

	err = fileLogger.createDir(logFile)
	if err != nil {
		return nil, err
	}

	fileLogger.writer, err = os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error: Open config file %v", err))
	}

	fileLogger.LoggerWriter = NewLoggerWriter(fileLogger.writer, level)

	return fileLogger, nil
}

func NewFileLoggerWithConfig(v Logger) (*FileLogger, error) {
	var (
		fileLogger *FileLogger
		err        error
	)

	fileLogger, err = NewFileLogger(ConvertString2Level(v.Level.Allow), fileLogger.variableReplacer(v.FileName))
	if err == nil {
		fileLogger.SetDenyLevel(ConvertString2Level(v.Level.Deny))
		fileLogger.config = v
	}

	return fileLogger, err
}
// TODO: fileroll error
func (this *FileLogger) createDir(fileName string) error {
	var (
		isExist  bool
		filePath string
		err      error
	)
	filePath, err = helper.Path.Dir(fileName)
	if err != nil {
		return err
	}

	isExist, err = helper.Path.IsExist(filePath)
	if err != nil {
		return err
	}

	if !isExist {
		err = helper.Path.CreateDir(filePath, 0755)
		if err != nil && err.Error() != "file exists" {
			return err
		}
	}

	return nil
}

// Replace ${([a-zA-Z_][0-9a-zA-Z_]+)} and %{([0-9a-zA-Z_:-]+} variable
func (this *FileLogger) variableReplacer(fileName string) string {
	fileName = this.variableReplaceByConfig(fileName)
	fileName = this.variableReplaceBySystem(fileName)

	return fileName
}

// Replace ${([a-zA-Z_][0-9a-zA-Z_]+)} user defined variable
func (this *FileLogger) variableReplaceByConfig(str string) string {
	return VariableReplaceByConfig(str)
}

// Replace %{([a-zA-Z_][0-9a-zA-Z_:-]*} system defined variable
func (this *FileLogger) variableReplaceBySystem(str string) string {
	var (
		varPattern string
		varName    string
		vars       = make([]string, 2)
	)
	for {
		varPattern, varName = Variable("%", "([a-zA-Z_][0-9a-zA-Z_/:-]*)", str)
		if varName == "" {
			// no variable, exit loop
			break
		}

		if strings.Contains(varName, ":") {
			vars = strings.Split(varName, ":")
		} else {
			vars[0] = varName
			vars[1] = ""
		}

		switch strings.ToLower(vars[0]) {
		case "date":
			{
				if vars[1] != "" {
					varName = helper.Time.Format(strings.Join(vars[1:], ":"), time.Now())
				} else {
					varName = time.Now().Format(DefaultLogTimeFormat)
				}
			}
		case "i":
			{
				this.storeIndex = this.storeIndex + 1
				varName = fmt.Sprintf("%02d", this.storeIndex)
			}
		default:
			{
				Errorf("unsupported function %s", vars[0])
				varName = ""
			}
		}

		str = strings.Replace(str, varPattern, varName, -1)
	}

	return str
}

func (this *FileLogger) CheckRollingSize() {
	// if XMLName is empty, maybe not initial from config file
	// and the config maybe be empty
	if this.config.XMLName.Local != "" {
		if this.config.Rolling.SizeBased <= 0 {
			this.config.Rolling.SizeBased = 1
		}

		fileInfo, err := this.writer.Stat()
		if err == nil {
			// check file size
			if fileInfo.Size() >= int64(this.config.Rolling.SizeBased)*1024*1024 {
				this.RollingFile()
			}
		} else {
			Errorf("check file error with %s", err.Error())
		}
	}
}

// Rolling a new file to write logger
func (this *FileLogger) RollingFile() {
	var (
		storeFileName string
		newFileName   string
		err           error
		isExist       bool
	)

	for {
		storeFileName = this.variableReplacer(this.config.FilePattern)

		isExist, err = helper.Path.IsExist(storeFileName)
		if err != nil {
			Errorf("check file exist error: %s", err.Error())
			return
		}

		if isExist {
			// The file already exists,
			// continue execute to get the next file name
			continue
		}

		break // rolling complete, exit loop
	}
	// if storage file same name with log file,
	// there is no need to rolling the file
	if storeFileName == this.writer.Name() {
		Trace("store file name is same as log file, rolling file will be ignored")
		return
	}

	// if log file has no content,
	// there is no need to rolling the file
	fileInfo, err := this.writer.Stat()
	if err != nil {
		Errorf("check log file error with %s", err.Error())
		return
	}

	// check file size
	if fileInfo.Size() <= 0 {
		return
	}

	err = this.createDir(storeFileName)
	if err != nil {
		Errorf("create storage path error: %s", err.Error())
		return
	}

	newPath, err := helper.Path.Dir(this.writer.Name())
	if err != nil {
		Errorf("get log file base path error: %s", err.Error())
		return
	}

	newName, err := helper.Path.FileName(storeFileName)
	if err != nil {
		Errorf("get log file name error: %s", err.Error())
		return
	}

	newFileName = path.Join(newPath, newName)
	err = os.Rename(this.writer.Name(), newFileName)
	if err != nil {
		Errorf("rename file error: %s", err.Error())
		return
	}

	// create a new log file
	file, err := os.OpenFile(this.writer.Name(), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		Errorf("create log file error: %s", err.Error())
		return
	}

	this.LoggerWriter.NewLogger(file)

	this.writer.Close()
	this.writer = file

	switch this.config.Compress {
	case "gzip":
		{
			// gzip file to store path
			err = helper.File.GZipFile(newFileName, storeFileName+".gz")
			if err != nil {
				Errorf("gzip file error: %s", err.Error())
				return
			}
		}
	default:
		{
			// copyFile file to store path
			_, err = helper.File.CopyFile(newFileName, storeFileName)
			if err != nil {
				Errorf("copy log file error: %s", err.Error())
				return
			}
		}
	}

	// check keep count
	this.keepFile()
}

func (this *FileLogger) keepFile() {
	var (
		storeFile string
		logFile   = this.variableReplacer(this.config.FileName)
	)
	if (this.storeIndex - this.storeFirst - this.config.Rolling.KeepCount) >= 0 {
		for i := this.storeFirst; i <= this.storeIndex-this.config.Rolling.KeepCount; i++ {

			this.storeFirst = i + 1

			if i == 0 {
				continue
			}

			storeFile = strings.Replace(this.config.FilePattern, "%{i}", fmt.Sprintf("%02d", i), -1)
			storeFile = this.variableReplacer(storeFile)

			inx := strings.LastIndex(storeFile, ".")
			ext := storeFile[inx:]
			if ext == ".gz" {
				storeFile = storeFile[:inx]
			}

			newPath, err := helper.Path.Dir(logFile)
			if err != nil {
				Errorf("get log file base path error: %s", err.Error())
				return
			}

			newName, err := helper.Path.FileName(storeFile)
			if err != nil {
				Errorf("get log file name error: %s", err.Error())
				return
			}

			storeFile = path.Join(newPath, newName)

			err = os.Remove(storeFile)
			if err != nil {
				Errorf("delete file error: %s", err.Error())
			}
		}

	}
}
