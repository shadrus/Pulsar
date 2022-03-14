package main

import (
	"flag"
	"tester/src/config"
	"tester/src/internal/logger"
	"tester/src/internal/metrics"
	"tester/src/internal/scheduler"
	"tester/src/internal/tester"

	log "github.com/sirupsen/logrus"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yml", "config filename")
	flag.Parse()
}

func main() {
	resultChan := make(chan tester.TestResult)
	defer close(resultChan)
	configuration := config.LoadConfiguration(configFile)
	//all log settings must be set there
	logger.PrepareLogger(configuration.LogLevel)
	log.Debug(configuration)
	log.Info("Starting tester...")
	scheduler.StartJobs(configuration, resultChan)
	metrics.StartPrometheus(resultChan)
}
