package fakehttpserver

import (
	"net/http"
	"net/http/httptest"
)

// FakeHTTPServer fake http server struct
type FakeHTTPServer struct {
	server *httptest.Server
}

// New instantiate new fakehttpserver
func New(body string, statusCode int) *FakeHTTPServer {
	return &FakeHTTPServer{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(body))
			w.WriteHeader(statusCode)
		})),
	}
}

// GetServerAddr http server address getter
func (f *FakeHTTPServer) GetServerAddr() string {
	return f.server.URL
}

// GetServer return fake http server
func (f *FakeHTTPServer) GetServer() *httptest.Server {
	return f.server
}
