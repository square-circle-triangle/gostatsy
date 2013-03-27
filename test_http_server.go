package statsy

import (
	"io/ioutil"
	"log"
	"net/http"
)

const (
	HttpPort = "30456"
)

type TestHttpServer struct {
	Body           string
	ContentType    string
	MockBody       string
	MockStatusCode int
}

func NewHTTPServer() *TestHttpServer {
	server := &TestHttpServer{}
	server.ResetMocks()

	http.HandleFunc("/api/v1/multiple_event", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Panicln("Error reading body", err)
		}

		server.Body = string(bodyBytes)
		server.ContentType = r.Header.Get("Content-Type")

		w.WriteHeader(server.MockStatusCode)
		_, err = w.Write([]byte(server.MockBody))

		// reset for next test
		server.ResetMocks()
	})

	go func() {
		err := http.ListenAndServe(":"+HttpPort, nil)
		if err != nil {
			log.Panic("Error starting test http server", err)
		}
	}()
	return server
}

func (t *TestHttpServer) ResetMocks() {
	t.MockBody = ""
	t.MockStatusCode = 200
}

func (t *TestHttpServer) BaseURL() string {
	return "http://localhost:" + HttpPort + "/"
}
