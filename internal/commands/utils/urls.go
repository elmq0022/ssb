package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func BuildEndpoint(args ...string) (string, error) {
	baseUrl := MustReadConfig().URL
	endpoint, err := url.JoinPath(baseUrl, args...)
	if err != nil {
		return "", fmt.Errorf("couldn't build the url: %w", err)
	}
	return endpoint, nil
}

func BuildJSONRequest(method, endpoint string, data any) (*http.Request, error) {
	var body io.Reader

	if data != nil {
		payload, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(payload)
	}

	token := MustReadJWTToken().Token
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
