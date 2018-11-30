package logger

type Writer interface {
	Tracef(format string, v ...interface{})

	Debugf(format string, v ...interface{})

	Infof(format string, v ...interface{})

	Warnf(format string, v ...interface{})

	Errorf(format string, v ...interface{})

	Fatalf(format string, v ...interface{})

	Trace(message string)

	Debug(message string)

	Info(message string)

	Warn(message string)

	Error(message string)

	Fatal(message string)
}
