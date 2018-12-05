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
	"time"
)

var (
	logger  *LoggerWriter
	config  *Config
	rolling = false
	crontab = cron.New()
)

type LoggerWriter struct {
	writers []Writer

	properties map[string]string
}

func GetLogger() *LoggerWriter {
	return logger
}

func InitLogger(configFile string) error {
	var (
		err error
	)

	logger = &LoggerWriter{
		properties: map[string]string{},
	}

	config, err = NewConfig(configFile)
	if err != nil {
		return err
	}

	if config.Properties != nil {
		for _, v := range config.Properties {
			logger.properties[v.Name] = v.Value
		}
	}

	if config.Loggers != nil {
		crontab.Start()

		for _, v := range config.Loggers {
			switch v.Target {
			case "STDOUT":
				{
					// Console Log
					var (
						consoleLogger *ConsoleLogger
					)

					consoleLogger = NewConsoleLogger(ConvertLevelName(v.Level.Allow))
					consoleLogger.SetDenyLevel(ConvertLevelName(v.Level.Deny))

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
							err = crontab.AddFunc(fmt.Sprintf("0 0 */%d * * ?", v.Rolling.TimeBased), fileLogger.rollingFile)
							if err != nil {
								Errorf("create cron error %s", err.Error())
							}
						}

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
		go rollingFile()
	}

	return nil
}

func rollingFile() {
	for {
		if logger.writers != nil && rolling {
			for _, v := range logger.writers {
				v.rolling()
			}
		}
		if config.RollingInterval <= 0 {
			config.RollingInterval = 60
		}
		time.Sleep(time.Duration(config.RollingInterval) * time.Second)
	}
}

func StartRolling() {
	rolling = true
	crontab.Start()
}

func StopRolling() {
	rolling = false
	crontab.Stop()
}

func Tracef(format string, v ...interface{}) {
	for _, value := range logger.writers {
		value.Tracef(format, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	for _, value := range logger.writers {
		value.Debugf(format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	for _, value := range logger.writers {
		value.Infof(format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	for _, value := range logger.writers {
		value.Warnf(format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	for _, value := range logger.writers {
		value.Errorf(format, v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	for _, value := range logger.writers {
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
