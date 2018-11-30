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
	"github.com/skygangsta/go-logger/config"
	"regexp"
	"strings"
)

type Logger struct {
	logger []Writer

	properties map[string]string
}

var (
	this *Logger
)

func GetLogger() *Logger {
	return this
}

func InitLogger(configFile string) error {
	var (
		err    error
		config *config.Config
	)

	this = &Logger{
		properties: map[string]string{},
	}

	config, err = NewConfig(configFile)
	if err != nil {
		return err
	}

	if config.Properties != nil {
		for _, v := range config.Properties {
			this.properties[v.Name] = v.Value
		}
	}

	if config.Loggers != nil {
		for _, v := range config.Loggers {
			switch v.Target {
			case "STDOUT":
				{
					// Console Log
					var (
						consoleLogger *ConsoleLogger
					)

					consoleLogger = NewConsoleLogger(convertLevelName(v.Level.Allow))
					consoleLogger.SetDenyLevel(convertLevelName(v.Level.Deny))

					this.logger = append(this.logger, consoleLogger)
				}
			case "FILE":
				{
					// File Log
					var (
						fileLogger *FileLogger
					)

					fileLogger, err = NewFileLogger(convertLevelName(v.Level.Allow), this.varReplacer(v.FileName), 0644)
					if err == nil {
						fileLogger.SetDenyLevel(convertLevelName(v.Level.Deny))

						this.logger = append(this.logger, fileLogger)
					} else {
						DefaultConsoleLogger().Error(err.Error())
					}
				}
			default:
				DefaultConsoleLogger().Warnf("unsupported log target %s", v.Target)
			}
		}
	}

	return nil
}

func (this *Logger) varReplacer(str string) string {
	r := regexp.MustCompile(`\${.*}`)
	fileName := r.FindString(str)

	r = regexp.MustCompile(`([A-Z_A-Z]+)`)
	varName := r.FindString(fileName)

	varName = this.properties[varName]
	varName = strings.Replace(varName, "\n", "", -1)
	varName = strings.Trim(varName, " ")

	fileName = strings.Replace(str, fileName, varName, -1)
	fileName = strings.Replace(fileName, "//", "/", -1)

	return fileName
}

func Tracef(format string, v ...interface{}) {
	for _, value := range this.logger {
		value.Tracef(format, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	for _, value := range this.logger {
		value.Debugf(format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	for _, value := range this.logger {
		value.Infof(format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	for _, value := range this.logger {
		value.Warnf(format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	for _, value := range this.logger {
		value.Errorf(format, v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	for _, value := range this.logger {
		value.Fatalf(format, v...)
	}
}

func Trace(message string) {
	Tracef("%s", message)
}

func Debug(message string) {
	Debugf("%s", message)
}

func Info(message string) {
	Infof("%s", message)
}

func Warn(message string) {
	Warnf("%s", message)
}

func Error(message string) {
	Errorf("%s", message)
}

func Fatal(message string) {
	Fatalf("%s", message)
}

// ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
func convertLevelName(levelName string) int {
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
