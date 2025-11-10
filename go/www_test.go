package na01_test

import (
	"io"
	"na01"
	"net/http"
	"testing"
)

const URL = "http://localhost/"

func TestWWW(t *testing.T) {
	go func() {
		if err := na01.Start(); err != nil {
			panic(err)
		}
	}()
	defer na01.Stop()

	resp, err := http.Get(URL + "alive")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if len(resp.Header.Get("X-Health-Check")) == 0 {
		t.Fatal("Header not found!")
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "Ok" {
		t.Fatal("Health Check failed")
	}
}
