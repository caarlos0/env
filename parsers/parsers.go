// Package parsers contains custom parser funcs for common, non-built-in types
package parsers

import (
	"fmt"
	"net/url"
	"reflect"
)

var (
	// URLType is a helper var that represents the `reflect.Type`` of `url.URL`
	URLType = reflect.TypeOf(url.URL{})
)

// URLFunc is a basic parser for the url.URL type that should be used with `env.ParseWithFuncs()`
func URLFunc(v string) (interface{}, error) {
	u, err := url.Parse(v)
	if err != nil {
		return nil, fmt.Errorf("Unable to complete URL parse: %v", err)
	}

	return *u, nil
}
