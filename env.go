package env

import (
	"encoding"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// nolint: gochecknoglobals
var (
	// ErrNotAStructPtr is returned if you pass something that is not a pointer to a
	// Struct to Parse
	ErrNotAStructPtr = errors.New("env: expected a pointer to a Struct")

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
			return i, err
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

	defaultTypeParsers = map[reflect.Type]ParserFunc{
		reflect.TypeOf(url.URL{}): func(v string) (interface{}, error) {
			u, err := url.Parse(v)
			if err != nil {
				return nil, fmt.Errorf("unable to parse URL: %v", err)
			}
			return *u, nil
		},
		reflect.TypeOf(time.Nanosecond): func(v string) (interface{}, error) {
			s, err := time.ParseDuration(v)
			if err != nil {
				return nil, fmt.Errorf("unable to parse duration: %v", err)
			}
			return s, err
		},
	}
)

// ParserFunc defines the signature of a function that can be used within `CustomParsers`
type ParserFunc func(v string) (interface{}, error)

// Parse parses a struct containing `env` tags and loads its values from
// environment variables.
func Parse(v interface{}) error {
	return ParseWithFuncs(v, map[reflect.Type]ParserFunc{})
}

// Decryptor is used to decrypt variables tagged with encrypted when using the
// ParseWithDecrypt function. It wraps Parse otherwise.
type Decryptor interface {
	Decrypt(val string) (string, error)
}

// ParseWithDecrypt is the same as `Parse` except it allows you to supply an decryptor
// to be used for decrypting any vars tagged with `encrypted`.
func ParseWithDecrypt(v interface{}, decryptor Decryptor) error {
	return ParseWithDecryptFuncs(v, map[reflect.Type]ParserFunc{}, decryptor)
}

// parseWithFuncsCommon is the common function for ParseWithFuncs and ParseWithDecryptFuncs
// to get ref, parses and return them.
func parseWithFuncsCommon(v interface{}, funcMap map[reflect.Type]ParserFunc) (reflect.Value, map[reflect.Type]ParserFunc, error) {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return reflect.Value{}, nil, ErrNotAStructPtr
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return reflect.Value{}, nil, ErrNotAStructPtr
	}
	var parsers = defaultTypeParsers
	for k, v := range funcMap {
		parsers[k] = v
	}
	return ref, parsers, nil
}

// ParseWithFuncs is the same as `Parse` except it also allows the user to pass
// in custom parsers.
func ParseWithFuncs(v interface{}, funcMap map[reflect.Type]ParserFunc) error {
	ref, parsers, err := parseWithFuncsCommon(v, funcMap)
	if err != nil {
		return err
	}
	return doParse(ref, parsers, nil)
}

// ParseWithDecryptFuncs is the same as `ParseWithDecrypt` except it also
// allows the user to pass in custom parsers.
func ParseWithDecryptFuncs(v interface{}, funcMap map[reflect.Type]ParserFunc, decryptor Decryptor) error {
	ref, parsers, err := parseWithFuncsCommon(v, funcMap)
	if err != nil {
		return err
	}
	return doParse(ref, parsers, decryptor)
}

func doParse(ref reflect.Value, funcMap map[reflect.Type]ParserFunc, decryptor Decryptor) error {
	var refType = ref.Type()

	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		if !refField.CanSet() {
			continue
		}
		if reflect.Ptr == refField.Kind() && !refField.IsNil() {
			err := ParseWithFuncs(refField.Interface(), funcMap)
			if err != nil {
				return err
			}
			continue
		}
		if reflect.Struct == refField.Kind() && refField.CanAddr() && refField.Type().Name() == "" {
			err := Parse(refField.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}
		refTypeField := refType.Field(i)
		value, err := get(refTypeField, decryptor)
		if err != nil {
			return err
		}
		if value == "" {
			if reflect.Struct == refField.Kind() {
				if err := doParse(refField, funcMap, decryptor); err != nil {
					return err
				}
			}
			continue
		}
		if err := set(refField, refTypeField, value, funcMap); err != nil {
			return err
		}
	}
	return nil
}

func get(field reflect.StructField, decryptor Decryptor) (val string, err error) {
	var required bool
	var exists bool
	var loadFile bool
	var decrypt bool
	var expand = strings.EqualFold(field.Tag.Get("envExpand"), "true")

	key, opts := parseKeyForOption(field.Tag.Get("env"))

	for _, opt := range opts {
		switch opt {
		case "":
			break
		case "file":
			loadFile = true
		case "required":
			required = true
		case "decrypt":
			decrypt = true
		default:
			return "", fmt.Errorf("env: tag option %q not supported", opt)
		}
	}

	defaultValue := field.Tag.Get("envDefault")
	val, exists = getOr(key, defaultValue)

	if decrypt && !loadFile {
		decryptedVal, err := decryptVal(val, decryptor)
		if err != nil {
			return "", err
		}
		val = decryptedVal
	}

	if expand {
		val = os.ExpandEnv(val)
	}

	if required && !exists {
		return "", fmt.Errorf(`env: required environment variable %q is not set`, key)
	}

	if loadFile && val != "" {
		filename := val
		val, err = getFromFile(filename)
		if err != nil {
			return "", fmt.Errorf(`env: could not load content of file "%s" from variable %s: %v`, filename, key, err)
		}

		if decrypt {
			decryptedVal, err := decryptVal(val, decryptor)
			if err != nil {
				return "", err
			}
			val = decryptedVal
		}
	}

	return val, err
}

