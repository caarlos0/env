// Package env is a simple, zero-dependencies library to parse environment
// variables into structs.
//
// Example:
//
//	type config struct {
//		Home string `env:"HOME"`
//	}
//	// parse
//	var cfg config
//	err := env.Parse(&cfg)
//	// or parse with generics
//	cfg, err := env.ParseAs[config]()
//
// Check the examples and README for more detailed usage.
package env

import (
	"encoding"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// nolint: gochecknoglobals
var (
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
)

func defaultTypeParsers() map[reflect.Type]ParserFunc {
	return map[reflect.Type]ParserFunc{
		reflect.TypeOf(url.URL{}):       parseURL,
		reflect.TypeOf(time.Nanosecond): parseDuration,
		reflect.TypeOf(time.Location{}): parseLocation,
	}
}

func parseURL(v string) (interface{}, error) {
	u, err := url.Parse(v)
	if err != nil {
		return nil, newParseValueError("unable to parse URL", err)
	}
	return *u, nil
}

func parseDuration(v string) (interface{}, error) {
	d, err := time.ParseDuration(v)
	if err != nil {
		return nil, newParseValueError("unable to parse duration", err)
	}
	return d, err
}

func parseLocation(v string) (interface{}, error) {
	loc, err := time.LoadLocation(v)
	if err != nil {
		return nil, newParseValueError("unable to parse location", err)
	}
	return *loc, nil
}

// ParserFunc defines the signature of a function that can be used within
// `Options`' `FuncMap`.
type ParserFunc func(v string) (interface{}, error)

// OnSetFn is a hook that can be run when a value is set.
type OnSetFn func(tag string, value interface{}, isDefault bool)

// processFieldFn is a function which takes all information about a field and processes it.
type processFieldFn func(
	refField reflect.Value,
	refTypeField reflect.StructField,
	opts Options,
	fieldParams FieldParams,
) error

// Options for the parser.
type Options struct {
	// Environment keys and values that will be accessible for the service.
	Environment map[string]string

	// TagName specifies another tag name to use rather than the default 'env'.
	TagName string

	// PrefixTagName specifies another prefix tag name to use rather than the default 'envPrefix'.
	PrefixTagName string

	// DefaultValueTagName specifies another default tag name to use rather than the default 'envDefault'.
	DefaultValueTagName string

	// RequiredIfNoDef automatically sets all fields as required if they do not
	// declare 'envDefault'.
	RequiredIfNoDef bool

	// OnSet allows to run a function when a value is set.
	OnSet OnSetFn

	// Prefix define a prefix for every key.
	Prefix string

	// UseFieldNameByDefault defines whether or not `env` should use the field
	// name by default if the `env` key is missing.
	// Note that the field name will be "converted" to conform with environment
	// variable names conventions.
	UseFieldNameByDefault bool

	// SetDefaultsForZeroValuesOnly defines whether to set defaults for zero values
	// If the `env` variable for the value is not set
	// and `envDefault` is set
	// and the value is not a zero value for the the type
	// and SetDefaultsForZeroValuesOnly=true
	// the value from `envDefault` will be ignored
	// Usefull for mixing default values from `envDefault` and struct initialization
	SetDefaultsForZeroValuesOnly bool

	// Custom parse functions for different types.
	FuncMap map[reflect.Type]ParserFunc

	// Used internally. maps the env variable key to its resolved string value.
	// (for env var expansion)
	rawEnvVars map[string]string
}

func (opts *Options) getRawEnv(s string) string {
	val := opts.rawEnvVars[s]
	if val == "" {
		val = opts.Environment[s]
	}
	return os.Expand(val, opts.getRawEnv)
}

func defaultOptions() Options {
	return Options{
		TagName:             "env",
		PrefixTagName:       "envPrefix",
		DefaultValueTagName: "envDefault",
		Environment:         toMap(os.Environ()),
		FuncMap:             defaultTypeParsers(),
		rawEnvVars:          make(map[string]string),
	}
}

func mergeOptions[T any](target, source *T) {
	targetPtr := reflect.ValueOf(target).Elem()
	sourcePtr := reflect.ValueOf(source).Elem()

	targetType := targetPtr.Type()
	for i := 0; i < targetPtr.NumField(); i++ {
		fieldName := targetType.Field(i).Name
		targetField := targetPtr.Field(i)
		sourceField := sourcePtr.FieldByName(fieldName)

		if targetField.CanSet() && !isZero(sourceField) {
			// FuncMaps are being merged, while Environments must be overwritten
			if fieldName == "FuncMap" {
				if !sourceField.IsZero() {
					iter := sourceField.MapRange()
					for iter.Next() {
						targetField.SetMapIndex(iter.Key(), iter.Value())
					}
				}
			} else {
				targetField.Set(sourceField)
			}
		}
	}
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	default:
		zero := reflect.Zero(v.Type())
		return v.Interface() == zero.Interface()
	}
}

