//go:build go1.18
// +build go1.18

package env

// ParseAs parses the given struct type containing `env` tags and loads its
// values from environment variables.
func ParseAs[T any]() (T, error) {
	var t T
	return t, Parse(&t)
}

// ParseWithOptions parses the given struct type containing `env` tags and
// loads its values from environment variables.
func ParseAsWithOptions[T any](opts Options) (T, error) {
	var t T
	return t, ParseWithOptions(&t, opts)
}

// Must panic is if err is not nil, and returns t otherwise.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
