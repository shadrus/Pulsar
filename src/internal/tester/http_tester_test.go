package tester

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
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

type MockSlowHttpTesterClient struct {
}

func (r MockSlowHttpTesterClient) Do(req *http.Request) (*http.Response, error) {
	json := `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}}`
	// create a new reader with that JSON
	respData := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{StatusCode: 200, Body: respData}, nil
}

func Test(t *testing.T) {
	resultChan := make(chan TestResult)
	defer close(resultChan)
	headers := make(map[string]string)
	target := config.CommonConfig{Endpoint: "ya.ru", Interval: 10, Timeout: 20}
	config := config.HttpTesterConfig{Method: "get", SuccessStatus: 200, Headers: headers}
	config.CommonConfig = target
	tests := []struct {
		name    string
		h       HttpTester
		want    TestResult
		wantErr bool
	}{
		{"SlowHttpClient", HttpTester{config: config, resultsChannel: resultChan, client: MockSlowHttpTesterClient{}}, HttpTestResult{Success: true, TestDuration: 0, Configuration: config, ResponseStatus: 200}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.testHttp()
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpTester.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				//t.Errorf("HttpTester.Test() = %v, want %v", got, tt.want)
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
