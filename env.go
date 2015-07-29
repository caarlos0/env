package env

import "os"

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
