package env

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type unmarshaler struct {
	time.Duration
}

// TextUnmarshaler implements encoding.TextUnmarshaler
func (d *unmarshaler) UnmarshalText(data []byte) (err error) {
	if len(data) != 0 {
		d.Duration, err = time.ParseDuration(string(data))
	} else {
		d.Duration = 0
	}
	return err
}

// nolint: maligned
type Config struct {
	String          string          `env:"STRING"`
	StringPtr       *string         `env:"STRING_PTR"`
	Strings         []string        `env:"STRINGS"`
	Bool            bool            `env:"BOOL"`
	BoolPtr         *bool           `env:"BOOL_PTR"`
	Bools           []bool          `env:"BOOLS"`
	Int             int             `env:"INT"`
	IntPtr          *int            `env:"INT_PTR"`
	Ints            []int           `env:"INTS"`
	Int8            int8            `env:"INT8"`
	Int8Ptr         *int8           `env:"INT8_PTR"`
	Int8s           []int8          `env:"INT8S"`
	Int16           int16           `env:"INT16"`
	Int16s          []int16         `env:"INT16S"`
	Int16Ptr        *int16          `env:"INT16_PTR"`
	Int32           int32           `env:"INT32"`
	Int32s          []int32         `env:"INT32S"`
	Int32Ptr        *int32          `env:"INT32_PTR"`
	Int64           int64           `env:"INT64"`
	Int64s          []int64         `env:"INT64S"`
	Int64Ptr        *int64          `env:"INT64_PTR"`
	Uint            uint            `env:"UINT"`
	Uints           []uint          `env:"UINTS"`
	UintPtr         *uint           `env:"UINT_PTR"`
	Uint8           uint8           `env:"UINT8"`
	Uint8s          []uint8         `env:"UINT8S"`
	Uint8Ptr        *uint8          `env:"UINT8_PTR"`
	Uint16          uint16          `env:"UINT16"`
	Uint16s         []uint16        `env:"UINT16S"`
	Uint16Ptr       *uint16         `env:"UINT16_PTR"`
	Uint32          uint32          `env:"UINT32"`
	Uint32s         []uint32        `env:"UINT32S"`
	Uint32Ptr       *uint32         `env:"UINT32_PTR"`
	Uint64          uint64          `env:"UINT64"`
	Uint64s         []uint64        `env:"UINT64S"`
	Uint64Ptr       *uint64         `env:"UINT64_PTR"`
	Float32         float32         `env:"FLOAT32"`
	Float32Ptr      *float32        `env:"FLOAT32_PTR"`
	Float32s        []float32       `env:"FLOAT32S"`
	Float64         float64         `env:"FLOAT64"`
	Float64Ptr      *float64        `env:"FLOAT64_PTR"`
	Float64s        []float64       `env:"FLOAT64S"`
	DatabaseURL     string          `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/db"`
	SepStrings      []string        `env:"SEPSTRINGS" envSeparator:":"`
	Duration        time.Duration   `env:"DURATION"`
	Durations       []time.Duration `env:"DURATIONS"`
	DurationPtr     *time.Duration  `env:"DURATION_PTR"`
	Unmarshaler     unmarshaler     `env:"UNMARSHALER"`
	UnmarshalerPtr  *unmarshaler    `env:"UNMARSHALER_PTR"`
	Unmarshalers    []unmarshaler   `env:"UNMARSHALERS"`
	UnmarshalerPtrs []*unmarshaler  `env:"UNMARSHALER_PTRS"`
	URL             url.URL         `env:"URL"`
	URLPtr          *url.URL        `env:"URL_PTR"`
	URLs            []url.URL       `env:"URLS"`
	URLPtrs         []*url.URL      `env:"URL_PTRS"`

	NotAnEnv   string
	unexported string `env:"FOO"`
}

type ParentStruct struct {
	InnerStruct *InnerStruct
	unexported  *InnerStruct
	Ignored     *http.Client
}

type InnerStruct struct {
	Inner  string `env:"innervar"`
	Number uint   `env:"innernum"`
}

type ForNestedStruct struct {
	NestedStruct
}
type NestedStruct struct {
	NestedVar string `env:"nestedvar"`
}

