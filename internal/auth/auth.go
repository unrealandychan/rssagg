package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extract an API Key from the header of an HTTP request
// Example
// Authorization: API Key {insert API key here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("API Key not found in Authorization header")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first time of API Key")
	}
	return vals[1], nil
}
