package env

import (
	"fmt"
	"reflect"
	"strings"
)

// List of available errors
// 	AggregateError - aggregates and contains errors below:
// 		ParseError
//		NotStructPtrErrorType
//  	NoParserError
//		NoSupportedTagOptionError
//		EnvVarIsNotSetError
//		LoadFileContentError
// 		ParseValueError

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

type ParseError struct {
	msg string
}

func newParseError(sf reflect.StructField, err error) error {
	return ParseError{fmt.Sprintf(`parse error on field "%s" of type "%s": %v`, sf.Name, sf.Type, err)}
}

func (e ParseError) Error() string {
	return e.msg
}

// ErrNotAStructPtr is returned if you pass something that is not a pointer to a Struct to Parse.
type NotStructPtrError struct {
	msg string
}

func newNotStructPtrError() error {
	return NotStructPtrError{"expected a pointer to a Struct"}
}

func (e NotStructPtrError) Error() string {
	return e.msg
}

type NoParserError struct {
	msg string
}

func newNoParserError(sf reflect.StructField) error {
	return NoParserError{fmt.Sprintf(`no parser found for field "%s" of type "%s"`, sf.Name, sf.Type)}
}

func (e NoParserError) Error() string {
	return e.msg
}

type NoSupportedTagOptionError struct {
	msg string
}

func newNoSupportedTagOptionError(tag string) error {
	return NoSupportedTagOptionError{fmt.Sprintf("tag option %q not supported", tag)}
}

func (e NoSupportedTagOptionError) Error() string {
	return e.msg
}

type EnvVarIsNotSetError struct {
	msg string
}

func newEnvVarIsNotSet(key string) error {
	return EnvVarIsNotSetError{fmt.Sprintf(`required environment variable %q is not set`, key)}
}

func (e EnvVarIsNotSetError) Error() string {
	return e.msg
}

type EmptyEnvVarError struct {
	msg string
}

func newEmptyEnvVarError(key string) error {
	return EmptyEnvVarError{fmt.Sprintf("environment variable %q should not be empty", key)}
}

func (e EmptyEnvVarError) Error() string {
	return e.msg
}

type LoadFileContentError struct {
	msg string
}

func newLoadFileContentError(filename, key string, err error) error {
	return LoadFileContentError{fmt.Sprintf(`could not load content of file "%s" from variable %s: %v`, filename, key, err)}
}

func (e LoadFileContentError) Error() string {
	return e.msg
}

type ParseValueError struct {
	msg string
}

func newParseValueError(message string) error {
	return ParseValueError{message}
}

func (e ParseValueError) Error() string {
	return e.msg
}