func TestParsesEnv(t *testing.T) {
	os.Setenv("STRING", "somevalue")
	os.Setenv("BOOL", "true")
	os.Setenv("INT", "8080")
	os.Setenv("STRINGS", "string1,string2,string3")
	os.Setenv("SEPSTRINGS", "string1:string2:string3")
	os.Setenv("INTS", "1,2,3,4")
	os.Setenv("INT64S", "1,2,2147483640,-2147483640")
	os.Setenv("UINT64S", "1,2,214748364011,9147483641")
	os.Setenv("BOOLS", "t,TRUE,0,1")
	os.Setenv("FLOAT32", "3.40282346638528859811704183484516925440e+38")
	os.Setenv("FLOAT64", "1.797693134862315708145274237317043567981e+308")
	os.Setenv("FLOAT32S", "1.0,2.0,3.0")
	os.Setenv("FLOAT64S", "1.0,2.0,3.0")
	os.Setenv("UINT", "44")
	os.Setenv("UINT8", "88")
	os.Setenv("UINT16", "1616")
	os.Setenv("UINT32", "3232")
	os.Setenv("UINT64", "6464")
	os.Setenv("INT8", "-88")
	os.Setenv("INT16", "-1616")
	os.Setenv("INT32", "-3232")
	os.Setenv("INT64", "-7575")
	os.Setenv("DURATION", "1s")
	os.Setenv("DURATION_PTR", "1s")
	os.Setenv("DURATIONS", "1s,2s,3s")
	os.Setenv("UNMARSHALER", "1s")
	os.Setenv("UNMARSHALER_PTR", "1m")
	os.Setenv("UNMARSHALERS", "2m,3m")
	os.Setenv("UNMARSHALER_PTRS", "2m,3m")
	os.Setenv("URL", "https://carlosbecker.dev")
	os.Setenv("URL_PTR", "https://carlosbecker.dev/ptr")
	os.Setenv("URLS", "https://carlosbecker.dev,https://carlosbecker.com")

	defer os.Clearenv()

	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "somevalue", cfg.String)
	assert.Equal(t, true, cfg.Bool)
	assert.Equal(t, 8080, cfg.Int)
	assert.Equal(t, uint(44), cfg.Uint)
	assert.Equal(t, uint8(88), cfg.Uint8)
	assert.Equal(t, uint16(1616), cfg.Uint16)
	assert.Equal(t, uint32(3232), cfg.Uint32)
	assert.Equal(t, uint64(6464), cfg.Uint64)
	assert.Equal(t, int8(-88), cfg.Int8)
	assert.Equal(t, int16(-1616), cfg.Int16)
	assert.Equal(t, int32(-3232), cfg.Int32)
	assert.Equal(t, int64(-7575), cfg.Int64)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.Strings)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.SepStrings)
	assert.Equal(t, []int{1, 2, 3, 4}, cfg.Ints)
	assert.Equal(t, []int64{1, 2, 2147483640, -2147483640}, cfg.Int64s)
	assert.Equal(t, []uint64{1, 2, 214748364011, 9147483641}, cfg.Uint64s)
	assert.Equal(t, []bool{true, true, false, true}, cfg.Bools)
	d1, _ := time.ParseDuration("1s")
	assert.Equal(t, d1, cfg.Duration)
	assert.Equal(t, &d1, cfg.DurationPtr)
	f32 := float32(3.40282346638528859811704183484516925440e+38)
	assert.Equal(t, f32, cfg.Float32)
	f64 := float64(1.797693134862315708145274237317043567981e+308)
	assert.Equal(t, f64, cfg.Float64)
	assert.Equal(t, []float32{float32(1.0), float32(2.0), float32(3.0)}, cfg.Float32s)
	assert.Equal(t, []float64{float64(1.0), float64(2.0), float64(3.0)}, cfg.Float64s)
	d2, _ := time.ParseDuration("2s")
	d3, _ := time.ParseDuration("3s")
	assert.Equal(t, []time.Duration{d1, d2, d3}, cfg.Durations)
	assert.Equal(t, time.Second, cfg.Unmarshaler.Duration)
	assert.Equal(t, time.Minute, cfg.UnmarshalerPtr.Duration)
	assert.Equal(t, []unmarshaler{{time.Minute * 2}, {time.Minute * 3}}, cfg.Unmarshalers)
	assert.Equal(t, []*unmarshaler{{time.Minute * 2}, {time.Minute * 3}}, cfg.UnmarshalerPtrs)

	assert.Equal(t, "https://carlosbecker.dev", cfg.URL.String())
	assert.Equal(t, "https://carlosbecker.dev/ptr", cfg.URLPtr.String())
	assert.Equal(t, "https://carlosbecker.dev", cfg.URLs[0].String())
	assert.Equal(t, "https://carlosbecker.com", cfg.URLs[1].String())
}

