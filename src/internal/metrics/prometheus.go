// package metrics converts different tests results to Prometheus metrics
package metrics

import (
	"net/http"
	"strconv"
	"tester/src/config"
	"tester/src/internal/tester"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	// httpTest is Summary type metric.
	// It provides info from http tests.
	httpTest = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "pulsar_http_request_seconds",
			Help: "Time to get response from the endpoint in seconds",
		},
		[]string{"endpoint", "success", "status"},
	)
	// httpTest is Gauge type metric.
	// It provides info from certificate validity test.
	certTest = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pulsar_days_to_expire_cert",
			Help: "Day left to expire certificate",
		},
		[]string{"endpoint", "success"},
	)
)

func serveHttpTestResults(conf config.HttpTesterConfig, result tester.TestResult) {
	httpResult := result.(tester.HttpTestResult)
	httpTest.WithLabelValues(conf.Endpoint, strconv.FormatBool(result.WasSuccessful()), strconv.Itoa(httpResult.ResponseStatus)).Observe(httpResult.TestDuration.Seconds())
}

func serveCertificateTestResults(conf config.CertificateTesterConfig, result tester.TestResult) {
	certificateTestResult := result.(tester.CertificateTestResult)
	certTest.WithLabelValues(conf.Endpoint, strconv.FormatBool(result.WasSuccessful())).Set(certificateTestResult.DaysToExpire)
}

// serveTestResults decides what metric must be updates based on test configuration
func serveTestResults(resultsChannel <-chan tester.TestResult) {
	for result := range resultsChannel {
		log.Debugf("Got results %v", result)
		switch conf := result.GetConfig().(type) {
		case config.HttpTesterConfig:
			serveHttpTestResults(conf, result)
		case config.CertificateTesterConfig:
			serveCertificateTestResults(conf, result)
		default:
			log.Error("unknown tester type")
		}

	}
}

// StartPrometheus creates job for gather test results and exposes metrics endpoint
func StartPrometheus(resultsChannel <-chan tester.TestResult) {
	prometheus.MustRegister(httpTest)
	prometheus.MustRegister(certTest)
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// TODO must be based on config
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	go serveTestResults(resultsChannel)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
