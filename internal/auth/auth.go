package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeaderVal := headers.Get("Authorization")
	if authHeaderVal == "" {
		return "", errors.New("no Authorization Header found")
	}

	authStrSplit := strings.Split(authHeaderVal, " ")
	if len(authStrSplit) != 2 || authStrSplit[0] != "ApiKey" {
		return "", errors.New("invalid Authorization Header")
	}
	return authStrSplit[1], nil
}
