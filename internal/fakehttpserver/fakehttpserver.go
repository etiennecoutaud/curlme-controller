package fakehttpserver

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// FakeHTTPServer fake http server struct
type FakeHTTPServer struct {
	server *http.Server
	body string
	statusCode int
	wg *sync.WaitGroup

}

// New instantiate new fakehttpserver
func New() *FakeHTTPServer {
	return &FakeHTTPServer{
		server: &http.Server{
			Addr: "localhost:8181",
		},
		wg: &sync.WaitGroup{},
	}
}

// Run start fake http server
func (f *FakeHTTPServer) Run() {
	f.wg.Add(1)
	go func(body string, sc int) {
		defer f.wg.Done()
		//f.server.Handler = http.HandlerFunc(f.handler)
		f.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(f.body))
			w.WriteHeader(f.statusCode)
		})
		f.server.ListenAndServe()
	}(f.body, f.statusCode)
	// Make sur Http server has time to start before handle request
	time.Sleep(5 * time.Second)
}

// Stop shutdown http server
func (f *FakeHTTPServer) Stop() {
	f.server.Shutdown(context.Background())
	f.wg.Wait()
}

// GetServerAddr http server address getter
func (f *FakeHTTPServer) GetServerAddr() string {
	return f.server.Addr
}

// SetBody http server body setter
func (f *FakeHTTPServer) SetBody(b string) {
	f.body = b
}

// SetStatusCode status code http setter
func (f *FakeHTTPServer) SetStatusCode(sc int) {
	f.statusCode = sc
}
