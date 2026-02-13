package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func PatchFields(dto interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	v := reflect.ValueOf(dto)
	t := reflect.TypeOf(dto)

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() != reflect.Ptr || field.IsNil() {
			continue
		}
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		key := strings.Split(tag, ",")[0]
		result[key] = field.Elem().Interface()
	}
	return result
}

//  returns : set of valid json tag names for a struct
func AllowedJSONFields(dto interface{}) map[string]bool {
	allowed := map[string]bool{}
	t := reflect.TypeOf(dto)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		key := strings.Split(tag, ",")[0]
		allowed[key] = true
	}
	return allowed
}

// for checking raw JSON keys w/ the struct's json tags
func ValidateNoUnknownFields(body map[string]interface{}, dto interface{}) error {
	allowed := AllowedJSONFields(dto)
	var unknown []string
	for key := range body {
		if !allowed[key] {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("unknown field(s): %s", strings.Join(unknown, ", "))
	}
	return nil
}