func TestParsesEnvInner(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
		unexported:  &InnerStruct{},
	}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "someinnervalue", cfg.InnerStruct.Inner)
}

func TestParsesEnvInnerNil(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{}
	assert.NoError(t, Parse(&cfg))
}

func TestParsesEnvInnerInvalid(t *testing.T) {
	os.Setenv("innernum", "-547")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
	}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Number\" of type \"uint\": strconv.ParseUint: parsing \"-547\": invalid syntax")
}

func TestParsesEnvNested(t *testing.T) {
	os.Setenv("nestedvar", "somenestedvalue")
	defer os.Clearenv()
	var cfg ForNestedStruct
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "somenestedvalue", cfg.NestedVar)
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "", cfg.String)
	assert.Equal(t, false, cfg.Bool)
	assert.Equal(t, 0, cfg.Int)
	assert.Equal(t, uint(0), cfg.Uint)
	assert.Equal(t, uint64(0), cfg.Uint64)
	assert.Equal(t, int64(0), cfg.Int64)
	assert.Equal(t, 0, len(cfg.Strings))
	assert.Equal(t, 0, len(cfg.SepStrings))
	assert.Equal(t, 0, len(cfg.Ints))
	assert.Equal(t, 0, len(cfg.Bools))
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	assert.EqualError(t, Parse(&thisShouldBreak), "env: expected a pointer to a Struct")
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	assert.EqualError(t, Parse(cfg), "env: expected a pointer to a Struct")
}

func TestInvalidBool(t *testing.T) {
	os.Setenv("BOOL", "should-be-a-bool")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Bool\" of type \"bool\": strconv.ParseBool: parsing \"should-be-a-bool\": invalid syntax")
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("INT", "should-be-an-int")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Int\" of type \"int\": strconv.ParseInt: parsing \"should-be-an-int\": invalid syntax")
}

func TestInvalidUint(t *testing.T) {
	os.Setenv("UINT", "-44")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Uint\" of type \"uint\": strconv.ParseUint: parsing \"-44\": invalid syntax")
}

func TestInvalidFloat32(t *testing.T) {
	os.Setenv("FLOAT32", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Float32\" of type \"float32\": strconv.ParseFloat: parsing \"AAA\": invalid syntax")
}

func TestInvalidFloat64(t *testing.T) {
	os.Setenv("FLOAT64", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Float64\" of type \"float64\": strconv.ParseFloat: parsing \"AAA\": invalid syntax")
}

func TestInvalidUint64(t *testing.T) {
	os.Setenv("UINT64", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Uint64\" of type \"uint64\": strconv.ParseUint: parsing \"AAA\": invalid syntax")
}

func TestInvalidInt64(t *testing.T) {
	os.Setenv("INT64", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Int64\" of type \"int64\": strconv.ParseInt: parsing \"AAA\": invalid syntax")
}

func TestInvalidInt64Slice(t *testing.T) {
	type config struct {
		BadFloats []int64 `env:"BADINTS"`
	}

	os.Setenv("BADINTS", "A,2,3")
	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: parse error on field \"BadFloats\" of type \"[]int64\": strconv.ParseInt: parsing \"A\": invalid syntax")
}

func TestInvalidUInt64Slice(t *testing.T) {
	type config struct {
		BadFloats []uint64 `env:"BADINTS"`
	}

	os.Setenv("BADFLOATS", "A,2,3")
	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: parse error on field \"BadFloats\" of type \"[]uint64\": strconv.ParseUint: parsing \"A\": invalid syntax")
}

func TestInvalidFloat32Slice(t *testing.T) {
	type config struct {
		BadFloats []float32 `env:"BADFLOATS"`
	}

	os.Setenv("BADFLOATS", "A,2.0,3.0")
	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: parse error on field \"BadFloats\" of type \"[]float32\": strconv.ParseFloat: parsing \"A\": invalid syntax")
}

