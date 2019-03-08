package env

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/caarlos0/env/parsers"
)

// nolint: gochecknoglobals
var (
	// ErrNotAStructPtr is returned if you pass something that is not a pointer to a
	// Struct to Parse
	ErrNotAStructPtr = errors.New("expected a pointer to a Struct")
	// ErrUnsupportedType if the struct field type is not supported by env
	ErrUnsupportedType = errors.New("type is not supported")
	// ErrUnsupportedSliceType if the slice element type is not supported by env
	ErrUnsupportedSliceType = errors.New("unsupported slice type")
	// OnEnvVarSet is an optional convenience callback, such as for logging purposes.
	// If not nil, it's called after successfully setting the given field from the given value.
	OnEnvVarSet func(reflect.StructField, string)

	defaultBuiltInParsers = map[reflect.Kind]ParserFunc{
		reflect.Bool: func(v string) (interface{}, error) {
			return strconv.ParseBool(v)
		},
		reflect.String: func(v string) (interface{}, error) {
			return v, nil
		},
		reflect.Int: func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 32)
			return int(i), err
		},
		reflect.Int16: func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 16)
			return int16(i), err
		},
		reflect.Int32: func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 32)
			return int32(i), err
		},
		reflect.Int64: func(v string) (interface{}, error) {
			return strconv.ParseInt(v, 10, 64)
		},
		reflect.Int8: func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 8)
			return int8(i), err
		},
		reflect.Uint: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 32)
			return uint(i), err
		},
		reflect.Uint16: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 16)
			return uint16(i), err
		},
		reflect.Uint32: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 32)
			return uint32(i), err
		},
		reflect.Uint64: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 64)
			return uint64(i), err
		},
		reflect.Uint8: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 8)
			return uint8(i), err
		},
		reflect.Float64: func(v string) (interface{}, error) {
			return strconv.ParseFloat(v, 64)
		},
		reflect.Float32: func(v string) (interface{}, error) {
			f, err := strconv.ParseFloat(v, 32)
			return float32(f), err
		},
	}
)

func defaultCustomParsers() CustomParsers {
	return CustomParsers{
		parsers.URLType:      parsers.URLFunc,
		parsers.DurationType: parsers.DurationFunc,
	}
}

// CustomParsers is a friendly name for the type that `ParseWithFuncs()` accepts
type CustomParsers map[reflect.Type]ParserFunc

// ParserFunc defines the signature of a function that can be used within `CustomParsers`
type ParserFunc func(v string) (interface{}, error)

// Parse parses a struct containing `env` tags and loads its values from
// environment variables.
func Parse(v interface{}) error {
	return ParseWithFuncs(v, defaultCustomParsers())
}

// ParseWithFuncs is the same as `Parse` except it also allows the user to pass
// in custom parsers.
func ParseWithFuncs(v interface{}, funcMap CustomParsers) error {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotAStructPtr
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrNotAStructPtr
	}
	return doParse(ref, funcMap)
}

func doParse(ref reflect.Value, funcMap CustomParsers) error {
	refType := ref.Type()
	var errorList []string

	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		if reflect.Ptr == refField.Kind() && !refField.IsNil() && refField.CanSet() {
			err := Parse(refField.Interface())
			if nil != err {
				return err
			}
			continue
		}
		refTypeField := refType.Field(i)
		value, err := get(refTypeField)
		if err != nil {
			errorList = append(errorList, err.Error())
			continue
		}
		if value == "" {
			if reflect.Struct == refField.Kind() {
				err := doParse(refField, funcMap)
				if nil != err {
					errorList = append(errorList, err.Error())
				}
			}
			continue
		}
		if err := set(refField, refTypeField, value, funcMap); err != nil {
			errorList = append(errorList, err.Error())
			continue
		}
		if OnEnvVarSet != nil {
			OnEnvVarSet(refTypeField, value)
		}
	}
	if len(errorList) == 0 {
		return nil
	}
	return errors.New(strings.Join(errorList, ". "))
}

func get(field reflect.StructField) (string, error) {
	var (
		val string
		err error
	)

	key, opts := parseKeyForOption(field.Tag.Get("env"))

	defaultValue := field.Tag.Get("envDefault")
	val = getOr(key, defaultValue)

	expandVar := field.Tag.Get("envExpand")
	if strings.ToLower(expandVar) == "true" {
		val = os.ExpandEnv(val)
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			// The only option supported is "required".
			switch opt {
			case "":
				break
			case "required":
				val, err = getRequired(key)
			default:
				err = fmt.Errorf("env tag option %q not supported", opt)
			}
		}
	}

	return val, err
}

// split the env tag's key into the expected key and desired option, if any.
func parseKeyForOption(key string) (string, []string) {
	opts := strings.Split(key, ",")
	return opts[0], opts[1:]
}

func getRequired(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("required environment variable %q is not set", key)
}

func getOr(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return defaultValue
}

func set(field reflect.Value, refType reflect.StructField, value string, funcMap CustomParsers) error {
	parserFunc, ok := funcMap[refType.Type]
	if ok {
		val, err := parserFunc(value)
		if err != nil {
			return fmt.Errorf("custom parser error: %v", err)
		}
		field.Set(reflect.ValueOf(val))
		return nil
	}

	parserFunc, ok = defaultBuiltInParsers[field.Kind()]
	if ok {
		val, err := parserFunc(value)
		if err != nil {
			return fmt.Errorf("parser error: %v", err)
		}
		field.Set(reflect.ValueOf(val))
		return nil
	}

	if field.Kind() == reflect.Slice {
		return handleSlice(field, value, refType.Tag.Get("envSeparator"), funcMap)
	}

	return handleTextUnmarshaler(field, value)
}

func handleSlice(field reflect.Value, value, separator string, funcMap CustomParsers) error {
	if separator == "" {
		separator = ","
	}
	parts := strings.Split(value, separator)
	result := reflect.MakeSlice(field.Type(), 0, len(parts))

	parserFunc, ok := funcMap[field.Type().Elem()]
	if !ok {
		parserFunc, ok = defaultBuiltInParsers[field.Type().Elem().Kind()]
		if !ok {
			return fmt.Errorf("no parser for slice of %s", field.Type().Elem().Kind().String())
		}
	}

	for _, part := range parts {
		r, err := parserFunc(part)
		if err != nil {
			return fmt.Errorf("parser error: %v", err)
		}
		result = reflect.Append(result, reflect.ValueOf(r))
	}

	field.Set(result)
	return nil
}

func handleTextUnmarshaler(field reflect.Value, value string) error {
	if reflect.Ptr == field.Kind() {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
	} else if field.CanAddr() {
		field = field.Addr()
	}

	tm, ok := field.Interface().(encoding.TextUnmarshaler)
	if !ok {
		return ErrUnsupportedType
	}

	return tm.UnmarshalText([]byte(value))
}
