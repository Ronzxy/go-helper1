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
	"github.com/robfig/cron"
	"os"
	"strings"
	"time"
)

var (
	config      *Config
	job         = cron.New()
	propertyMap = map[string]string{}
	writerMap   = map[string]Writer{}
	rolling     = false
	initialized = false
)

func Init(configFile string) error {
	var (
		err error
	)

	config, err = NewConfig(configFile)
	if err != nil {
		return err
	}

	initProperties()

	job.Start()

	if config.Loggers != nil {
		// Initialize the Writer of the package filter reference
		if config.PackageFilters != nil {
			for _, filter := range config.PackageFilters {
				if len(filter.Loggers) <= 0 {
					break
				}

				for _, loggerName := range filter.Loggers {
					if writerMap[loggerName] != nil {
						continue
					}

					logger := initLogger(loggerName)
					if logger != nil {
						writerMap[loggerName] = logger
					}
				}
			}
		}

		// Initialize the Writer of the default filter reference
		if len(config.DefaultFilter.Loggers) > 0 {
			for _, loggerName := range config.DefaultFilter.Loggers {
				if writerMap[loggerName] != nil {
					continue
				}

				logger := initLogger(loggerName)
				if logger != nil {
					writerMap[loggerName] = logger
				}
			}
		}

		initialized = true
		// rolling log file
		StartRolling()
	}

	return nil
}

func GetByPackage(packageName string) []Writer {
	if !Initialized() {
		return nil
	}

	if config.PackageFilters == nil {
		return nil
	}

	var writers []Writer

	for _, filter := range config.PackageFilters {
		//packageName := GetPackageName(frame.Function)
		if packageName == filter.Name {
			if len(filter.Loggers) == 0 {
				// No Logger define
				return nil
			}

			for _, name := range filter.Loggers {
				writer := writerMap[name]
				if writer == nil {
					continue
				}
				writers = append(writers, writer)
			}
		}
	}

	return writers
}

func initProperties() {
	if len(config.Properties) > 0 {
		for _, v := range config.Properties {
			propertyMap[v.Name] = v.Value
		}
	}
}

func initLogger(name string) Writer {
	var (
		err       error
		formatter Formatter
	)

	for _, v := range config.Loggers {
		if name == v.Name {
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
					consoleLogger.name = v.Name

					return consoleLogger
				}
			case "FILE":
				{
					// File Log
					var (
						fileLogger *FileLogger
					)

					fileLogger, err = NewFileLoggerWithConfig(v)
					if err == nil {
						if v.Rolling.TimeBased == "" {
							v.Rolling.TimeBased = "@daily"
						}

						err = job.AddFunc(v.Rolling.TimeBased, fileLogger.RollingFile)
						if err != nil {
							Errorf("create cron error %s", err.Error())
						}

						fileLogger.SetFormatter(formatter)
						fileLogger.name = v.Name

						return fileLogger
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

func rollingFileSize() {
	if config.RollingInterval <= 0 {
		config.RollingInterval = 60
	}

	for {
		select {
		case <-time.After(time.Duration(config.RollingInterval) * time.Second):
			// rolling file
			if writerMap != nil && rolling {
				for _, v := range writerMap {
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
	job.Start()

	go rollingFileSize()
}

func StopRolling() {
	rolling = false
	job.Stop()
}

func Initialized() bool {
	return initialized
}

func Tracef(format string, args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Tracef(format, args...)
		}
	} else {
		DefaultConsoleLogger().Tracef(format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Debugf(format, args...)
		}
	} else {
		DefaultConsoleLogger().Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Infof(format, args...)
		}
	} else {
		DefaultConsoleLogger().Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Warnf(format, args...)
		}
	} else {
		DefaultConsoleLogger().Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Errorf(format, args...)
		}
	} else {
		DefaultConsoleLogger().Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.FatalfWithExit(false, format, args...)
		}
	} else {
		DefaultConsoleLogger().FatalfWithExit(false, format, args...)
	}
	os.Exit(-1)
}

func Trace(args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Trace(args...)
		}
	} else {
		DefaultConsoleLogger().Trace(args...)
	}
}

func Debug(args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Debug(args...)
		}
	} else {
		DefaultConsoleLogger().Debug(args...)
	}
}

func Info(args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Info(args...)
		}
	} else {
		DefaultConsoleLogger().Info(args...)
	}
}

func Warn(args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Warn(args...)
		}
	} else {
		DefaultConsoleLogger().Warn(args...)
	}
}

func Error(args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.Error(args...)
		}
	} else {
		DefaultConsoleLogger().Error(args...)
	}
}

func Fatal(args ...interface{}) {
	if Initialized() {
		for _, value := range writerMap {
			value.FatalWithExit(false, args...)
		}
	} else {
		DefaultConsoleLogger().FatalWithExit(false, args...)
	}
	os.Exit(-1)
}
