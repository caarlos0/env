package env

import (
	"errors"
	"reflect"
	"strconv"
)

// ErrNotAStructPtr is returned if you pass something that is not a pointer to a
// Struct to Parse
var ErrNotAStructPtr = errors.New("Expected a pointer to a Struct")

// ErrUnsuportedType if the struct field type is not supported by env
var ErrUnsuportedType = errors.New("Type is not supported")

// Parse parses a struct containing `env` tags and loads its values from
// environment variables.
func Parse(val interface{}) error {
	ptrRef := reflect.ValueOf(val)
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotAStructPtr
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrNotAStructPtr
	}
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		value := get(refType.Field(i))
		if value == "" {
			continue
		}
		if err := set(ref.Field(i), value); err != nil {
			return err
		}
	}
	return nil
}

func get(field reflect.StructField) string {
	defaultValue := field.Tag.Get("default")
	return GetOr(field.Tag.Get("env"), defaultValue)
}

func set(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(bvalue)
	case reflect.Int:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	}
	return nil
}
