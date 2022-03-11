package main

import (
	log "github.com/sirupsen/logrus"
	"tester/src/config"
	"tester/src/internal/metrics"
	"tester/src/internal/scheduler"
	"tester/src/internal/tester"
)

func main() {
	resultChan := make(chan tester.TestResult)
	defer close(resultChan)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:  false,
		DisableQuote:   true,
		FullTimestamp:  true,
		DisableSorting: false,
	})
	configuration := config.LoadConfiguration("config.yml")
	log.SetLevel(configuration.GetLogLevel())
	log.Info("Starting tester...")
	log.Debug(configuration)
	scheduler.StartJobs(configuration, resultChan)
	metrics.StartPrometheus(resultChan)
}
