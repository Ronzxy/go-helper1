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
	"encoding/xml"
	"os"
)

type Config struct {
	XMLName         xml.Name   `xml:"Configuration"`
	RollingInterval int        `xml:"rollingInterval,attr"`
	Properties      []Property `xml:"Properties>Property"`
	Loggers         []Logger   `xml:"Loggers>Logger"`
	DefaultFilter   Filter     `xml:"Filters>DefaultFilter>Filter"`
	PackageFilters  []Filter   `xml:"Filters>PackageFilter>Filter"`
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
	Format      Format   `xml:"Format"`
	Level       Level    `xml:"Level"`
	Rolling     Rolling  `xml:"Rolling"`
}

type Format struct {
	XMLName xml.Name `xml:"Format"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:",innerxml"`
}

type Level struct {
	XMLName xml.Name `xml:"Level"`
	Allow   string   `xml:"Allow"`
	Deny    string   `xml:"Deny"`
}

type Rolling struct {
	XMLName   xml.Name `xml:"Rolling"`
	TimeBased string   `xml:"TimeBased"`
	SizeBased int      `xml:"SizeBased"`
	KeepCount int      `xml:"KeepCount"`
}

type Filter struct {
	XMLName xml.Name `xml:"Filter"`
	Name    string   `xml:"name,attr"`
	Loggers []string `xml:"Logger"`
}

func NewConfig(configFile string) (*Config, error) {
	file, err := os.OpenFile(configFile, os.O_RDONLY, 0)
	if err != nil {
		DefaultConsoleLogger().Errorf("error: Open config file %v", err)
		return nil, err
	}

	config := &Config{}
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		DefaultConsoleLogger().Errorf("error: Decode xml %v", err)
		return nil, err
	}

	return config, nil
}
