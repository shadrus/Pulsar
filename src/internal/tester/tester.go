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
	WasSuccessful() bool
	GetConfig() config.Configurator
}

func NewTester(configuration config.Configurator, resultsChannel chan TestResult) (Tester, error) {
	var tester Tester
	switch conf := configuration.(type) {
	case config.HttpTesterConfig:
		tester = HttpTester{config: conf, resultsChannel: resultsChannel}
	case config.CertificateTesterConfig:
		tester = CertificateTester{config: conf, resultsChannel: resultsChannel}
	default:
		return nil, errors.New("unknown tester type")
	}
	err := tester.Validate()
	return tester, err
}