func customOptions(opts Options) Options {
	defOpts := defaultOptions()
	mergeOptions(&defOpts, &opts)
	return defOpts
}

func optionsWithPrefix(opts Options, prefix string) Options {
	return Options{
		Environment:                  opts.Environment,
		TagName:                      opts.TagName,
		PrefixTagName:                opts.PrefixTagName,
		DefaultValueTagName:          opts.DefaultValueTagName,
		RequiredIfNoDef:              opts.RequiredIfNoDef,
		OnSet:                        opts.OnSet,
		Prefix:                       prefix,
		UseFieldNameByDefault:        opts.UseFieldNameByDefault,
		SetDefaultsForZeroValuesOnly: opts.SetDefaultsForZeroValuesOnly,
		FuncMap:                      opts.FuncMap,
		rawEnvVars:                   opts.rawEnvVars,
	}
}

func optionsWithSliceEnvPrefix(opts Options, index int) Options {
	return optionsWithPrefix(opts, fmt.Sprintf("%s%d_", opts.Prefix, index))
}

func optionsWithMapEnvPrefix(opts Options, mapKey string) Options {
	return optionsWithPrefix(opts, fmt.Sprintf("%s%s_", opts.Prefix, mapKey))
}

func optionsWithEnvPrefix(field reflect.StructField, opts Options) Options {
	return optionsWithPrefix(opts, opts.Prefix+field.Tag.Get(opts.PrefixTagName))
}

// Parse parses a struct containing `env` tags and loads its values from
// environment variables.
func Parse(v interface{}) error {
	return parseInternal(v, setField, defaultOptions())
}

// ParseWithOptions parses a struct containing `env` tags and loads its values from
// environment variables.
func ParseWithOptions(v interface{}, opts Options) error {
	return parseInternal(v, setField, customOptions(opts))
}

// ParseAs parses the given struct type containing `env` tags and loads its
// values from environment variables.
func ParseAs[T any]() (T, error) {
	var t T
	err := Parse(&t)
	return t, err
}

// ParseWithOptions parses the given struct type containing `env` tags and
// loads its values from environment variables.
func ParseAsWithOptions[T any](opts Options) (T, error) {
	var t T
	err := ParseWithOptions(&t, opts)
	return t, err
}

// Must panic is if err is not nil, and returns t otherwise.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// GetFieldParams parses a struct containing `env` tags and returns information about
// tags it found.
func GetFieldParams(v interface{}) ([]FieldParams, error) {
	return GetFieldParamsWithOptions(v, defaultOptions())
}

// GetFieldParamsWithOptions parses a struct containing `env` tags and returns information about
// tags it found.
func GetFieldParamsWithOptions(v interface{}, opts Options) ([]FieldParams, error) {
	var result []FieldParams
	err := parseInternal(
		v,
		func(_ reflect.Value, _ reflect.StructField, _ Options, fieldParams FieldParams) error {
			if fieldParams.OwnKey != "" {
				result = append(result, fieldParams)
			}
			return nil
		},
		customOptions(opts),
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseInternal(v interface{}, processField processFieldFn, opts Options) error {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return newAggregateError(NotStructPtrError{})
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return newAggregateError(NotStructPtrError{})
	}

	return doParse(ref, processField, opts)
}

func doParse(ref reflect.Value, processField processFieldFn, opts Options) error {
	refType := ref.Type()

	var agrErr AggregateError

	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		refTypeField := refType.Field(i)

		if err := doParseField(refField, refTypeField, processField, opts); err != nil {
			if val, ok := err.(AggregateError); ok {
				agrErr.Errors = append(agrErr.Errors, val.Errors...)
			} else {
				agrErr.Errors = append(agrErr.Errors, err)
			}
		}
	}

	if len(agrErr.Errors) == 0 {
		return nil
	}

	return agrErr
}

