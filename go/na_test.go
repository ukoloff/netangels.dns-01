package na01_test

import (
	"fmt"
	"na01"
	"testing"
)

func TestAuth(t *testing.T) {
	token, err := na01.Auth()
	if err != nil {
		panic(err)
	}
	if len(token) == 0 {
		panic("Empty token")
	}
}

func TestZones(t *testing.T) {
	zz, err := na01.Zones()
	if err != nil {
		panic(err)
	}
	fmt.Println(zz)
}
