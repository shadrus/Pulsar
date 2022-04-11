// package metrics converts different tests results to Prometheus metrics
package metrics

import (
	"net/http"
	"tester/src/config"
	"tester/src/internal/tester"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	certTestNotAfter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pulsar_cert_not_after",
			Help: "The date after which a peer certificate expires. Expressed as a Unix Epoch Time.",
		},
		[]string{"endpoint", "serial_no", "issuer_cn", "cn", "o", "ou"},
	)

	certTestNotBefore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pulsar_cert_not_before",
			Help: "The date before which a peer certificate is not valid. Expressed as a Unix Epoch Time.",
		},
		[]string{"endpoint", "serial_no", "issuer_cn", "cn", "o", "ou"},
	)
)

func mapToArray(labels map[string]string) []string {
	v := make([]string, 0, len(labels))
	for _, value := range labels {
		v = append(v, value)
	}
	return v
}

func serveHttpTestResults(conf config.HttpTesterConfig, result tester.HttpTestResult) {
	httpTest.WithLabelValues(mapToArray(result.PrepareLabels())...).Observe(result.TestDuration.Seconds())
}

func serveCertificateTestResults(conf config.CertificateTesterConfig, result tester.CertificateTestResult) {
	labels := result.PrepareLabels()
	certTestNotAfter.WithLabelValues(labels["endpoint"], labels["serial_no"], labels["issuer_cn"], labels["cn"], labels["o"], labels["ou"]).Set(result.CertNotAfter)
	certTestNotBefore.WithLabelValues(labels["endpoint"], labels["serial_no"], labels["issuer_cn"], labels["cn"], labels["o"], labels["ou"]).Set(result.CertNotBefore)
}

// serveTestResults decides what metric must be updates based on test configuration
func serveTestResults(resultsChannel <-chan tester.TestResult) {
	for result := range resultsChannel {
		log.Debugf("Got results %v", result)
		switch conf := result.GetConfig().(type) {
		case config.HttpTesterConfig:
			serveHttpTestResults(conf, result.(tester.HttpTestResult))
		case config.CertificateTesterConfig:
			serveCertificateTestResults(conf, result.(tester.CertificateTestResult))
		default:
			log.Error("unknown tester type")
		}

	}
}

// StartPrometheus creates job for gather test results and exposes metrics endpoint
func StartPrometheus(resultsChannel <-chan tester.TestResult) {
	prometheus.MustRegister(httpTest)
	prometheus.MustRegister(certTestNotAfter)
	prometheus.MustRegister(certTestNotBefore)
	prometheus.Unregister(collectors.NewGoCollector())
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
