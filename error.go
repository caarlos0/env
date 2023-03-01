package env

import (
	"fmt"
	"reflect"
	"strings"
)

// An aggregated error wrapper to combine gathered errors. This allows either to display all errors or convert them individually
// List of the available errors
// ParseError
// NotStructPtrError
// NoParserError
// NoSupportedTagOptionError
// EnvVarIsNotSetError
// EmptyEnvVarError
// LoadFileContentError
// ParseValueError
type AggregateError struct {
	Errors []error
}

func newAggregateError(initErr error) error {
	return AggregateError{
		[]error{
			initErr,
		},
	}
}

func (e AggregateError) Error() string {
	var sb strings.Builder

	sb.WriteString("env:")

	for _, err := range e.Errors {
		sb.WriteString(fmt.Sprintf(" %v;", err.Error()))
	}

	return strings.TrimRight(sb.String(), ";")
}

// Is conforms with errors.Is.
func (e AggregateError) Is(err error) bool {
	for _, ie := range e.Errors {
		if reflect.TypeOf(ie) == reflect.TypeOf(err) {
			return true
		}
	}
	return false
}

// The error occurs when it's impossible to convert the value for given type.
type ParseError struct {
	Name string
	Type reflect.Type
	Err  error
}

func newParseError(sf reflect.StructField, err error) error {
	return ParseError{sf.Name, sf.Type, err}
}

func (e ParseError) Error() string {
	return fmt.Sprintf(`parse error on field "%s" of type "%s": %v`, e.Name, e.Type, e.Err)
}

// The error occurs when pass something that is not a pointer to a Struct to Parse
type NotStructPtrError struct{}

func (e NotStructPtrError) Error() string {
	return "expected a pointer to a Struct"
}

// This error occurs when there is no parser provided for given type
// Supported types and defaults: https://github.com/caarlos0/env#supported-types-and-defaults
// How to create a custom parser: https://github.com/caarlos0/env#custom-parser-funcs
type NoParserError struct {
	Name string
	Type reflect.Type
}

func newNoParserError(sf reflect.StructField) error {
	return NoParserError{sf.Name, sf.Type}
}

func (e NoParserError) Error() string {
	return fmt.Sprintf(`no parser found for field "%s" of type "%s"`, e.Name, e.Type)
}

// This error occurs when the given tag is not supported
// In-built supported tags: "", "file", "required", "unset", "notEmpty", "envDefault", "envExpand", "envSeparator"
// How to create a custom tag: https://github.com/caarlos0/env#changing-default-tag-name
type NoSupportedTagOptionError struct {
	Tag string
}

func newNoSupportedTagOptionError(tag string) error {
	return NoSupportedTagOptionError{tag}
}

func (e NoSupportedTagOptionError) Error() string {
	return fmt.Sprintf("tag option %q not supported", e.Tag)
}

// This error occurs when the required variable is not set
// Read about required fields: https://github.com/caarlos0/env#required-fields
type EnvVarIsNotSetError struct {
	Key string
}

func newEnvVarIsNotSet(key string) error {
	return EnvVarIsNotSetError{key}
}

func (e EnvVarIsNotSetError) Error() string {
	return fmt.Sprintf(`required environment variable %q is not set`, e.Key)
}

// This error occurs when the variable which must be not empty is existing but has an empty value
// Read about not empty fields: https://github.com/caarlos0/env#not-empty-fields
type EmptyEnvVarError struct {
	Key string
}

func newEmptyEnvVarError(key string) error {
	return EmptyEnvVarError{key}
}

func (e EmptyEnvVarError) Error() string {
	return fmt.Sprintf("environment variable %q should not be empty", e.Key)
}

// This error occurs when it's impossible to load the value from the file
// Read about From file feature: https://github.com/caarlos0/env#from-file
type LoadFileContentError struct {
	Filename string
	Key      string
	Err      error
}

func newLoadFileContentError(filename, key string, err error) error {
	return LoadFileContentError{filename, key, err}
}

func (e LoadFileContentError) Error() string {
	return fmt.Sprintf(`could not load content of file "%s" from variable %s: %v`, e.Filename, e.Key, e.Err)
}

// This error occurs when it's impossible to convert value using given parser
// Supported types and defaults: https://github.com/caarlos0/env#supported-types-and-defaults
// How to create a custom parser: https://github.com/caarlos0/env#custom-parser-funcs
type ParseValueError struct {
	Msg string
	Err error
}

func newParseValueError(message string, err error) error {
	return ParseValueError{message, err}
}

func (e ParseValueError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}
