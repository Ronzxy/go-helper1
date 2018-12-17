package logger

type Formatter interface {
	Message(data map[string]interface{}, args ...interface{}) string
}
