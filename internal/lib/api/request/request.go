package request

import (
	"fmt"
	"reflect"
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
