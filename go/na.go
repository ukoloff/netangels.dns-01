package na01

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	AUTH    = "https://panel.netangels.ru/api/gateway/token/"
	API     = "https://api-ms.netangels.ru/api/v1/dns/"
	ENV_API = "NETANGELS_API_KEY"
)

type auth struct {
	Token string `json:"token"`
}

func Auth() (string, error) {
	body := url.Values{}
	body.Add("api_key", os.Getenv(ENV_API))
	r, _ := http.NewRequest(http.MethodPost, AUTH, strings.NewReader(body.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	j, _ := io.ReadAll(resp.Body)
	var jj struct {
		Token string `json:"token"`
	}
	json.Unmarshal(j, &jj)
	return jj.Token, nil
}
