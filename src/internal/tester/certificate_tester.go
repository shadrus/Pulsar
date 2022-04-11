package tester

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"strings"
	"tester/src/config"
	"time"

	log "github.com/sirupsen/logrus"
)

type CertificateTestResult struct {
	Endpoint      string
	Configuration config.Configurator
	cert          *x509.Certificate
	CertNotAfter  float64
	CertNotBefore float64
}

func (r CertificateTestResult) GetConfig() config.Configurator {
	return r.Configuration
}

func (h CertificateTestResult) PrepareLabels() map[string]string {
	expectedLabels := map[string]string{"endpoint": h.Endpoint}
	if h.cert == nil {
		expectedLabels = map[string]string{
			"endpoint":  h.Endpoint,
			"serial_no": "",
			"issuer_cn": "",
			"o":         "",
			"cn":        "",
			"ou":        "",
		}
	} else {
		expectedLabels = map[string]string{
			"endpoint":  h.Endpoint,
			"serial_no": h.cert.SerialNumber.String(),
			"issuer_cn": h.cert.Issuer.CommonName,
			"o":         strings.Join(h.cert.Issuer.Organization, ","),
			"cn":        h.cert.Subject.CommonName,
			"ou":        strings.Join(h.cert.Subject.OrganizationalUnit, ","),
		}
	}

	return expectedLabels

}

type CertificateTester struct {
	config         config.CertificateTesterConfig
	resultsChannel chan TestResult
}

func NewCertificateTester(config config.CertificateTesterConfig, resultsChannel chan TestResult) *CertificateTester {
	return &CertificateTester{config: config, resultsChannel: resultsChannel}
}

func (h CertificateTester) validateEndpoint() error {
	u, err := url.Parse(h.config.GetEndpoint())
	if err != nil {
		return err
	}
	if u.Path == "" {
		return fmt.Errorf("wrong certificate url: %s. It mst be like domain.com", h.config.GetEndpoint())
	}
	return nil
}

func (h CertificateTester) Validate() error {
	return h.validateEndpoint()
}

func (h CertificateTester) testCert() (TestResult, error) {
	testResult := CertificateTestResult{Configuration: h.config, Endpoint: h.config.GetEndpoint()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.config.GetTimeout())*time.Second)
	defer cancel()
	d := tls.Dialer{
		Config: nil,
	}
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:443", h.config.GetEndpoint()))
	if err != nil {
		log.Warning(err)
		return testResult, err
	}
	defer conn.Close()
	tlsConn := conn.(*tls.Conn)
	err = tlsConn.VerifyHostname(h.config.GetEndpoint())
	if err != nil {
		log.Warning(err)
		return testResult, err
	}
	testResult.cert = tlsConn.ConnectionState().PeerCertificates[0]
	testResult.CertNotAfter = float64(tlsConn.ConnectionState().PeerCertificates[0].NotAfter.Unix())
	testResult.CertNotBefore = float64(tlsConn.ConnectionState().PeerCertificates[0].NotBefore.Unix())
	return testResult, nil
}

func (h CertificateTester) Test() (TestResult, error) {
	testResult, err := h.testCert()
	h.resultsChannel <- testResult
	return testResult, err
}