func doParseField(
	refField reflect.Value,
	refTypeField reflect.StructField,
	processField processFieldFn,
	opts Options,
) error {
	if !refField.CanSet() {
		return nil
	}
	if refField.Kind() == reflect.Ptr && refField.Elem().Kind() == reflect.Struct && !refField.IsNil() {
		return parseInternal(refField.Interface(), processField, optionsWithEnvPrefix(refTypeField, opts))
	}
	if refField.Kind() == reflect.Struct && refField.CanAddr() && refField.Type().Name() == "" {
		return parseInternal(refField.Addr().Interface(), processField, optionsWithEnvPrefix(refTypeField, opts))
	}

	params, err := parseFieldParams(refTypeField, opts)
	if err != nil {
		return err
	}

	if params.Ignored {
		return nil
	}

	if err := processField(refField, refTypeField, opts, params); err != nil {
		return err
	}

	if params.Init && isInvalidPtr(refField) {
		refField.Set(reflect.New(refField.Type().Elem()))
		refField = refField.Elem()
	}

	if refField.Kind() == reflect.Struct {
		return doParse(refField, processField, optionsWithEnvPrefix(refTypeField, opts))
	}

	if isSliceOfStructs(refTypeField) {
		return doParseSlice(refField, processField, optionsWithEnvPrefix(refTypeField, opts))
	}

	if isMapOfStructs(refTypeField) {
		return doParseMap(refField, processField, optionsWithEnvPrefix(refTypeField, opts), refTypeField)
	}

	return nil
}

func isSliceOfStructs(refTypeField reflect.StructField) bool {
	field := refTypeField.Type

	// *[]struct
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
		if field.Kind() == reflect.Slice && field.Elem().Kind() == reflect.Struct {
			return true
		}
	}

	// []struct{}
	if field.Kind() == reflect.Slice && field.Elem().Kind() == reflect.Struct {
		return true
	}

	return false
}

func doParseSlice(ref reflect.Value, processField processFieldFn, opts Options) error {
	if opts.Prefix != "" && !strings.HasSuffix(opts.Prefix, string(underscore)) {
		opts.Prefix += string(underscore)
	}

	var environments []string
	for environment := range opts.Environment {
		if strings.HasPrefix(environment, opts.Prefix) {
			environments = append(environments, environment)
		}
	}

	if len(environments) > 0 {
		counter := 0
		for finished := false; !finished; {
			finished = true
			prefix := fmt.Sprintf("%s%d%c", opts.Prefix, counter, underscore)
			for _, variable := range environments {
				if strings.HasPrefix(variable, prefix) {
					counter++
					finished = false
					break
				}
			}
		}

		sliceType := ref.Type()
		var initialized int
		if reflect.Ptr == ref.Kind() {
			sliceType = sliceType.Elem()
			// Due to the rest of code the pre-initialized slice has no chance for this situation
			initialized = 0
		} else {
			initialized = ref.Len()
		}

		var capacity int
		if capacity = initialized; counter > initialized {
			capacity = counter
		}
		result := reflect.MakeSlice(sliceType, capacity, capacity)
		for i := 0; i < capacity; i++ {
			item := result.Index(i)
			if i < initialized {
				item.Set(ref.Index(i))
			}
			if err := doParse(item, processField, optionsWithSliceEnvPrefix(opts, i)); err != nil {
				return err
			}
		}

		if result.Len() > 0 {
			if reflect.Ptr == ref.Kind() {
				resultPtr := reflect.New(sliceType)
				resultPtr.Elem().Set(result)
				result = resultPtr
			}
			ref.Set(result)
		}
	}

	return nil
}

func setField(refField reflect.Value, refTypeField reflect.StructField, opts Options, fieldParams FieldParams) error {
	value, err := get(fieldParams, opts)
	if err != nil {
		return err
	}

	if value != "" && (!opts.SetDefaultsForZeroValuesOnly || refField.IsZero()) {
		return set(refField, refTypeField, value, opts.FuncMap)
	}

	return nil
}

const underscore rune = '_'

func toEnvName(input string) string {
	var output []rune
	for i, c := range input {
		if c == underscore {
			continue
		}
		if len(output) > 0 && unicode.IsUpper(c) {
			if len(input) > i+1 {
				peek := rune(input[i+1])
				if unicode.IsLower(peek) || unicode.IsLower(rune(input[i-1])) {
					output = append(output, underscore)
				}
			}
		}
		output = append(output, unicode.ToUpper(c))
	}
	return string(output)
}