func TestInvalidFloat64Slice(t *testing.T) {
	type config struct {
		BadFloats []float64 `env:"BADFLOATS"`
	}

	os.Setenv("BADFLOATS", "A,2.0,3.0")
	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: parse error on field \"BadFloats\" of type \"[]float64\": strconv.ParseFloat: parsing \"A\": invalid syntax")
}

func TestInvalidBoolsSlice(t *testing.T) {
	type config struct {
		BadBools []bool `env:"BADBOOLS"`
	}

	os.Setenv("BADBOOLS", "t,f,TRUE,faaaalse")
	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: parse error on field \"BadBools\" of type \"[]bool\": strconv.ParseBool: parsing \"faaaalse\": invalid syntax")
}

func TestInvalidDuration(t *testing.T) {
	os.Setenv("DURATION", "should-be-a-valid-duration")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Duration\" of type \"time.Duration\": unable to parser duration: time: invalid duration should-be-a-valid-duration")
}

func TestInvalidDurations(t *testing.T) {
	os.Setenv("DURATIONS", "1s,contains-an-invalid-duration,3s")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"Durations\" of type \"[]time.Duration\": unable to parser duration: time: invalid duration contains-an-invalid-duration")
}

func TestParsesDefaultConfig(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "postgres://localhost:5432/db", cfg.DatabaseURL)
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Empty(t, cfg.NotAnEnv)
}

func TestParseStructWithInvalidFieldKind(t *testing.T) {
	type config struct {
		WontWorkByte byte `env:"BLAH"`
	}
	os.Setenv("BLAH", "a")
	cfg := config{}
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"WontWorkByte\" of type \"uint8\": strconv.ParseUint: parsing \"a\": invalid syntax")
}

func TestUnsupportedSliceType(t *testing.T) {
	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}

	os.Setenv("WONTWORK", "1,2,3")
	defer os.Clearenv()

	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: no parser found for field \"WontWork\" of type \"[]map[int]int\"")
}

func TestBadSeparator(t *testing.T) {
	type config struct {
		WontWork []int `env:"WONTWORK" envSeparator:":"`
	}

	cfg := &config{}
	os.Setenv("WONTWORK", "1,2,3,4")
	defer os.Clearenv()

	assert.EqualError(t, Parse(cfg), "env: parse error on field \"WontWork\" of type \"[]int\": strconv.ParseInt: parsing \"1,2,3,4\": invalid syntax")
}

func TestNoErrorRequiredSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}

	cfg := &config{}

	os.Setenv("IS_REQUIRED", "")
	defer os.Clearenv()
	assert.NoError(t, Parse(cfg))
	assert.Equal(t, "", cfg.IsRequired)
}

func TestErrorRequiredNotSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}

	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: required environment variable \"IS_REQUIRED\" is not set")
}

func TestParseExpandOption(t *testing.T) {
	type config struct {
		Host        string `env:"HOST" envDefault:"localhost"`
		Port        int    `env:"PORT" envDefault:"3000" envExpand:"True"`
		SecretKey   string `env:"SECRET_KEY" envExpand:"True"`
		ExpandKey   string `env:"EXPAND_KEY"`
		CompoundKey string `env:"HOST_PORT" envDefault:"${HOST}:${PORT}" envExpand:"True"`
		Default     string `env:"DEFAULT" envDefault:"def1"  envExpand:"True"`
	}
	defer os.Clearenv()

	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "3000")
	os.Setenv("EXPAND_KEY", "qwerty12345")
	os.Setenv("SECRET_KEY", "${EXPAND_KEY}")

	cfg := config{}
	err := Parse(&cfg)

	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 3000, cfg.Port)
	assert.Equal(t, "qwerty12345", cfg.SecretKey)
	assert.Equal(t, "qwerty12345", cfg.ExpandKey)
	assert.Equal(t, "localhost:3000", cfg.CompoundKey)
	assert.Equal(t, "def1", cfg.Default)
}

