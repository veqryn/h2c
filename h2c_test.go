package h2c

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"testing"

	"golang.org/x/net/http2"
)

func TestHandlerH2CServeHTTP(t *testing.T) {
	t.Parallel()

	startServer()

	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
	}

	resp, err := client.Get("http://127.0.0.1:8080/")
	if err != nil {
		t.Fatal(err)
	}

	if resp.ProtoMajor != 2 {
		t.Errorf("Expected protocol version: %d; Got: %d", 2, resp.ProtoMajor)
	}
}

func startServer() {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	h2cWrapper := &HandlerH2C{
		Handler:  router,
		H2Server: &http2.Server{},
	}

	srv := http.Server{
		Addr:    ":8080",
		Handler: h2cWrapper,
	}

	go func() {
		srv.ListenAndServe()
	}()

	// Wait for it to start
	for {
		resp, _ := http.Get("http://127.0.0.1:8080/")
		if resp != nil && resp.StatusCode == http.StatusOK {
			break
		}
	}
}
