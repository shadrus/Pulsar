package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"tester/src/config"
	"tester/src/internal/metrics"
	"tester/src/internal/scheduler"
	"tester/src/internal/tester"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yml", "config filename")
	flag.Parse()
}

func main() {
	resultChan := make(chan tester.TestResult)
	defer close(resultChan)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:  false,
		DisableQuote:   true,
		FullTimestamp:  true,
		DisableSorting: false,
	})
	configuration := config.LoadConfiguration(configFile)
	log.SetLevel(configuration.GetLogLevel())
	log.Info("Starting tester...")
	log.Debug(configuration)
	scheduler.StartJobs(configuration, resultChan)
	metrics.StartPrometheus(resultChan)
}
