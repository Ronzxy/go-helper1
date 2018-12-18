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
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/skygangsta/go-utils"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type FileLogger struct {
	*LogWriter

	writer     *os.File
	config     Logger
	storeIndex int
	storeFirst int
}

func NewFileLogger(level int, configFile string, perm uint32) (*FileLogger, error) {
	var (
		this = &FileLogger{}
		err  error
		//p    = util.NewPath()
	)

	err = this.createDir(configFile)
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

func NewFileLoggerWithConfig(v Logger) (*FileLogger, error) {
	var (
		fileLogger *FileLogger
		err        error
	)

	fileLogger, err = NewFileLogger(ConvertString2Level(v.Level.Allow), fileLogger.varReplacer(v.FileName), 0644)
	if err == nil {
		fileLogger.SetDenyLevel(ConvertString2Level(v.Level.Deny))
		fileLogger.config = v
	}

	return fileLogger, err
}

func (this *FileLogger) createDir(fileName string) error {
	var (
		p        = util.NewPath()
		b        bool
		s        []string
		filePath string
		err      error
	)
	filePath, err = p.Dir(fileName)
	if err != nil {
		return err
	}

	b, err = p.IsExist(filePath)
	if err != nil {
		return err
	}

	if !b {
		s = p.Split(filePath)

		filePath = ""
		for _, v := range s {
			if v == "" {
				filePath = "/"
			} else {
				filePath = path.Join(filePath, v)

				err = p.Create(filePath, 0755)
				if err != nil {
					b, _ = p.IsExist(filePath)
					if !b {
						return err
					}
				}
			}
		}
	}

	return nil
}

// Replace ${([a-zA-Z_][0-9a-zA-Z_]+)} and %{([0-9a-zA-Z_:-]+} variable
func (this *FileLogger) varReplacer(fileName string) string {
	fileName = this.varDefineByConfig(fileName)
	fileName = this.varDefineBySystem(fileName)

	return fileName
}

// Replace ${([a-zA-Z_][0-9a-zA-Z_]+)} user defined variable
func (this *FileLogger) varDefineByConfig(str string) string {

	return VariableReplaceByConfig(str)
}

// Replace %{([a-zA-Z_][0-9a-zA-Z_:-]*} system defined variable
func (this *FileLogger) varDefineBySystem(str string) string {
	for {
		varPattern, varName := Variable("%", "([a-zA-Z_][0-9a-zA-Z_/:-]*)", str)
		if varName == "" {
			// no variable, exit loop
			break
		}

		if strings.Contains(varName, ":") { // %{date:Y-m-d} variable
			slice := strings.Split(varName, ":")
			if len(slice) != 2 {
				Error("unsupported function define")
				varName = ""
				continue
			}
			switch slice[0] {
			case "date":
				{
					varName = util.NewDate().Format(time.Now(), slice[1])
				}
			default:
				{
					Errorf("unsupported function %s", slice[0])
					varName = ""
				}
			}
		} else { // %{i}
			switch varName {
			case "i":
				{
					this.storeIndex = this.storeIndex + 1
					varName = fmt.Sprintf("%02d", this.storeIndex)
				}
			default:
				{
					Errorf("unsupported variable %s", varName)
					varName = ""
				}
			}
		}

		str = strings.Replace(str, varPattern, varName, -1)
	}

	return str
}

func (this *FileLogger) copyFile(src, dst string) (int64, error) {
	stat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !stat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	reader, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	writer, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer writer.Close()

	return io.Copy(writer, reader)
}

// gzip file to store path
func (this *FileLogger) gzipFile(read, write string) error {
	reader, err := os.Open(read)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := os.Create(write)
	if err != nil {
		return err
	}
	defer writer.Close()

	gw, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gw.Close()

	var bytes = make([]byte, 4096)
	for {
		n, err := reader.Read(bytes)
		if err != nil {
			if err.Error() != "EOF" {
				return err
			}

			break
		}

		gw.Write(bytes[:n])
		gw.Flush()
	}

	return nil
}

func (this *FileLogger) rolling() {
	// if XMLName is empty, maybe not initial from config
	if this.config.XMLName.Local != "" {
		fileInfo, err := this.writer.Stat()
		if err == nil {
			// check file size
			if fileInfo.Size() >= int64(this.config.Rolling.SizeBased)*1024*1024 {
				this.rollingFile()
			}
		} else {
			Errorf("check file error with %s", err.Error())
		}
	}
}

func (this *FileLogger) rollingFile() {
	var (
		storeFile string
		newFile   string
		err       error
		b         bool
		logFile   = this.varReplacer(this.config.FileName)
		utilPath  = util.NewPath()
	)
	for {
		storeFile = this.varReplacer(this.config.FilePattern)

		b, err = utilPath.IsExist(storeFile)
		if err != nil {
			Errorf("check file exist error: %s", err.Error())
			return
		}

		if b {
			// The file already exists, continue execute to get the next file name
			continue
		}

		break // rolling complete, exit loop
	}
	// if storage file same name with log file,
	// there is no need to rolling the file
	if storeFile == logFile {
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

	err = this.createDir(storeFile)
	if err != nil {
		Errorf("create storage path error: %s", err.Error())
		return
	}

	inx := strings.LastIndex(storeFile, ".")
	ext := storeFile[inx:]

	if ext == ".gz" {
		newFile = storeFile[:inx]
	}

	newPath, err := utilPath.Dir(logFile)
	if err != nil {
		Errorf("get log file base path error: %s", err.Error())
		return
	}

	newName, err := utilPath.FileName(newFile)
	if err != nil {
		Errorf("get log file name error: %s", err.Error())
		return
	}

	newFile = path.Join(newPath, newName)
	err = os.Rename(logFile, newFile)
	if err != nil {
		Errorf("rename file error: %s", err.Error())
		return
	}

	// create a new log file
	this.writer, err = os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		Errorf("create log file error: %s", err.Error())
		return
	}

	this.LogWriter = NewLogWriter(this.writer, ConvertString2Level(this.config.Level.Allow))
	this.SetDenyLevel(ConvertString2Level(this.config.Level.Deny))

	if ext == ".gz" {
		// gzip file to store path
		err = this.gzipFile(newFile, storeFile)
		if err != nil {
			Errorf("gzip file error: %s", err.Error())
			return
		}
	} else {
		// copyFile file to store path
		_, err = this.copyFile(newFile, storeFile)
		if err != nil {
			Errorf("copy log file error: %s", err.Error())
			return
		}
	}
	// check keep count
	this.autoKeepFile()
}

func (this *FileLogger) autoKeepFile() {
	var (
		storeFile string
		logFile   = this.varReplacer(this.config.FileName)
		utilPath  = util.NewPath()
	)
	if (this.storeIndex - this.storeFirst - this.config.Rolling.KeepCount) >= 0 {
		for i := this.storeFirst; i <= this.storeIndex-this.config.Rolling.KeepCount; i++ {

			this.storeFirst = i + 1

			if i == 0 {
				continue
			}

			storeFile = strings.Replace(this.config.FilePattern, "%{i}", fmt.Sprintf("%02d", i), -1)
			storeFile = this.varReplacer(storeFile)

			inx := strings.LastIndex(storeFile, ".")
			ext := storeFile[inx:]
			if ext == ".gz" {
				storeFile = storeFile[:inx]
			}

			newPath, err := utilPath.Dir(logFile)
			if err != nil {
				Errorf("get log file base path error: %s", err.Error())
				return
			}

			newName, err := utilPath.FileName(storeFile)
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

/***** LogWrite *****/
// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *FileLogger) Tracef(format string, args ...interface{}) {
	this.Println(TRACE, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Debugf(format string, args ...interface{}) {
	this.Println(DEBUG, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Infof(format string, args ...interface{}) {
	this.Println(INFO, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Warnf(format string, args ...interface{}) {
	this.Println(WARN, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Errorf(format string, args ...interface{}) {
	this.Println(ERROR, fmt.Sprintf(format, args...))
}

func (this *FileLogger) Fatalf(exit bool, format string, args ...interface{}) {
	this.Println(FATAL, fmt.Sprintf(format, args...))

	if exit {
		os.Exit(1)
	}
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func (this *FileLogger) Trace(args ...interface{}) {
	this.Println(TRACE, args...)
}

func (this *FileLogger) Debug(args ...interface{}) {
	this.Println(DEBUG, args...)
}

func (this *FileLogger) Info(args ...interface{}) {
	this.Println(INFO, args...)
}

func (this *FileLogger) Warn(args ...interface{}) {
	this.Println(WARN, args...)
}

func (this *FileLogger) Error(args ...interface{}) {
	this.Println(ERROR, args...)
}

func (this *FileLogger) Fatal(exit bool, args ...interface{}) {
	this.Println(FATAL, args...)

	if exit {
		os.Exit(1)
	}
}
