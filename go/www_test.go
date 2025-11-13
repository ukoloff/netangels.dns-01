package na01_test

import (
	"na01"
	"testing"
	"time"
)

func TestWWW(t *testing.T) {
	go func() {
		if err := na01.Start(); err != nil {
			panic(err)
		}
	}()
	defer na01.Stop()

	time.Sleep(108 * time.Millisecond)
	err := na01.FireAlive()
	if err != nil {
		t.Fatal(err)
	}
}
