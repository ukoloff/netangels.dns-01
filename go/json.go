package main

import "encoding/json"

func Stringify(data any) (string, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
