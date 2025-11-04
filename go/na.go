package na01

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

type api struct {
	path   string
	method string
	in     any
	out    any
}

var token string

func (entry api) invoke() error {
	for stage := range 2 {
		if len(token) == 0 {
			t, err := Auth()
			if err != nil {
				return err
			}
			token = t
		}
		method := entry.method
		if len(method) == 0 {
			method = http.MethodGet
		}
		var in io.Reader
		if entry.in != nil {
			data, err := json.Marshal(entry.in)
			if err != nil {
				return err
			}
			in = bytes.NewBuffer(data)
		}
		req, err := http.NewRequest(method, API+entry.path, in)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		if entry.in != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 403 {
			if stage == 0 {
				token = ""
				continue
			}
		}
		if resp.StatusCode >= 300 {
			return errors.New("HTTP Error: " + strconv.Itoa(resp.StatusCode))
		}
		if entry.out != nil {
			err := json.NewDecoder(resp.Body).Decode(entry.out)
			if err != nil {
				return err
			}
		}
		break
	}
	return nil
}

type entities[T any] struct {
	Count int `json:"count"`
	List  []T `json:"entities"`
}

type Zone struct {
	ID       int                `json:"id,omitzero"`
	Name     string             `json:"name"`
	Comment  string             `json:"comment,omitzero"`
	Count    int                `json:"records_count,omitzero"`
	TTL      int                `json:"ttl,omitzero"`
	Email    string             `json:"soa_email,omitzero"`
	Transfer bool               `json:"is_in_transfer,omitzero"`
	Tech     bool               `json:"is_technical_zone,omitzero"`
	DNS      entities[[]string] `json:"secondary_dns,omitzero"`
}

type RR struct {
	ID   int            `json:"id,omitzero"`
	Zone int            `json:"zone_id,omitzero"`
	Name string         `json:"name,omitzero"`
	Type string         `json:"type,omitzero"`
	TTL  int            `json:"ttl,omitzero"`
	Data map[string]any `json:"details,omitzero"`
}

type RRa struct {
	RR
	IP string `json:"ip"`
}

type RRtxt struct {
	RR
	Value string `json:"value"`
}

type RRcname struct {
	RR
	Domain string `json:"domain"`
}

func Zones() ([]Zone, error) {
	var r entities[Zone]
	err := api{path: "zones", out: &r}.invoke()
	if err != nil {
		return nil, err
	}
	return r.List, nil
}

func GetZone(id int) (Zone, error) {
	var z Zone
	err := api{path: "zones/" + strconv.Itoa(id), out: &z}.invoke()
	return z, err
}

func NewZone(name string) (Zone, error) {
	in := Zone{Name: name}
	var out Zone
	err := api{path: "zones", method: http.MethodPost, in: in, out: &out}.invoke()
	return out, err
}

func DropZone(id int) (Zone, error) {
	var z Zone
	err := api{path: "zones/" + strconv.Itoa(id), method: http.MethodDelete, out: &z}.invoke()
	return z, err
}

func ZoneRRs(id int) ([]RR, error) {
	var result []RR
	for {
		var r entities[RR]
		err := api{
			path: "zones/" + strconv.Itoa(id) + "/records/?offset=" + strconv.Itoa(len(result)),
			out:  &r}.invoke()
		if err != nil {
			return nil, err
		}
		if len(r.List) == 0 {
			break
		}
		result = append(result, r.List...)
	}
	return result, nil
}

func NewRR(rr any) (RR, error) {
	var res RR
	err := api{path: "records", method: http.MethodPost, in: rr, out: &res}.invoke()
	return res, err
}
