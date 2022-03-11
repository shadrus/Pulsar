package scheduler

import (
	"github.com/go-co-op/gocron"
	"tester/src/config"
	"tester/src/internal/tester"
	"time"
)

func startJob(configuration config.Configurator, resultsChannel chan tester.TestResult, scheduler *gocron.Scheduler)  {
	test, err := tester.NewTester(configuration, resultsChannel)
	if err != nil {
		panic(err.Error())
	}
	scheduler.Every(configuration.GetInterval()).Seconds().Do(test.Test)
}

func StartJobs(config *config.Config, resultsChannel chan tester.TestResult) {
	sched := gocron.NewScheduler(time.UTC)
	for _, conf := range config.HttpConfig {
		startJob(conf, resultsChannel, sched)
	}
	for _, conf := range config.CertificateConfig {
		startJob(conf, resultsChannel, sched)
	}
	sched.StartAsync()
}