// FieldParams contains information about parsed field tags.
type FieldParams struct {
	OwnKey          string
	Key             string
	DefaultValue    string
	HasDefaultValue bool
	Required        bool
	LoadFile        bool
	Unset           bool
	NotEmpty        bool
	Expand          bool
	Init            bool
	Ignored         bool
}

func parseFieldParams(field reflect.StructField, opts Options) (FieldParams, error) {
	ownKey, tags := parseKeyForOption(field.Tag.Get(opts.TagName))
	if ownKey == "" && opts.UseFieldNameByDefault {
		ownKey = toEnvName(field.Name)
	}

	defaultValue, hasDefaultValue := field.Tag.Lookup(opts.DefaultValueTagName)

	result := FieldParams{
		OwnKey:          ownKey,
		Key:             opts.Prefix + ownKey,
		Required:        opts.RequiredIfNoDef,
		DefaultValue:    defaultValue,
		HasDefaultValue: hasDefaultValue,
		Ignored:         ownKey == "-",
	}

	for _, tag := range tags {
		switch tag {
		case "":
			continue
		case "file":
			result.LoadFile = true
		case "required":
			result.Required = true
		case "unset":
			result.Unset = true
		case "notEmpty":
			result.NotEmpty = true
		case "expand":
			result.Expand = true
		case "init":
			result.Init = true
		case "-":
			result.Ignored = true
		default:
			return FieldParams{}, newNoSupportedTagOptionError(tag)
		}
	}

	return result, nil
}

func get(fieldParams FieldParams, opts Options) (val string, err error) {
	var exists, isDefault bool

	val, exists, isDefault = getOr(
		fieldParams.Key,
		fieldParams.DefaultValue,
		fieldParams.HasDefaultValue,
		opts.Environment,
	)

	if fieldParams.Expand {
		val = os.Expand(val, opts.getRawEnv)
	}

	opts.rawEnvVars[fieldParams.OwnKey] = val

	if fieldParams.Unset {
		defer os.Unsetenv(fieldParams.Key)
	}

	if fieldParams.Required && !exists && fieldParams.OwnKey != "" {
		return "", newVarIsNotSetError(fieldParams.Key)
	}

	if fieldParams.NotEmpty && val == "" {
		return "", newEmptyVarError(fieldParams.Key)
	}

	if fieldParams.LoadFile && val != "" {
		filename := val
		val, err = getFromFile(filename)
		if err != nil {
			return "", newLoadFileContentError(filename, fieldParams.Key, err)
		}
	}

	if opts.OnSet != nil {
		if fieldParams.OwnKey != "" {
			opts.OnSet(fieldParams.Key, val, isDefault)
		}
	}
	return val, err
}

// split the env tag's key into the expected key and desired option, if any.
func parseKeyForOption(key string) (string, []string) {
	opts := strings.Split(key, ",")
	return opts[0], opts[1:]
}

func getFromFile(filename string) (value string, err error) {
	b, err := os.ReadFile(filename)
	return string(b), err
}

func getOr(key, defaultValue string, defExists bool, envs map[string]string) (val string, exists, isDefault bool) {
	value, exists := envs[key]
	switch {
	case (!exists || key == "") && defExists:
		return defaultValue, true, true
	case exists && value == "" && defExists:
		return defaultValue, true, true
	case !exists:
		return "", false, false
	}

	return value, true, false
}

func set(field reflect.Value, sf reflect.StructField, value string, funcMap map[reflect.Type]ParserFunc) error {
	if tm := asTextUnmarshaler(field); tm != nil {
		if err := tm.UnmarshalText([]byte(value)); err != nil {
			return newParseError(sf, err)
		}
		return nil
	}

	typee := sf.Type
	fieldee := field
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

	switch field.Kind() {
	case reflect.Slice:
		return handleSlice(field, value, sf, funcMap)
	case reflect.Map:
		return handleMap(field, value, sf, funcMap)
	}

	return newNoParserError(sf)
}

