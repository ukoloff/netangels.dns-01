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

type Zone struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Comment  string `json:"comment"`
	Count    int    `json:"records_count"`
	TTL      int    `json:"ttl"`
	Email    string `json:"soa_email"`
	Transfer bool   `json:"is_in_transfer"`
	Tech     bool   `json:"is_technical_zone"`
	DNS      struct {
		List []string `json:"entities"`
	} `json:"secondary_dns"`
}

func Zones() ([]Zone, error) {
	token, err := Auth()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, API+"/zones", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r struct {
		Count int    `json:"count"`
		List  []Zone `json:"entities"`
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r.List, nil
}