func TestCustomParser(t *testing.T) {
	type foo struct {
		name string
	}

	type bar struct {
		Name string `env:"OTHER"`
		Foo  *foo   `env:"BLAH"`
	}

	type config struct {
		Var   foo  `env:"VAR"`
		Foo   *foo `env:"BLAH"`
		Other *bar
	}

	os.Setenv("VAR", "test")
	defer os.Unsetenv("VAR")
	os.Setenv("OTHER", "test2")
	defer os.Unsetenv("OTHER")
	os.Setenv("BLAH", "test3")
	defer os.Unsetenv("BLAH")

	cfg := &config{
		Other: &bar{},
	}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(foo{}): func(v string) (interface{}, error) {
			return foo{name: v}, nil
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, cfg.Var.name, "test")
	assert.Equal(t, cfg.Foo.name, "test3")
	assert.Equal(t, cfg.Other.Name, "test2")
	assert.Equal(t, cfg.Other.Foo.name, "test3")
}

func TestParseWithFuncsNoPtr(t *testing.T) {
	type foo struct{}
	err := ParseWithFuncs(foo{}, nil)
	assert.EqualError(t, err, "env: expected a pointer to a Struct")
}

func TestParseWithFuncsInvalidType(t *testing.T) {
	var c int
	err := ParseWithFuncs(&c, nil)
	assert.EqualError(t, err, "env: expected a pointer to a Struct")
}

func TestCustomParserError(t *testing.T) {
	type foo struct {
		name string
	}

	customParserFunc := func(v string) (interface{}, error) {
		return nil, errors.New("something broke")
	}

	t.Run("single", func(t *testing.T) {
		type config struct {
			Var foo `env:"VAR"`
		}

		os.Setenv("VAR", "single")
		cfg := &config{}
		err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
			reflect.TypeOf(foo{}): customParserFunc,
		})

		assert.Empty(t, cfg.Var.name)
		assert.EqualError(t, err, "env: parse error on field \"Var\" of type \"env.foo\": something broke")
	})

	t.Run("slice", func(t *testing.T) {
		type config struct {
			Var []foo `env:"VAR2"`
		}
		os.Setenv("VAR2", "slice,slace")

		cfg := &config{}
		err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
			reflect.TypeOf(foo{}): customParserFunc,
		})

		assert.Empty(t, cfg.Var)
		assert.EqualError(t, err, "env: parse error on field \"Var\" of type \"[]env.foo\": something broke")
	})
}

func TestCustomParserBasicType(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_"`
	}

	exp := ConstT(123)
	os.Setenv("CONST_", fmt.Sprintf("%d", exp))

	customParserFunc := func(v string) (interface{}, error) {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		r := ConstT(i)
		return r, nil
	}

	cfg := &config{}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	assert.NoError(t, err)
	assert.Equal(t, exp, cfg.Const)
}

func TestCustomParserUint64Alias(t *testing.T) {
	type T uint64

	var one T = 1

	type config struct {
		Val T `env:"" envDefault:"1x"`
	}

	parserCalled := false

	tParser := func(value string) (interface{}, error) {
		parserCalled = true
		trimmed := strings.TrimSuffix(value, "x")
		i, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, err
		}
		return T(i), nil
	}

	cfg := config{}

	err := ParseWithFuncs(&cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(one): tParser,
	})

	assert.True(t, parserCalled, "tParser should have been called")
	assert.NoError(t, err)
	assert.Equal(t, T(1), cfg.Val)
}

func TestTypeCustomParserBasicInvalid(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_"`
	}

	os.Setenv("CONST_", "foobar")

	customParserFunc := func(_ string) (interface{}, error) {
		return nil, errors.New("random error")
	}

	cfg := &config{}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	assert.Empty(t, cfg.Const)
	assert.EqualError(t, err, "env: parse error on field \"Const\" of type \"env.ConstT\": random error")
}

func TestCustomParserNotCalledForNonAlias(t *testing.T) {
	type T uint64
	type U uint64

	type config struct {
		Val   uint64 `env:"" envDefault:"33"`
		Other U      `env:"OTHER" envDefault:"44"`
	}

	tParserCalled := false

	tParser := func(value string) (interface{}, error) {
		tParserCalled = true
		return T(99), nil
	}

	cfg := config{}

	err := ParseWithFuncs(&cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(T(0)): tParser,
	})

	assert.False(t, tParserCalled, "tParser should not have been called")
	assert.NoError(t, err)
	assert.Equal(t, uint64(33), cfg.Val)
	assert.Equal(t, U(44), cfg.Other)
}