func handleSlice(field reflect.Value, value string, sf reflect.StructField, funcMap map[reflect.Type]ParserFunc) error {
	separator := sf.Tag.Get("envSeparator")
	if separator == "" {
		separator = ","
	}
	parts := strings.Split(value, separator)

	typee := sf.Type.Elem()
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

	result := reflect.MakeSlice(sf.Type, 0, len(parts))
	for _, part := range parts {
		r, err := parserFunc(part)
		if err != nil {
			return newParseError(sf, err)
		}
		v := reflect.ValueOf(r).Convert(typee)
		if sf.Type.Elem().Kind() == reflect.Ptr {
			v = reflect.New(typee)
			v.Elem().Set(reflect.ValueOf(r).Convert(typee))
		}
		result = reflect.Append(result, v)
	}
	field.Set(result)
	return nil
}

func handleMap(field reflect.Value, value string, sf reflect.StructField, funcMap map[reflect.Type]ParserFunc) error {
	keyType := sf.Type.Key()
	keyParserFunc, ok := funcMap[keyType]
	if !ok {
		keyParserFunc, ok = defaultBuiltInParsers[keyType.Kind()]
		if !ok {
			return newNoParserError(sf)
		}
	}

	elemType := sf.Type.Elem()
	elemParserFunc, ok := funcMap[elemType]
	if !ok {
		elemParserFunc, ok = defaultBuiltInParsers[elemType.Kind()]
		if !ok {
			return newNoParserError(sf)
		}
	}

	separator := sf.Tag.Get("envSeparator")
	if separator == "" {
		separator = ","
	}

	keyValSeparator := sf.Tag.Get("envKeyValSeparator")
	if keyValSeparator == "" {
		keyValSeparator = ":"
	}

	result := reflect.MakeMap(sf.Type)
	for _, part := range strings.Split(value, separator) {
		pairs := strings.SplitN(part, keyValSeparator, 2)
		if len(pairs) != 2 {
			return newParseError(sf, fmt.Errorf(`%q should be in "key%svalue" format`, part, keyValSeparator))
		}

		key, err := keyParserFunc(pairs[0])
		if err != nil {
			return newParseError(sf, err)
		}

		elem, err := elemParserFunc(pairs[1])
		if err != nil {
			return newParseError(sf, err)
		}

		result.SetMapIndex(reflect.ValueOf(key).Convert(keyType), reflect.ValueOf(elem).Convert(elemType))
	}

	field.Set(result)
	return nil
}

