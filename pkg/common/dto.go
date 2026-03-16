package common

import (
	"fmt"
	"reflect"
	"strings"
)

func DTOSchemaValidation(v any) error {
	if v == nil {
		return fmt.Errorf("expected struct or pointer to struct")
	}

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("nil pointer provided")
		}

		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct or pointer to struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		bindingTag := fieldType.Tag.Get("binding")
		if !strings.Contains(bindingTag, "required") {
			continue
		}

		if fieldValue.IsZero() {
			jsonName := fieldType.Tag.Get("json")
			if jsonName == "" {
				jsonName = fieldType.Name
			}
			return fmt.Errorf("field '%v' is required", jsonName)
		}
	}

	return nil
}
