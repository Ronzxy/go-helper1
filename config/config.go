package config

import (
	"encoding/xml"
)

type Config struct {
	XMLName    xml.Name   `xml:"Configuration"`
	Properties []Property `xml:"Properties>Property"`
	Loggers    []Logger   `xml:"Loggers>Logger"`
}

type Property struct {
	XMLName xml.Name `xml:"Property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",innerxml"`
}

type Logger struct {
	XMLName     xml.Name `xml:"Logger"`
	Name        string   `xml:"name,attr"`
	Target      string   `xml:"target,attr"`
	FileName    string   `xml:"fileName,attr"`
	FilePattern string   `xml:"filePattern,attr"`
	Level       Level    `xml:"Level"`
	Rolling     Rolling  `xml:"Rolling"`
}

type Level struct {
	XMLName xml.Name `xml:"Level"`
	Allow   string   `xml:"Allow"`
	Deny    string   `xml:"Deny"`
}

type Rolling struct {
	XMLName   xml.Name `xml:"Rolling"`
	TimeBased int      `xml:"TimeBased"`
	SizeBased string   `xml:"SizeBased"`
	KeepCount int      `xml:"KeepCount"`
}

func NewConfig() *Config {
	return &Config{}
}
