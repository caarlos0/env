// Package parsers contains custom parser funcs for common, non-built-in types
package parsers

import (
	"fmt"
	"net/url"
	"reflect"
	"time"
)

// nolint: gochecknoglobals
var (
	// URLType is a helper var that represents the `reflect.Type` of `url.URL`
	URLType = reflect.TypeOf(url.URL{})

	//  DurationType is a helper var that represents the `reflect.Type` of `time.Duration`
	DurationType = reflect.TypeOf(time.Nanosecond)
)

// URLFunc is a basic parser for the url.URL type that should be used with `env.ParseWithFuncs()`
func URLFunc(v string) (interface{}, error) {
	u, err := url.Parse(v)
	if err != nil {
		return nil, fmt.Errorf("unable parse URL: %v", err)
	}
	return *u, nil
}

// DurationFunc is a basic parser for the time.Duration type that should be used with `env.ParseWithFuncs()`
func DurationFunc(v string) (interface{}, error) {
	s, err := time.ParseDuration(v)
	if err != nil {
		return nil, fmt.Errorf("unable to parser duration: %v", err)
	}
	return s, err
}
