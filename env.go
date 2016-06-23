package env

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ErrNotAStructPtr is returned if you pass something that is not a pointer to a
// Struct to Parse
var ErrNotAStructPtr = errors.New("Expected a pointer to a Struct")

// ErrUnsupportedType if the struct field type is not supported by env
var ErrUnsupportedType = errors.New("Type is not supported")

// ErrUnsupportedSliceType if the slice element type is not supported by env
var ErrUnsupportedSliceType = errors.New("Unsupported slice type")

// Friendly names for reflect types
var sliceOfInts = reflect.TypeOf([]int(nil))
var sliceOfStrings = reflect.TypeOf([]string(nil))
var sliceOfBools = reflect.TypeOf([]bool(nil))

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
	return doParse(ref, val)
}

func doParse(ref reflect.Value, val interface{}) error {
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		value, err := get(refType.Field(i))
		if err != nil {
			return err
		}
		if value == "" {
			continue
		}
		if err := set(ref.Field(i), refType.Field(i), value); err != nil {
			return err
		}
	}
	return nil
}

func get(field reflect.StructField) (string, error) {
	var (
		val string
		err error
	)

	key, opt := parseKeyForOption(field.Tag.Get("env"))
	// The only option supported is "required".
	switch opt {
	case "":
		defaultValue := field.Tag.Get("envDefault")
		val = getOr(key, defaultValue)
	case "required":
		val, err = getRequired(key)
	default:
		err = errors.New("Env tag option " + opt + " not supported.")
	}
	return val, err
}

// split the env tag's key into the expected key and desired option, if any.
func parseKeyForOption(key string) (string, string) {
	i := strings.Index(key, ",")
	if i == -1 {
		return key, ""
	}
	if i == len(key)-1 {
		return key[:i], ""
	}
	return key[:i], key[i+1:]
}

func getRequired(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	// We do not use fmt.Errorf to avoid another import.
	return "", errors.New("Required environment variable " + key + " is not set")
}

func getOr(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func set(field reflect.Value, refType reflect.StructField, value string) error {
	switch field.Kind() {
	case reflect.Slice:
		separator := refType.Tag.Get("envSeparator")
		return handleSlice(field, value, separator)
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
	default:
		return ErrUnsupportedType
	}
	return nil
}

func handleSlice(field reflect.Value, value, separator string) error {
	if separator == "" {
		separator = ","
	}

	splitData := strings.Split(value, separator)

	switch field.Type() {
	case sliceOfStrings:
		field.Set(reflect.ValueOf(splitData))
	case sliceOfInts:
		intData, err := parseInts(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(intData))
	case sliceOfBools:
		boolData, err := parseBools(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(boolData))
	default:
		return ErrUnsupportedSliceType
	}
	return nil
}

func parseInts(data []string) ([]int, error) {
	var intSlice []int

	for _, v := range data {
		intValue, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, int(intValue))
	}
	return intSlice, nil
}

func parseBools(data []string) ([]bool, error) {
	var boolSlice []bool

	for _, v := range data {
		bvalue, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}

		boolSlice = append(boolSlice, bvalue)
	}
	return boolSlice, nil
}
