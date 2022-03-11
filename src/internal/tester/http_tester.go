package tester

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"tester/src/config"
	"time"
)

type HttpTestResult struct {
	Success        bool
	TestDuration   time.Duration
	Configuration  config.Configurator
	ResponseStatus int
}

func (r HttpTestResult) WasSuccessful() bool {
	return r.Success
}

func (r HttpTestResult) GetConfig() config.Configurator {
	return r.Configuration
}

type HttpTester struct {
	config         config.HttpTesterConfig
	resultsChannel chan TestResult
}

func (h HttpTester) validateEndpoint() error {
	_, err := url.ParseRequestURI(h.config.Endpoint)
	return err
}

func (h HttpTester) Validate() error {
	return h.validateEndpoint()
}

func (h HttpTester) Test() (TestResult, error) {
	client := &http.Client{}
	testResult := HttpTestResult{Configuration: h.config, Success: false}
	req, err := http.NewRequest(h.config.Method, h.config.Endpoint, nil)
	if err != nil {
		return testResult, err
	}
	for key, value := range h.config.Headers {
		req.Header.Add(key, value)
	}
	t1 := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	t2 := time.Now()
	testResult.TestDuration = t2.Sub(t1)
	if err != nil {
		log.Warn(err)
	} else {
		testResult.ResponseStatus = resp.StatusCode
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Debug(bodyString)
		if resp.StatusCode == h.config.SuccessStatus {
			testResult.Success = true
		}
	}
	log.Debug(testResult)
	h.resultsChannel <- testResult
	return testResult, nil
}
