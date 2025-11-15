package na01

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return strings.Trim(
		strings.ReplaceAll(
			base64.RawURLEncoding.EncodeToString(b),
			"_", "-"), "-")
}
