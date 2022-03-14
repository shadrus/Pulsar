package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Configurator interface {
	GetInterval() int
	GetTimeout() int
	GetEndpoint() string
}

type Config struct {
	LogLevel          string                    `yaml:"log_level"`
	HttpConfig        []HttpTesterConfig        `yaml:"http_config"`
	CertificateConfig []CertificateTesterConfig `yaml:"certificate_config"`
}

type CommonConfig struct {
	Endpoint string `yaml:"endpoint"`
	Interval int    `yaml:"interval"`
	Timeout  int    `yaml:"timeout"`
}

func (t CommonConfig) GetInterval() int {
	return t.Interval
}

func (t CommonConfig) GetTimeout() int {
	return t.Timeout
}

func (t CommonConfig) GetEndpoint() string {
	return t.Endpoint
}

type HttpTesterConfig struct {
	CommonConfig  `yaml:"target"`
	Method        string            `yaml:"method"`
	SuccessStatus int               `yaml:"success_status"`
	Headers       map[string]string `yaml:"headers"`
	CheckText     string            `yaml:"check_text,omitempty"`
}

type CertificateTesterConfig struct {
	CommonConfig `yaml:"target"`
	DaysForWarn  int `yaml:"days_for_warn"`
}

func LoadConfiguration(filePath string) *Config {
	var config Config
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &config
}
