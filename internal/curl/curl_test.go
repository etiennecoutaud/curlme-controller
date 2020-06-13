package curl_test

import (
	"github.com/etiennecoutaud/curlme-controller/internal/curl"
	"github.com/etiennecoutaud/curlme-controller/internal/fakehttpserver"
	"net/http"
	"testing"
)

func TestCurl_CallingURL(t *testing.T) {

	fakeHTTPServer := fakehttpserver.New()
	fakeHTTPServer.Run()


	tests := []struct{
		body string
		statusCode int
		isExpectedErrNil bool
	}{
		{
			body: "ok",
			statusCode: http.StatusOK,
			isExpectedErrNil: true,
		},
		{
			body: "toto",
			statusCode: http.StatusOK,
			isExpectedErrNil: true,
		},
		{
			body: "tutu",
			statusCode: http.StatusOK,
			isExpectedErrNil: true,
		},
		{
			body: "",
			statusCode: http.StatusInternalServerError,
			isExpectedErrNil: false,
		},
		{
			body: "",
			statusCode: http.StatusNotFound,
			isExpectedErrNil: false,
		},

	}


	for _, test := range tests {
		fakeHTTPServer.SetBody(test.body)
		fakeHTTPServer.SetStatusCode(test.statusCode)

		c := curl.New()
		resBody, err := c.CallingURL(fakeHTTPServer.GetServerAddr())
		if resBody != test.body || (err != nil && test.isExpectedErrNil) {
			t.Errorf("expect %s, got: %s, %v", test.body, resBody, err)
		}
	}
	fakeHTTPServer.Stop()
}
