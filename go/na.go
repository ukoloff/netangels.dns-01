package na01

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
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
	req := url.Values{}
	req.Add("api_key", os.Getenv(ENV_API))
	client := &http.Client{}
	resp, err := client.PostForm(AUTH, req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	var jj struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jj); err != nil {
		return "", err
	}
	return jj.Token, nil
}
