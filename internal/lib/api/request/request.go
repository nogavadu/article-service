package request

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func IsStructEmpty(obj interface{}) (bool, error) {
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return false, fmt.Errorf("expected a struct type")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			return false, nil
		}

		if !field.IsZero() {
			return false, nil
		}
	}

	return true, nil
}

func GetAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization header")
	}

	token := tokenParts[1]
	return token, nil
}
