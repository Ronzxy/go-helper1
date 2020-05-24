# go-logger

[![Go Report Card](https://goreportcard.com/badge/github.com/ronzxy/go-logger)](https://goreportcard.com/report/github.com/ronzxy/go-logger)
[![GoDoc](https://godoc.org/github.com/ronzxy/go-logger?status.svg)](https://godoc.org/github.com/ronzxy/go-logger)
[![Github All Releases](https://img.shields.io/github/downloads/ronzxy/go-logger/total.svg)](https://github.com/ronzxy/go-logger/releases)
[![GitHub release](https://img.shields.io/github/release/ronzxy/go-logger/all.svg)](https://github.com/ronzxy/go-logger/releases)
[![GitHub Release Date](https://img.shields.io/github/release-date-pre/ronzxy/go-logger.svg)](https://github.com/ronzxy/go-logger/releases)
[![GitHub license](https://img.shields.io/github/license/ronzxy/go-logger.svg)](https://github.com/ronzxy/go-logger/blob/master/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/ronzxy/go-logger.svg)](https://github.com/ronzxy/go-logger/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/ronzxy/go-logger.svg)](https://github.com/ronzxy/go-logger/network)
[![Sourcegraph](https://sourcegraph.com/github.com/ronzxy/go-logger/-/badge.svg)](https://sourcegraph.com/github.com/ronzxy/go-logger?badge)

## Description

A log library for golang. Can be initialized from xml format configuration file, supports scrolling based on file size and time of log files.

## Installation

This package can be installed with the go get command:

```sh
    go get github.com/ronzxy/go-logger
```

### Logger Level

```text
    ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
```

### Initialized

LoggerWriter is a complete logger that supports automatic scrolling of logger files by time or file size. It's initial by config file with xml format "[example/logger.xml](https://github.com/ronzxy/go-logger/blob/master/example/logger.xml)":

```go
package main

import (
    "github.com/ronzxy/go-logger"
)

func main()  {
    err := logger.Init("example/logger.xml")

    if err != nil {
        logger.DefaultConsoleLogger().Error(err.Error())
        return
    }

    logger.Trace("Test Logger trace message")
    logger.Debug("Test Logger debug message")
    logger.Info("Test Logger info message")
    logger.Warn("Test Logger warn message")
    logger.Error("Test Logger error message")
}
```

### ConsoleLogger

Console logger outputs log information to stdout:

```go
package main

import (
    "github.com/ronzxy/go-logger"
)

func main()  {
    consoleLogger := logger.NewConsoleLogger(logger.ALL)

    consoleLogger.Info("ConsoleLogger info message") 
}
```

### FileLogger

File logger outputs log information to a file:

```go
package main

import (
    "github.com/ronzxy/go-logger"
)

func main()  {
    fileLogger, err := logger.NewFileLogger(logger.ALL, "logs/fileLogger.log")
    if err != nil {
        logger.DefaultConsoleLogger().Error(err.Error())
        return
    }

    fileLogger.Info("FileLogger info message")
}
```
