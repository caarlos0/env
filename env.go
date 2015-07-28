package env

import (
	"errors"
	"os"
	"reflect"
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

// ErrNotAStruct is returned if you pass something that is not a Struct to
// ParseEnv
var ErrNotAStruct = errors.New("Expected a Struct")

// ParseEnv parses a struct containing `env` tags and loads its values from
// environment variables.
func ParseEnv(t interface{}, v interface{}) error {
	it := reflect.TypeOf(t)
	for i := 0; i < it.NumField(); i++ {
		field := it.Field(i)
		value := Get(field.Tag.Get("env"))
		reflect.ValueOf(v).Elem().FieldByName(field.Name).SetString(value)
	}
	return nil
}