func asTextUnmarshaler(field reflect.Value) encoding.TextUnmarshaler {
	if field.Kind() == reflect.Ptr {
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

// ToMap Converts list of env vars as provided by os.Environ() to map you
// can use as Options.Environment field
func ToMap(env []string) map[string]string {
	return toMap(env)
}

func isInvalidPtr(v reflect.Value) bool {
	return reflect.Ptr == v.Kind() && v.Elem().Kind() == reflect.Invalid
}

func isMapOfStructs(refTypeField reflect.StructField) bool {
	field := refTypeField.Type

	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}

	if field.Kind() == reflect.Map {
		kind := field.Elem().Kind()
		if kind == reflect.Struct {
			return true
		}

		if kind == reflect.Ptr {
			ptrField := field.Elem()
			if ptrField.Elem().Kind() == reflect.Struct {
				return true
			}
		}
	}

	return false
}

func doParseMap(ref reflect.Value, processField processFieldFn, opts Options, sf reflect.StructField) error {
	if opts.Prefix != "" && !strings.HasSuffix(opts.Prefix, string(underscore)) {
		opts.Prefix += string(underscore)
	}

	var environments []string
	for environment := range opts.Environment {
		if strings.HasPrefix(environment, opts.Prefix) {
			environments = append(environments, environment)
		}
	}

	// There are no map keys that match
	if len(environments) == 0 {
		return nil
	}

	// Create a new map if it's nil
	if ref.IsNil() {
		ref.Set(reflect.MakeMap(ref.Type()))
	}

	// Get the key and value types
	keyType := ref.Type().Key()
	valueType := ref.Type().Elem()

	keyGroups := make(map[string]bool)

	structInnerSubEnvVars := getPossibleEnvVars(valueType, opts)

	for _, env := range environments {
		currKey := ""
		key := strings.TrimPrefix(env, opts.Prefix)

		// A user can have multiple environment variables which match to multiple keys
		// for example BAR_KEY_STR and BAR_KEY_NEW_STR are valid envars
		// If the struct has both "STR" and "NEW_STR" this would mean that
		// "STR" matches to both as a suffix and would result in two map keys
		// KEY_NEW and KEY, thus we match the suffix that would give the smallest key
		// since the smallest suffix that gives the largest key may have its own
		// different environment variable
		for _, innerEnvVar := range structInnerSubEnvVars {
			// If we are using a map of a map (we don't use the absolute path value, instead we use the prefix value)
			suffix := string(underscore) + innerEnvVar.value
			if innerEnvVar.useContains {
				idx := strings.LastIndex(key, suffix)
				if idx != -1 {
					newKey := key[:idx]
					// We had a better match which trimmed the key further
					if newKey != "" && (currKey == "" || len(currKey) > len(newKey)) {
						currKey = newKey
					}
				}
			} else if strings.HasSuffix(key, innerEnvVar.value) {
				if key == innerEnvVar.value {
					// If the key is exactly the innerEnvVar, this means that the env var was malformed
					return newParseError(sf, fmt.Errorf(`malformed complex map struct for %q`, key))
				}
				newKey := strings.TrimSuffix(key, suffix)
				// We had a better match which trimmed the key further
				if newKey != "" && (currKey == "" || len(currKey) > len(newKey)) {
					currKey = newKey
				}
			}
		}

		// If a key match has been found
		if currKey != "" {
			keyGroups[currKey] = true
		}
	}

	// Process each key group
	for mapKey := range keyGroups {
		value := reflect.New(valueType).Elem()
		keyOpts := optionsWithMapEnvPrefix(opts, mapKey)

		initialKind := value.Kind()
		if initialKind == reflect.Ptr {
			if value.IsNil() {
				value.Set(reflect.New(valueType.Elem()))
			}
			value = value.Elem()
		}

		err := doParse(value, processField, keyOpts)
		if err != nil {
			return err
		}

		parserFunc, ok := opts.FuncMap[keyType]
		if !ok {
			kind := keyType.Kind()
			if parserFunc, ok = defaultBuiltInParsers[kind]; !ok {
				return newNoParserError(sf)
			}
		}

		parsedKey, err := parserFunc(mapKey)
		if err != nil {
			return newParseError(sf, fmt.Errorf("failed to parse map key %q: %w", mapKey, err))
		}

		keyValue := reflect.ValueOf(parsedKey).Convert(keyType)

		if initialKind == reflect.Ptr {
			valuePtr := reflect.New(valueType.Elem())
			valuePtr.Elem().Set(value)
			value = valuePtr
		}

		ref.SetMapIndex(keyValue, value)
	}

	return nil
}

type SuffixType struct {
	useContains bool
	value       string
}

// getPossibleEnvVars returns all possible environment variables that could be set for a given struct type.
func getPossibleEnvVars(v reflect.Type, opts Options) []SuffixType {
	envVars := make(map[string]bool)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// The lib does not allow recursive structs technically
	// Recursive structs need to have the parent reference type as Pointer,
	// which means since pointer struct types do not get initialized by the parser by default,
	// and only with the `env:,init` tag. However, when the `init` attribute is set
	// the lib goes into an infinite loop because it does not support recursive structs
	// Thus we do not handle recursive structs here
	traverseStruct(v, "", opts, envVars)

	// Convert map keys to slice and sort for deterministic order
	result := make([]SuffixType, 0, len(envVars))
	for k, val := range envVars {
		entry := SuffixType{
			value:       k,
			useContains: val,
		}
		result = append(result, entry)
	}

	slices.SortFunc(result, func(i, j SuffixType) int {
		if i.useContains != j.useContains {
			if i.useContains {
				return 1
			}
			return -1
		}
		return strings.Compare(i.value, j.value)
	})

	return result
}

// traverseStruct recursively traverses a struct type and collects all possible environment variables.
func traverseStruct(t reflect.Type, prefix string, opts Options, envVars map[string]bool) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get field prefix if exists
		fieldPrefix := field.Tag.Get(opts.PrefixTagName)
		if fieldPrefix != "" {
			prefix = prefix + fieldPrefix
		}

		// Get env tag if exists
		envTag := field.Tag.Get(opts.TagName)
		key := prefix + envTag
		if envTag != "" {
			envVars[key] = false
		}

		// Handle nested structs and maps of structs
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct {
			traverseStruct(fieldType, prefix, opts, envVars)
		}

		if fieldType.Kind() == reflect.Map {
			elemType := fieldType.Elem()
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			if elemType.Kind() == reflect.Struct {
				envVars[key] = true
			}
		}
	}
}
