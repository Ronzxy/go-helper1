package logger

import (
	"encoding/json"
	"time"
)

type JSONFormatter struct {
	Indent bool
}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		Indent: false,
	}
}

func (this *JSONFormatter) Message(data map[string]interface{}, args ...interface{}) string {
	var (
		buf []byte
		err error
	)

	data["Time"] = time.Now().Format(DefaultLogTimeFormat)
	data["Message"] = args

	if this.Indent {
		buf, err = json.MarshalIndent(data, "", "")
	} else {
		buf, err = json.Marshal(data)
	}

	if err != nil {
		Errorf("")
	}

	return string(buf)
}
