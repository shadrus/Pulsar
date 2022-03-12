package config

import (
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Configurator interface {
	GetInterval() int
}

type Config struct {
	LogLevel          string                    `yaml:"log_level"`
	HttpConfig        []HttpTesterConfig        `yaml:"http_config"`
	CertificateConfig []CertificateTesterConfig `yaml:"certificate_config"`
}

func (c Config) GetLogLevel() logrus.Level {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	return level
}

type HttpTesterConfig struct {
	Endpoint      string            `yaml:"endpoint"`
	Interval      int               `yaml:"interval"`
	Method        string            `yaml:"method"`
	SuccessStatus int               `yaml:"success_status"`
	Headers       map[string]string `yaml:"headers"`
}

func (t HttpTesterConfig) GetInterval() int {
	return t.Interval
}

type CertificateTesterConfig struct {
	Endpoint    string `yaml:"endpoint"`
	Interval    int    `yaml:"interval"`
	DaysForWarn int    `yaml:"days_for_warn"`
}

func (t CertificateTesterConfig) GetInterval() int {
	return t.Interval
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
