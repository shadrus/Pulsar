package tester

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"tester/src/config"
	"time"

	log "github.com/sirupsen/logrus"
)

type CertificateTestResult struct {
	Success       bool
	DaysToExpire  float64
	Configuration config.Configurator
}

func (r CertificateTestResult) WasSuccessful() bool {
	return r.Success
}

func (r CertificateTestResult) GetConfig() config.Configurator {
	return r.Configuration
}

type CertificateTester struct {
	config         config.CertificateTesterConfig
	resultsChannel chan TestResult
}

func NewCertificateTester(config config.CertificateTesterConfig, resultsChannel chan TestResult) *CertificateTester {
	return &CertificateTester{config: config, resultsChannel: resultsChannel}
}

func (h CertificateTester) validateEndpoint() error {
	u, err := url.Parse(h.config.Endpoint)
	if err != nil {
		return err
	}
	if u.Path == "" {
		return fmt.Errorf("wrong certificate url: %s. It mst be like domain.com", h.config.Endpoint)
	}
	return nil
}

func (h CertificateTester) Validate() error {
	return h.validateEndpoint()
}

func (h CertificateTester) Test() (TestResult, error) {
	testResult := CertificateTestResult{Configuration: h.config, Success: false}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	d := tls.Dialer{
		Config: nil,
	}
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:443", h.config.Endpoint))
	if err != nil {
		log.Warning(err)
		h.resultsChannel <- testResult
		return testResult, err
	}
	defer conn.Close()
	tlsConn := conn.(*tls.Conn)
	err = tlsConn.VerifyHostname(h.config.Endpoint)
	if err != nil {
		log.Warning(err)
		h.resultsChannel <- testResult
		return testResult, err
	}
	expiry := tlsConn.ConnectionState().PeerCertificates[0].NotAfter
	timeDiff := time.Until(expiry)
	daysToExpire := timeDiff.Hours() / 24
	testResult.DaysToExpire = daysToExpire
	if int(daysToExpire) > h.config.DaysForWarn {
		testResult.Success = true
	}
	log.Debug(expiry)
	h.resultsChannel <- testResult
	return testResult, nil
}
