package auth

import (
	"errors"
	"net/http"
	"strings"
)

// format => Authorization: ApiKey {apiKey}
func GetApiKey(headers http.Header) (string, error) {
	key := headers.Get("Authorization")
	if key == "" {
		return "", errors.New("no authentication information")
	}

	vals := strings.Split(key, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}

	return vals[1], nil
}
