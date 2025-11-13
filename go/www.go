package na01

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var server = http.Server{}

type acmeReq struct {
	FQDN  string `json:"fqdn"`
	Value string `json:"value"`
}

func Start() error {
	http.HandleFunc("/alive", alive)
	http.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Bye")
		go func() {
			time.Sleep(300 * time.Millisecond)
			Stop()
		}()
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<li><a
      href="https://github.com/ukoloff/netangels.dns-01.git">Source</a>
      <li><a href="/alive">Health Check<a>
      <li><a href="/quit">Quit</a>`)
	})
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func Stop() error {
	return server.Shutdown(context.Background())
}

func alive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Health-Check", "+")
	fmt.Fprint(w, "Ok")
}
