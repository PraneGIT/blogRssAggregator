package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Autharization: ApiKey {insert apikey here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no autharization header")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid autharization header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("invalid first part of header")
	}
	return vals[1], nil
}