func TestCustomParserBasicUnsupported(t *testing.T) {
	type ConstT struct {
		A int
	}

	type config struct {
		Const ConstT `env:"CONST_"`
	}

	os.Setenv("CONST_", "42")

	cfg := &config{}
	err := Parse(cfg)

	assert.Zero(t, cfg.Const)
	assert.EqualError(t, err, "env: no parser found for field \"Const\" of type \"env.ConstT\"")
}

func TestUnsupportedStructType(t *testing.T) {
	type config struct {
		Foo http.Client `env:"FOO"`
	}

	os.Setenv("FOO", "foo")

	cfg := &config{}
	err := Parse(cfg)

	assert.EqualError(t, err, "env: no parser found for field \"Foo\" of type \"http.Client\"")
}

func TestEmptyOption(t *testing.T) {
	type config struct {
		Var string `env:"VAR,"`
	}

	cfg := &config{}

	os.Setenv("VAR", "")
	defer os.Clearenv()
	assert.NoError(t, Parse(cfg))
	assert.Equal(t, "", cfg.Var)
}

func TestErrorOptionNotRecognized(t *testing.T) {
	type config struct {
		Var string `env:"VAR,not_supported!"`
	}

	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: tag option \"not_supported!\" not supported")
}

func TestTextUnmarshalerError(t *testing.T) {
	type config struct {
		Unmarshaler unmarshaler `env:"UNMARSHALER"`
	}
	os.Setenv("UNMARSHALER", "invalid")
	cfg := &config{}
	assert.EqualError(t, Parse(cfg), "env: parse error on field \"Unmarshaler\" of type \"env.unmarshaler\": time: invalid duration invalid")
}

func TestParseURL(t *testing.T) {
	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL" envDefault:"https://google.com"`
	}
	var cfg config
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "https://google.com", cfg.ExampleURL.String())
}

func TestParseInvalidURL(t *testing.T) {
	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL_2"`
	}
	var cfg config
	os.Setenv("EXAMPLE_URL_2", "nope://s s/")
	assert.EqualError(t, Parse(&cfg), "env: parse error on field \"ExampleURL\" of type \"url.URL\": unable parse URL: parse nope://s s/: invalid character \" \" in host name")
}

func ExampleParse() {
	type inner struct {
		Foo string `env:"FOO" envDefault:"foobar"`
	}
	type config struct {
		Home         string `env:"HOME,required"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
		Inner        inner
	}
	os.Setenv("HOME", "/tmp/fakehome")
	var cfg config
	if err := Parse(&cfg); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Home:/tmp/fakehome Port:3000 IsProduction:false Inner:{Foo:foobar}}
}

func TestIgnoresUnexported(t *testing.T) {
	type unexportedConfig struct {
		home  string `env:"HOME"`
		Home2 string `env:"HOME"`
	}
	cfg := unexportedConfig{}

	os.Setenv("HOME", "/tmp/fakehome")
	assert.NoError(t, Parse(&cfg))
	assert.Empty(t, cfg.home)
	assert.Equal(t, "/tmp/fakehome", cfg.Home2)
}

type LogLevel int8

func (l *LogLevel) UnmarshalText(text []byte) error {
	txt := string(text)
	switch txt {
	case "debug":
		*l = DebugLevel
	case "info":
		*l = InfoLevel
	default:
		return fmt.Errorf("unknown level: %q", txt)
	}

	return nil
}

const (
	DebugLevel LogLevel = iota - 1
	InfoLevel
)

func TestPrecedenceUnmarshalText(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_LEVELS", "debug,info")
	defer os.Unsetenv("LOG_LEVEL")
	defer os.Unsetenv("LOG_LEVELS")

	type config struct {
		LogLevel  LogLevel   `env:"LOG_LEVEL"`
		LogLevels []LogLevel `env:"LOG_LEVELS"`
	}
	var cfg config

	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, DebugLevel, cfg.LogLevel)
	assert.Equal(t, []LogLevel{DebugLevel, InfoLevel}, cfg.LogLevels)
}

func ExampleParseWithFuncs() {
	type thing struct {
		desc string
	}

	type conf struct {
		Thing thing `env:"THING"`
	}

	os.Setenv("THING", "my thing")

	var c = conf{}

	err := ParseWithFuncs(&c, map[reflect.Type]ParserFunc{
		reflect.TypeOf(thing{}): func(v string) (interface{}, error) {
			return thing{desc: v}, nil
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.Thing.desc)
	// Output:
	// my thing
}
