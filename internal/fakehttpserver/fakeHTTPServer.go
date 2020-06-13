package fakehttpserver

import (
	"context"
	"net/http"
	"sync"
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

func (f *FakeHTTPServer) handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(f.body))
	w.WriteHeader(f.statusCode)
}

// Run start fake http server
func (f *FakeHTTPServer) Run() {
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		f.server.Handler = http.HandlerFunc(f.handler)
		f.server.ListenAndServe()
	}()
	//time.Sleep(1 * time.Second)
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
