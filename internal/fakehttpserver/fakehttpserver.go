package fakehttpserver

import (
	"net/http"
	"net/http/httptest"
	//"sync"
	//"time"
)

// FakeHTTPServer fake http server struct
type FakeHTTPServer struct {
	//server *http.Server
	//body string
	//statusCode int
	//wg *sync.WaitGroup
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
	//	server: &http.Server{
	//		Addr: "localhost:8181",
	//	},
	//	wg: &sync.WaitGroup{},
	//}
}

//func (f *FakeHTTPServer) handler(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte(f.body))
//	w.WriteHeader(f.statusCode)
//}

//func (f *FakeHTTPServer)

// Run start fake http server
//func (f *FakeHTTPServer) Run() {
//	f.wg.Add(1)
//	go func(body string, sc int) {
//		defer f.wg.Done()
//		//f.server.Handler = http.HandlerFunc(f.handler)
//		f.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.Write([]byte(f.body))
//			w.WriteHeader(f.statusCode)
//		})
//		f.server.ListenAndServe()
//	}(f.body, f.statusCode)
//	// Make sur Http server has time to start before handle request
//	time.Sleep(1 * time.Second)
//}

// Stop shutdown http server
func (f *FakeHTTPServer) Stop() {
	f.server.Close()
}

// GetServerAddr http server address getter
func (f *FakeHTTPServer) GetServerAddr() string {
	return f.server.URL
}

// SetBody http server body setter
//func (f *FakeHTTPServer) SetBody(b string) {
//	f.body = b
//}
//
//// SetStatusCode status code http setter
//func (f *FakeHTTPServer) SetStatusCode(sc int) {
//	f.statusCode = sc
//}

// GetServer return fake http server
func (f *FakeHTTPServer) GetServer() *httptest.Server{
	return f.server
}
