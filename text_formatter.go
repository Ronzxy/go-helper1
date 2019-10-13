package logger

import (
	"fmt"
	"github.com/skygangsta/go-helper"
	"strings"
	"time"
)

type TextFormatter struct {
	Format string
}

var (
	DefaultFormat = "%{Prefix} - %{Time:yyyy-mm-dd HH:MM:SS.ms} - %{Level:5} - %{File}:%{Line:3} - %{Message}"
)

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		Format: DefaultFormat,
	}
}

func NewTextFormatterWithFormat(format string) *TextFormatter {
	this := NewTextFormatter()
	this.Format = VariableReplaceByConfig(format)

	return this
}

func (this *TextFormatter) SetFormat(format string) {
	this.Format = format
}

func (this *TextFormatter) Message(data map[string]interface{}, args ...interface{}) string {
	var (
		varPattern string
		varName    string
		vars       = make([]string, 2)
		str        = this.Format
	)
	for {
		varPattern, varName = Variable("%", "([a-zA-Z_][0-9a-zA-Z\\s\\._/:-]*)", str)
		if varName == "" {
			// no variable, exit loop
			break
		}

		if strings.Contains(varName, ":") {
			vars = strings.Split(varName, ":")
		} else {
			vars[0] = varName
			vars[1] = ""
		}

		switch strings.ToUpper(vars[0]) {
		case "PREFIX":
			{
				varName = fmt.Sprintf("%v", data["Prefix"])
			}
		case "TIME":
			{
				if vars[1] != "" {
					varName = helper.Time.Format(strings.Join(vars[1:], ":"), time.Now())
				} else {
					varName = time.Now().Format(DefaultLogTimeFormat)
				}
			}
		case "LEVEL":
			{
				format := "%s"
				if vars[1] != "" {
					format = fmt.Sprintf("%%-%ss", vars[1])
				}

				varName = fmt.Sprintf(format, data["Level"])
			}
		case "FILE":
			{
				format := "%s"
				if vars[1] != "" {
					format = fmt.Sprintf("%%-%ss", vars[1])
				}

				varName = fmt.Sprintf(format, data["File"])
			}
		case "LINE":
			{
				format := "%s"
				if vars[1] != "" {
					format = fmt.Sprintf("%%-%sd", vars[1])
				}

				varName = fmt.Sprintf(format, data["Line"])
			}
		case "MESSAGE":
			{
				varName = fmt.Sprint(args...)
			}
		default:
			{
				DefaultConsoleLogger().Errorf("unsupported log format %s", varName)
			}
		}

		str = strings.Replace(str, varPattern, varName, -1)
	}

	return str
}
