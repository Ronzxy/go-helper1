go-logger
=========

Description
-----------

A log library for golang.

Installation
------------

This package can be installed with the go get command:

    go get github.com/skygangsta/go-logger
    
Logger Level
------------

    ALL < TRACE < DEBUG < INFO < WARN < ERROR < FATAL < OFF
    
LoggerWriter
------------

LoggerWriter is a complete logger that supports automatic scrolling of logger files by time or file size. It's initial by config file "[example/logger.xml](https://github.com/skygangsta/go-logger/blob/v0.1.0/example/logger.xml)":

```go
package main

import (
	"github.com/skygangsta/go-logger"
)

func main()  {
    err := logger.InitLogger("example/logger.xml")
    
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

ConsoleLogger
-------------

Console logger outputs log information to stdout:

```go
package main

import (
	"github.com/skygangsta/go-logger"
)

func main()  {
    consoleLogger := logger.NewConsoleLogger(logger.ALL)
    
    consoleLogger.Info("ConsoleLogger info message") 
}
```

FileLogger
----------

Console logger outputs log information to a file:

```go
package main

import (
	"github.com/skygangsta/go-logger"
)

func main()  {
    fileLogger, err := logger.NewFileLogger(logger.ALL, "logs/fileLogger.log", 0644)
	if err != nil {
		logger.DefaultConsoleLogger().Error(err.Error())
		return
	}
    
    fileLogger.Info("FileLogger info message")
}
```
