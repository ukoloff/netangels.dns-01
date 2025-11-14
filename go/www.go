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

func Start() error {
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

func Stop() error {
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
		Stop()
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

func present(w http.ResponseWriter, r *http.Request) {
	var in acmeReq
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		// return err
	}
	rr, err := Present(in.FQDN, in.Value)
	if err != nil {
		// return err
	}
	data, err := json.Marshal(rr)
	if err != nil {
		// return err
	}
	w.Write(data)
}

func cleanup(w http.ResponseWriter, r *http.Request) {
	var in acmeReq
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		// return err
	}
	rrs, err := CleanUp(in.FQDN, in.Value)
	if err != nil {
		// return err
	}
	data, err := json.Marshal(rrs)
	if err != nil {
		// return err
	}
	w.Write(data)
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

func FirePresent(fqdn, text string, out any) error {
	return fire("present", fqdn, text, out)
}

func FireCleanUp(fqdn, text string, out any) error {
	return fire("cleanup", fqdn, text, out)
}
