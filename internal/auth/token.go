package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	if headers.Get("Authorization") == "" {
		return "", fmt.Errorf("No Authorization Token in Header")
	}
	token := strings.TrimPrefix(headers.Get("Authorization"), "Bearer ")

	return token, nil
}
