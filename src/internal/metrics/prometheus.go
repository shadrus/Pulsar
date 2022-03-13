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
	// Create a summary to track fictional interservice RPC latencies for three
	// distinct services with different latency distributions. These services are
	// differentiated via a "service" label.
	httpTest = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "pulsar_http_request_seconds",
			Help: "Time to get response from the endpoint in seconds",
		},
		[]string{"endpoint", "success", "status"},
	)
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

func StartPrometheus(resultsChannel <-chan tester.TestResult) {
	prometheus.MustRegister(httpTest)
	prometheus.MustRegister(certTest)
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	go serveTestResults(resultsChannel)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
