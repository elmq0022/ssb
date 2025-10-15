package utils

import (
	"fmt"
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
