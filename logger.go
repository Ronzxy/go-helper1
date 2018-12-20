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
	"github.com/robfig/cron"
	"os"
	"strings"
	"time"
)

var (
	config     *Config
	logger     *LoggerWriter
	properties map[string]string
	rolling    = false
	crontab    = cron.New()
)

type LoggerWriter struct {
	writers []Writer
}

func GetLogger() *LoggerWriter {
	return logger
}

func InitLogger(configFile string) error {
	var (
		err       error
		formatter Formatter
	)

	config, err = NewConfig(configFile)
	if err != nil {
		return err
	}

	properties = map[string]string{}
	logger = &LoggerWriter{}

	if config.Properties != nil {
		for _, v := range config.Properties {
			properties[v.Name] = v.Value
		}
	}

	if config.Loggers != nil {
		crontab.Start()

		for _, v := range config.Loggers {
			// 初始化 Formatter
			switch strings.ToLower(v.Format.Type) {
			case "text":
				formatter = NewTextFormatterWithFormat(v.Format.Value)
			case "json":
				formatter = NewJSONFormatter()
			default:
				formatter = NewTextFormatter()
			}

			switch v.Target {
			case "STDOUT":
				{
					// Console Log
					var (
						consoleLogger *ConsoleLogger
					)

					consoleLogger = NewConsoleLogger(ConvertString2Level(v.Level.Allow))
					consoleLogger.SetDenyLevel(ConvertString2Level(v.Level.Deny))
					consoleLogger.SetFormatter(formatter)

					logger.writers = append(logger.writers, consoleLogger)
				}
			case "FILE":
				{
					// File Log
					var (
						fileLogger *FileLogger
					)

					fileLogger, err = NewFileLoggerWithConfig(v)
					if err == nil {
						if v.Rolling.TimeBased > 0 {
							err = crontab.AddFunc(fmt.Sprintf("0 0 */%d * * ?", v.Rolling.TimeBased), fileLogger.RollingFile)
							if err != nil {
								Errorf("create cron error %s", err.Error())
							}
						}

						fileLogger.SetFormatter(formatter)

						logger.writers = append(logger.writers, fileLogger)
					} else {
						DefaultConsoleLogger().Error(err.Error())
					}
				}
			default:
				DefaultConsoleLogger().Warnf("unsupported log target %s", v.Target)
			}
		}
		// rolling log file
		StartRolling()
	}

	return nil
}

func rollingFileSize() {
	if config.RollingInterval <= 0 {
		config.RollingInterval = 60
	}

	for {
		select {
		case <-time.After(time.Duration(config.RollingInterval) * time.Second):
			// rolling file
			if logger.writers != nil && rolling {
				for _, v := range logger.writers {
					v.CheckRollingSize()
				}
			}
		}

		// if disable rolling break loop
		if !rolling {
			break
		}
	}
}

func StartRolling() {
	rolling = true
	crontab.Start()

	go rollingFileSize()
}

func StopRolling() {
	rolling = false
	crontab.Stop()
}

func Tracef(format string, args ...interface{}) {
	for _, value := range logger.writers {
		value.Tracef(format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	for _, value := range logger.writers {
		value.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	for _, value := range logger.writers {
		value.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	for _, value := range logger.writers {
		value.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	for _, value := range logger.writers {
		value.Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	for _, value := range logger.writers {
		value.Fatalf(false, format, args...)
	}
	os.Exit(-1)
}

func Trace(args ...interface{}) {
	for _, value := range logger.writers {
		value.Trace(args...)
	}
}

func Debug(args ...interface{}) {
	for _, value := range logger.writers {
		value.Debug(args...)
	}
}

func Info(args ...interface{}) {
	for _, value := range logger.writers {
		value.Info(args...)
	}
}

func Warn(args ...interface{}) {
	for _, value := range logger.writers {
		value.Warn(args...)
	}
}

func Error(args ...interface{}) {
	for _, value := range logger.writers {
		value.Error(args...)
	}
}

func Fatal(args ...interface{}) {
	for _, value := range logger.writers {
		value.Fatal(false, args...)
	}
	os.Exit(-1)
}
