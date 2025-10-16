package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type RequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    any
}

func NewRequestBuilder(method, url string) *RequestBuilder {
	return &RequestBuilder{
		method:  method,
		url:     url,
		headers: make(map[string]string),
	}
}

func (b *RequestBuilder) WithJSON(body any) *RequestBuilder {
	b.body = body
	b.headers["Content-Type"] = "application/json"
	return b
}

func (b *RequestBuilder) WithAuth() *RequestBuilder {
	token := MustReadJWTToken().Token
	b.headers["Authorization"] = "Bearer " + token
	return b
}

func (b *RequestBuilder) Build() (*http.Request, error) {
	var buf *bytes.Buffer
	if b.body != nil {
		data, err := json.Marshal(b.body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		buf = bytes.NewBuffer(data)
	} else {
		buf = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(b.method, b.url, buf)
	if err != nil {
		return nil, err
	}

	for k, v := range b.headers {
		req.Header.Set(k, v)
	}

	return req, nil
}
