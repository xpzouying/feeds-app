package main

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type serverConfig struct {
	HTTPAddr string `yaml:"http_addr"`
}

type dbConfig struct {
	// UseMem 是否使用内存
	UseMem bool `yaml:"use_mem"`

	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type Tracer struct {
	Type    string `yaml:"type"`
	Address string `yaml:"address"`
}

type allConfig struct {
	Server serverConfig `yaml:"server"`

	DB dbConfig `yaml:"database"`

	Tracer Tracer `yaml:"tracer"`
}

func newConfig(r io.Reader) (*allConfig, error) {
	var cfg allConfig
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func newConfigFromFile(filename string) (*allConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return newConfig(file)
}
