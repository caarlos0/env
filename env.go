package env

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

// Set alias to os.Setenv
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset alias to os.Unsetenv
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Get alias to os.Getenv
func Get(key string) string {
	return os.Getenv(key)
}

// GetOr alias to os.Getenv, returning the given default value case it's not set
func GetOr(key, defaultValue string) string {
	value := Get(key)
	if value != "" {
		return value
	}
	return defaultValue
}

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
		field := refType.Field(i)
		value := os.Getenv(field.Tag.Get("env"))
		if value == "" {
			continue
		}
		if err := set(ref.Field(i), value); err != nil {
			return err
		}
	}
	return nil
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
