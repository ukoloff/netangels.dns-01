package na01_test

import (
	"crypto/x509"
	"encoding/pem"
	"na01"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLego(t *testing.T) {
	me := exePath()
	cmd := exec.Command("go", "build", "-C", "main", "-ldflags", "-s -w", "-o", me)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(me)

	t.Run("cli", func(t *testing.T) {
		lego(t, "exec", map[string]string{
			"EXEC_PATH":                me,
			"EXEC_POLLING_INTERVAL":    "10",
			"EXEC_PROPAGATION_TIMEOUT": "300",
		})
	})

	t.Run("www", func(t *testing.T) {
		cmd := exec.Command(me, "www")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			resp, err := http.Get("http://localhost/quit")
			if err != nil {
				return
			}
			resp.Body.Close()
		}()
	})
}

func exePath() string {
	me, _ := os.Executable()
	f, _ := os.CreateTemp("", "na01-*"+filepath.Ext(me))
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

func lego(t *testing.T, provider string, env map[string]string) error {
	domain := strings.ToLower(na01.RandomString(7) + "." + provider + ".uralhimmash.com")

	cmd := exec.Command("lego", "-a",
		"-dns", provider,
		"-d", domain,
		"--pfx",
		"--dns.resolvers", "8.8.8.8",
		"run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Environ()
	env["LEGO_EMAIL"] = "stas@ekb.ru"
	env["LEGO_SERVER"] = "https://acme-staging-v02.api.letsencrypt.org/directory"
	for k, v := range env {
		e = append(e, k+"="+v)
	}
	cmd.Env = e
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	for _, ext := range []string{"key", "pfx", "crt", "json"} {
		_, err := os.Stat("./.lego/certificates/" + domain + "." + ext)
		if err != nil {
			t.Fatal(err)
		}
	}

	blob, err := os.ReadFile("./.lego/certificates/" + domain + ".crt")
	if err != nil {
		t.Fatal(err)
	}
	der, _ := pem.Decode(blob)
	crt, err := x509.ParseCertificate(der.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if crt.Subject.CommonName != domain {
		t.Fatal("Invalid CN")
	}
	if crt.DNSNames[0] != domain {
		t.Fatal("Invalid AltName")
	}
	if crt.SerialNumber == nil {
		t.Fatal("Empty Serial")
	}
	return nil
}
