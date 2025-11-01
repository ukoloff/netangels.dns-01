package na01_test

import (
	"na01"
	"testing"
)

func TestAuth(t *testing.T) {
	token, _ := na01.Auth()
	if len(token) == 0 {
		panic("Empty token")
	}
}

func TestZones(t *testing.T) {
	na01.Zones()
}
