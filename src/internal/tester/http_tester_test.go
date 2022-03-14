package tester

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"tester/src/config"
	"testing"
)

func TestHttpTestResult_WasSuccessful(t *testing.T) {
	tests := []struct {
		name string
		r    HttpTestResult
		want bool
	}{
		{"Ok", HttpTestResult{Success: true}, true},
		{"Error", HttpTestResult{Success: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.WasSuccessful(); got != tt.want {
				t.Errorf("HttpTestResult.WasSuccessful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpTester_validateEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		h        HttpTester
		endpoint string
		wantErr  bool
	}{
		{"https", HttpTester{}, "https://test.com", false},
		{"http", HttpTester{}, "http://test.com", false},
		{"no http", HttpTester{}, "test.com", true},
		{"no domain", HttpTester{}, "test", true},
		{"ip", HttpTester{}, "127.0.0.1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.validateEndpoint(tt.endpoint); (err != nil) != tt.wantErr {
				t.Errorf("HttpTester.validateEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MockHttpTesterClient struct{}

func (r MockHttpTesterClient) Do(req *http.Request) (*http.Response, error) {
	json := `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}}`
	// create a new reader with that JSON
	respData := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{StatusCode: 200, Body: respData}, nil
}

type BadStatusMockHttpTesterClient struct{}

func (r BadStatusMockHttpTesterClient) Do(req *http.Request) (*http.Response, error) {
	json := `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}}`
	// create a new reader with that JSON
	respData := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{StatusCode: 203, Body: respData}, nil
}

type ErrorMockHttpTesterClient struct{}

func (r ErrorMockHttpTesterClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, errors.New("Timeout")
}

func Test(t *testing.T) {
	resultChan := make(chan TestResult)
	defer close(resultChan)
	headers := make(map[string]string)
	target := config.CommonConfig{Endpoint: "ya.ru", Interval: 10, Timeout: 5}
	config := config.HttpTesterConfig{Method: "get", SuccessStatus: 200, Headers: headers}
	config.CommonConfig = target
	tests := []struct {
		name    string
		h       HttpTester
		want    TestResult
		wantErr bool
	}{
		{"HttpClient good request", HttpTester{config: config, resultsChannel: resultChan, client: MockHttpTesterClient{}}, HttpTestResult{Success: true, TestDuration: 0, Configuration: config, ResponseStatus: 200}, false},
		{"HttpClient wrong response code", HttpTester{config: config, resultsChannel: resultChan, client: BadStatusMockHttpTesterClient{}}, HttpTestResult{Success: false, TestDuration: 0, Configuration: config, ResponseStatus: 203}, false},
		{"HttpClient error on Do request", HttpTester{config: config, resultsChannel: resultChan, client: ErrorMockHttpTesterClient{}}, HttpTestResult{Success: false, Configuration: config}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.h.testHttp()
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpTester.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result.WasSuccessful() != tt.want.WasSuccessful() {
				t.Errorf("HttpTester.Test() Success = %v, want %v", result.WasSuccessful(), tt.want.WasSuccessful())
				return
			}

		})
	}
}

func Test_testResponseBody(t *testing.T) {
	goodBody := "<link rel=\"search\" href=\"//yandex.ru/opensearch.xml\" title=\"Яндекс\" type=\"application/opensearchdescription+xml\">"
	wrongBody := "<link rel=\"search\" title=\"Яндекс\" type=\"application/opensearchdescription+xml\">"
	type args struct {
		body        string
		checkString string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"No check test", args{body: goodBody, checkString: ""}, true},
		{"Valid check test", args{body: goodBody, checkString: "yandex"}, true},
		{"Wrong check test", args{body: wrongBody, checkString: "yandex"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testResponseBody(tt.args.body, tt.args.checkString); got != tt.want {
				t.Errorf("testResponseBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
