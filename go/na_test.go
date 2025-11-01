package na01_test

import (
	"fmt"
	"na01"
	"testing"
)

func TestAuth(t *testing.T) {
	token, _ := na01.Auth()
	fmt.Println(token)
}
