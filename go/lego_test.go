package na01_test

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"na01"
	"os"
	"os/exec"
	"strings"
)

func lego(provider string, env map[string]string) error {
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
		return err
	}
	err = x509check(domain)
	return err
}

func x509check(domain string) error {
	for _, ext := range []string{"key", "pfx", "crt", "json"} {
		_, err := os.Stat("./.lego/certificates/" + domain + "." + ext)
		if err != nil {
			return err
		}
	}

	blob, err := os.ReadFile("./.lego/certificates/" + domain + ".crt")
	if err != nil {
		return err
	}
	der, _ := pem.Decode(blob)
	crt, err := x509.ParseCertificate(der.Bytes)
	if err != nil {
		return err
	}
	if crt.Subject.CommonName != domain {
		return errors.New("Invalid CN")
	}
	if crt.DNSNames[0] != domain {
		return errors.New("Invalid AltName")
	}
	if crt.SerialNumber == nil {
		return errors.New("Empty Serial")
	}
	return nil
}
