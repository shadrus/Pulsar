package tester

import (
	"errors"
	"tester/src/config"
)

type Type int

type Tester interface {
	Validate() error
	Test() (TestResult, error)
}

type TestResult interface {
	GetConfig() config.Configurator
	PrepareLabels() map[string]string
}

func NewTester(configuration config.Configurator, resultsChannel chan TestResult) (Tester, error) {
	var tester Tester
	switch conf := configuration.(type) {
	case config.HttpTesterConfig:
		tester = NewHttpTester(conf, resultsChannel)
	case config.CertificateTesterConfig:
		tester = NewCertificateTester(conf, resultsChannel)
	default:
		return nil, errors.New("unknown tester type")
	}
	err := tester.Validate()
	return tester, err
}
