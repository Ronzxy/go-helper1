package logger

import (
	"encoding/xml"
	"github.com/skygangsta/go-logger/config"
	"os"
)

func NewConfig(configFile string) (*config.Config, error) {
	file, err := os.OpenFile(configFile, os.O_RDONLY, 0)
	if err != nil {
		DefaultConsoleLogger().Errorf("error: Open config file %v", err)
		return nil, err
	}

	config := config.NewConfig()
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		DefaultConsoleLogger().Errorf("error: Decode xml %v", err)
		return nil, err
	}

	return config, nil
}
