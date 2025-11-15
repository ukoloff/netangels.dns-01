package na01

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const URL = "http://localhost/"

var server = http.Server{}

type acmeReq struct {
	FQDN  string `json:"fqdn"`
	Value string `json:"value"`
}

func StartWWW() error {
	http.HandleFunc("/alive", alive)
	http.HandleFunc("/quit", quit)
	http.HandleFunc("/present", present)
	http.HandleFunc("/cleanup", cleanup)
	http.HandleFunc("/", home)
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func StopWWW() error {
	return server.Shutdown(context.Background())
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<li><a
      href="https://github.com/ukoloff/netangels.dns-01.git">Source</a>
      <li><a href="/alive">Health Check<a>
      <li><a href="/quit">Quit</a>`)
}

func quit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Bye")
	go func() {
		time.Sleep(300 * time.Millisecond)
		StopWWW()
	}()
}

func alive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Health-Check", "+")
	fmt.Fprint(w, "Ok")
}

func FireAlive() error {
	resp, err := http.Get(URL + "alive")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if len(resp.Header.Get("X-Health-Check")) == 0 {
		return errors.New("header not found")
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(b) != "Ok" {
		return errors.New("health check failed")
	}
	return nil
}

func sendError(w http.ResponseWriter, err error) {
	type croak struct {
		Message string `json:"message"`
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(croak{Message: err.Error()})
}

func verb(w http.ResponseWriter, r *http.Request, handler func(in acmeReq) (any, error)) {
	var in acmeReq
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		sendError(w, err)
		return
	}
	res, err := handler(in)
	if err != nil {
		sendError(w, err)
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		sendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func present(w http.ResponseWriter, r *http.Request) {
	verb(w, r, func(in acmeReq) (any, error) {
		return Present(in.FQDN, in.Value)
	})
}

func cleanup(w http.ResponseWriter, r *http.Request) {
	verb(w, r, func(in acmeReq) (any, error) {
		return CleanUp(in.FQDN, in.Value)
	})
}

func fire(cmd, fqdn, text string, out any) error {
	in := acmeReq{FQDN: fqdn, Value: text}
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(data)
	resp, err := http.Post(URL+cmd, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New("HTTP Error: " + strconv.Itoa(resp.StatusCode))
	}
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func FirePresent(fqdn, text string) (RR, error) {
	var r RR
	return r, fire("present", fqdn, text, &r)
}

func FireCleanUp(fqdn, text string) ([]RR, error) {
	var rs []RR
	return rs, fire("cleanup", fqdn, text, &rs)
}

func FireQuit() {
	resp, err := http.Get(URL + "quit")
	if err == nil {
		resp.Body.Close()
	}
}