// decryptVal will decrypt val using decryptor.
func decryptVal(val string, decryptor Decryptor) (string, error) {
	if decryptor == nil {
		return "", fmt.Errorf("env: detected decrypt tag on var but called with Parse. Use ParseWithDecrypt instead")
	}
	decryptedVal, err := decryptor.Decrypt(val)
	if err != nil {
		return "", fmt.Errorf("env: couldn't decrypt val using decryptor. %s", err.Error())
	}
	return decryptedVal, nil
}

// split the env tag's key into the expected key and desired option, if any.
func parseKeyForOption(key string) (string, []string) {
	opts := strings.Split(key, ",")
	return opts[0], opts[1:]
}

func getFromFile(filename string) (value string, err error) {
	b, err := ioutil.ReadFile(filename)
	return string(b), err
}

func getOr(key, defaultValue string) (value string, exists bool) {
	value, exists = os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value, exists
}

func set(field reflect.Value, sf reflect.StructField, value string, funcMap map[reflect.Type]ParserFunc) error {
	if field.Kind() == reflect.Slice {
		return handleSlice(field, value, sf, funcMap)
	}

	var tm = asTextUnmarshaler(field)
	if tm != nil {
		var err = tm.UnmarshalText([]byte(value))
		return newParseError(sf, err)
	}

	var typee = sf.Type
	var fieldee = field
	if typee.Kind() == reflect.Ptr {
		typee = typee.Elem()
		fieldee = field.Elem()
	}

	parserFunc, ok := funcMap[typee]
	if ok {
		val, err := parserFunc(value)
		if err != nil {
			return newParseError(sf, err)
		}

		fieldee.Set(reflect.ValueOf(val))
		return nil
	}

	parserFunc, ok = defaultBuiltInParsers[typee.Kind()]
	if ok {
		val, err := parserFunc(value)
		if err != nil {
			return newParseError(sf, err)
		}

		fieldee.Set(reflect.ValueOf(val).Convert(typee))
		return nil
	}

	return newNoParserError(sf)
}

func handleSlice(field reflect.Value, value string, sf reflect.StructField, funcMap map[reflect.Type]ParserFunc) error {
	var separator = sf.Tag.Get("envSeparator")
	if separator == "" {
		separator = ","
	}
	var parts = strings.Split(value, separator)

	var typee = sf.Type.Elem()
	if typee.Kind() == reflect.Ptr {
		typee = typee.Elem()
	}

	if _, ok := reflect.New(typee).Interface().(encoding.TextUnmarshaler); ok {
		return parseTextUnmarshalers(field, parts, sf)
	}

	parserFunc, ok := funcMap[typee]
	if !ok {
		parserFunc, ok = defaultBuiltInParsers[typee.Kind()]
		if !ok {
			return newNoParserError(sf)
		}
	}

	var result = reflect.MakeSlice(sf.Type, 0, len(parts))
	for _, part := range parts {
		r, err := parserFunc(part)
		if err != nil {
			return newParseError(sf, err)
		}
		var v = reflect.ValueOf(r).Convert(typee)
		if sf.Type.Elem().Kind() == reflect.Ptr {
			v = reflect.New(typee)
			v.Elem().Set(reflect.ValueOf(r).Convert(typee))
		}
		result = reflect.Append(result, v)
	}
	field.Set(result)
	return nil
}

func asTextUnmarshaler(field reflect.Value) encoding.TextUnmarshaler {
	if reflect.Ptr == field.Kind() {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
	} else if field.CanAddr() {
		field = field.Addr()
	}

	tm, ok := field.Interface().(encoding.TextUnmarshaler)
	if !ok {
		return nil
	}
	return tm
}

func parseTextUnmarshalers(field reflect.Value, data []string, sf reflect.StructField) error {
	s := len(data)
	elemType := field.Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(elemType), s, s)
	for i, v := range data {
		sv := slice.Index(i)
		kind := sv.Kind()
		if kind == reflect.Ptr {
			sv = reflect.New(elemType.Elem())
		} else {
			sv = sv.Addr()
		}
		tm := sv.Interface().(encoding.TextUnmarshaler)
		if err := tm.UnmarshalText([]byte(v)); err != nil {
			return newParseError(sf, err)
		}
		if kind == reflect.Ptr {
			slice.Index(i).Set(sv)
		}
	}

	field.Set(slice)

	return nil
}

func newParseError(sf reflect.StructField, err error) error {
	if err == nil {
		return nil
	}
	return parseError{
		sf:  sf,
		err: err,
	}
}

type parseError struct {
	sf  reflect.StructField
	err error
}

func (e parseError) Error() string {
	return fmt.Sprintf(`env: parse error on field "%s" of type "%s": %v`, e.sf.Name, e.sf.Type, e.err)
}

func newNoParserError(sf reflect.StructField) error {
	return fmt.Errorf(`env: no parser found for field "%s" of type "%s"`, sf.Name, sf.Type)
}
