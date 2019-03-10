package env_test

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

	"github.com/caarlos0/env"
	"github.com/caarlos0/env/parsers"
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
	Some            string `env:"somevar"`
	Other           bool   `env:"othervar"`
	Port            int    `env:"PORT"`
	Int8Val         int8   `env:"INT8VAL"`
	Int16Val        int16  `env:"INT16VAL"`
	Int32Val        int32  `env:"INT32VAL"`
	Int64Val        int64  `env:"INT64VAL"`
	UintVal         uint   `env:"UINTVAL"`
	Uint8Val        uint8  `env:"UINT8VAL"`
	Uint16Val       uint16 `env:"UINT16VAL"`
	Uint32Val       uint32 `env:"UINT32VAL"`
	Uint64Val       uint64 `env:"UINT64VAL"`
	NotAnEnv        string
	DatabaseURL     string          `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/db"`
	Strings         []string        `env:"STRINGS"`
	SepStrings      []string        `env:"SEPSTRINGS" envSeparator:":"`
	Numbers         []int           `env:"NUMBERS"`
	Numbers64       []int64         `env:"NUMBERS64"`
	UNumbers64      []uint64        `env:"UNUMBERS64"`
	Bools           []bool          `env:"BOOLS"`
	Duration        time.Duration   `env:"DURATION"`
	Float32         float32         `env:"FLOAT32"`
	Float64         float64         `env:"FLOAT64"`
	Float32s        []float32       `env:"FLOAT32S"`
	Float64s        []float64       `env:"FLOAT64S"`
	Durations       []time.Duration `env:"DURATIONS"`
	Unmarshaler     unmarshaler     `env:"UNMARSHALER"`
	UnmarshalerPtr  *unmarshaler    `env:"UNMARSHALER_PTR"`
	Unmarshalers    []unmarshaler   `env:"UNMARSHALERS"`
	UnmarshalerPtrs []*unmarshaler  `env:"UNMARSHALER_PTRS"`
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
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	os.Setenv("STRINGS", "string1,string2,string3")
	os.Setenv("SEPSTRINGS", "string1:string2:string3")
	os.Setenv("NUMBERS", "1,2,3,4")
	os.Setenv("NUMBERS64", "1,2,2147483640,-2147483640")
	os.Setenv("UNUMBERS64", "1,2,214748364011,9147483641")
	os.Setenv("BOOLS", "t,TRUE,0,1")
	os.Setenv("DURATION", "1s")
	os.Setenv("FLOAT32", "3.40282346638528859811704183484516925440e+38")
	os.Setenv("FLOAT64", "1.797693134862315708145274237317043567981e+308")
	os.Setenv("FLOAT32S", "1.0,2.0,3.0")
	os.Setenv("FLOAT64S", "1.0,2.0,3.0")
	os.Setenv("UINTVAL", "44")
	os.Setenv("UINT8VAL", "88")
	os.Setenv("UINT16VAL", "1616")
	os.Setenv("UINT32VAL", "3232")
	os.Setenv("UINT64VAL", "6464")
	os.Setenv("INT8VAL", "-88")
	os.Setenv("INT16VAL", "-1616")
	os.Setenv("INT32VAL", "-3232")
	os.Setenv("INT64VAL", "-7575")
	os.Setenv("DURATIONS", "1s,2s,3s")
	os.Setenv("UNMARSHALER", "1s")
	os.Setenv("UNMARSHALER_PTR", "1m")
	os.Setenv("UNMARSHALERS", "2m,3m")
	os.Setenv("UNMARSHALER_PTRS", "2m,3m")

	defer os.Clearenv()

	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "somevalue", cfg.Some)
	assert.Equal(t, true, cfg.Other)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, uint(44), cfg.UintVal)
	assert.Equal(t, uint8(88), cfg.Uint8Val)
	assert.Equal(t, uint16(1616), cfg.Uint16Val)
	assert.Equal(t, uint32(3232), cfg.Uint32Val)
	assert.Equal(t, uint64(6464), cfg.Uint64Val)
	assert.Equal(t, int8(-88), cfg.Int8Val)
	assert.Equal(t, int16(-1616), cfg.Int16Val)
	assert.Equal(t, int32(-3232), cfg.Int32Val)
	assert.Equal(t, int64(-7575), cfg.Int64Val)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.Strings)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.SepStrings)
	assert.Equal(t, []int{1, 2, 3, 4}, cfg.Numbers)
	assert.Equal(t, []int64{1, 2, 2147483640, -2147483640}, cfg.Numbers64)
	assert.Equal(t, []uint64{1, 2, 214748364011, 9147483641}, cfg.UNumbers64)
	assert.Equal(t, []bool{true, true, false, true}, cfg.Bools)
	d1, _ := time.ParseDuration("1s")
	assert.Equal(t, d1, cfg.Duration)
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
}

func TestParsesEnvInner(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
		unexported:  &InnerStruct{},
	}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "someinnervalue", cfg.InnerStruct.Inner)
}

func TestParsesEnvInnerNil(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{}
	assert.NoError(t, env.Parse(&cfg))
}

func TestParsesEnvInnerInvalid(t *testing.T) {
	os.Setenv("innernum", "-547")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
	}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Number\" of type \"uint\": strconv.ParseUint: parsing \"-547\": invalid syntax")
}

func TestParsesEnvNested(t *testing.T) {
	os.Setenv("nestedvar", "somenestedvalue")
	defer os.Clearenv()
	var cfg ForNestedStruct
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "somenestedvalue", cfg.NestedVar)
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "", cfg.Some)
	assert.Equal(t, false, cfg.Other)
	assert.Equal(t, 0, cfg.Port)
	assert.Equal(t, uint(0), cfg.UintVal)
	assert.Equal(t, uint64(0), cfg.Uint64Val)
	assert.Equal(t, int64(0), cfg.Int64Val)
	assert.Equal(t, 0, len(cfg.Strings))
	assert.Equal(t, 0, len(cfg.SepStrings))
	assert.Equal(t, 0, len(cfg.Numbers))
	assert.Equal(t, 0, len(cfg.Bools))
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	assert.EqualError(t, env.Parse(&thisShouldBreak), "env: expected a pointer to a Struct")
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	assert.EqualError(t, env.Parse(cfg), "env: expected a pointer to a Struct")
}

func TestInvalidBool(t *testing.T) {
	os.Setenv("othervar", "should-be-a-bool")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Other\" of type \"bool\": strconv.ParseBool: parsing \"should-be-a-bool\": invalid syntax")
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("PORT", "should-be-an-int")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Port\" of type \"int\": strconv.ParseInt: parsing \"should-be-an-int\": invalid syntax")
}

func TestInvalidUint(t *testing.T) {
	os.Setenv("UINTVAL", "-44")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"UintVal\" of type \"uint\": strconv.ParseUint: parsing \"-44\": invalid syntax")
}

func TestInvalidFloat32(t *testing.T) {
	os.Setenv("FLOAT32", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Float32\" of type \"float32\": strconv.ParseFloat: parsing \"AAA\": invalid syntax")
}

func TestInvalidFloat64(t *testing.T) {
	os.Setenv("FLOAT64", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Float64\" of type \"float64\": strconv.ParseFloat: parsing \"AAA\": invalid syntax")
}

func TestInvalidUint64(t *testing.T) {
	os.Setenv("UINT64VAL", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Uint64Val\" of type \"uint64\": strconv.ParseUint: parsing \"AAA\": invalid syntax")
}

func TestInvalidInt64(t *testing.T) {
	os.Setenv("INT64VAL", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Int64Val\" of type \"int64\": strconv.ParseInt: parsing \"AAA\": invalid syntax")
}

func TestInvalidInt64Slice(t *testing.T) {
	type config struct {
		BadFloats []int64 `env:"BADINTS"`
	}

	os.Setenv("BADINTS", "A,2,3")
	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: parse error on field \"BadFloats\" of type \"[]int64\": strconv.ParseInt: parsing \"A\": invalid syntax")
}

func TestInvalidUInt64Slice(t *testing.T) {
	type config struct {
		BadFloats []uint64 `env:"BADINTS"`
	}

	os.Setenv("BADFLOATS", "A,2,3")
	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: parse error on field \"BadFloats\" of type \"[]uint64\": strconv.ParseUint: parsing \"A\": invalid syntax")
}

func TestInvalidFloat32Slice(t *testing.T) {
	type config struct {
		BadFloats []float32 `env:"BADFLOATS"`
	}

	os.Setenv("BADFLOATS", "A,2.0,3.0")
	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: parse error on field \"BadFloats\" of type \"[]float32\": strconv.ParseFloat: parsing \"A\": invalid syntax")
}

func TestInvalidFloat64Slice(t *testing.T) {
	type config struct {
		BadFloats []float64 `env:"BADFLOATS"`
	}

	os.Setenv("BADFLOATS", "A,2.0,3.0")
	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: parse error on field \"BadFloats\" of type \"[]float64\": strconv.ParseFloat: parsing \"A\": invalid syntax")
}

func TestInvalidBoolsSlice(t *testing.T) {
	type config struct {
		BadBools []bool `env:"BADBOOLS"`
	}

	os.Setenv("BADBOOLS", "t,f,TRUE,faaaalse")
	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: parse error on field \"BadBools\" of type \"[]bool\": strconv.ParseBool: parsing \"faaaalse\": invalid syntax")
}

func TestInvalidDuration(t *testing.T) {
	os.Setenv("DURATION", "should-be-a-valid-duration")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Duration\" of type \"time.Duration\": unable to parser duration: time: invalid duration should-be-a-valid-duration")
}

func TestInvalidDurations(t *testing.T) {
	os.Setenv("DURATIONS", "1s,contains-an-invalid-duration,3s")
	defer os.Clearenv()

	cfg := Config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"Durations\" of type \"[]time.Duration\": unable to parser duration: time: invalid duration contains-an-invalid-duration")
}

func TestParsesDefaultConfig(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "postgres://localhost:5432/db", cfg.DatabaseURL)
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Empty(t, cfg.NotAnEnv)
}

func TestParseStructWithInvalidFieldKind(t *testing.T) {
	type config struct {
		WontWorkByte byte `env:"BLAH"`
	}
	os.Setenv("BLAH", "a")
	cfg := config{}
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"WontWorkByte\" of type \"uint8\": strconv.ParseUint: parsing \"a\": invalid syntax")
}

func TestUnsupportedSliceType(t *testing.T) {
	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}

	os.Setenv("WONTWORK", "1,2,3")
	defer os.Clearenv()

	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: no parser found for field \"WontWork\" of type \"[]map[int]int\"")
}

func TestBadSeparator(t *testing.T) {
	type config struct {
		WontWork []int `env:"WONTWORK" envSeparator:":"`
	}

	cfg := &config{}
	os.Setenv("WONTWORK", "1,2,3,4")
	defer os.Clearenv()

	assert.EqualError(t, env.Parse(cfg), "env: parse error on field \"WontWork\" of type \"[]int\": strconv.ParseInt: parsing \"1,2,3,4\": invalid syntax")
}

func TestNoErrorRequiredSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}

	cfg := &config{}

	os.Setenv("IS_REQUIRED", "val")
	defer os.Clearenv()
	assert.NoError(t, env.Parse(cfg))
	assert.Equal(t, "val", cfg.IsRequired)
}

func TestErrorRequiredNotSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}

	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: required environment variable \"\"IS_REQUIRED\"\" is not set")
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
	err := env.Parse(&cfg)

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

	type config struct {
		Var foo `env:"VAR"`
	}

	os.Setenv("VAR", "test")

	customParserFunc := func(v string) (interface{}, error) {
		return foo{name: v}, nil
	}

	cfg := &config{}
	err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(foo{}): customParserFunc,
	})

	assert.NoError(t, err)
	assert.Equal(t, cfg.Var.name, "test")
}

func TestParseWithFuncsNoPtr(t *testing.T) {
	type foo struct{}
	err := env.ParseWithFuncs(foo{}, nil)
	assert.EqualError(t, err, "env: expected a pointer to a Struct")
}

func TestParseWithFuncsInvalidType(t *testing.T) {
	var c int
	err := env.ParseWithFuncs(&c, nil)
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
		err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(foo{}): customParserFunc,
		})

		assert.Empty(t, cfg.Var.name)
		assert.EqualError(t, err, "env: parse error on field \"Var\" of type \"env_test.foo\": something broke")
	})

	t.Run("slice", func(t *testing.T) {
		type config struct {
			Var []foo `env:"VAR2"`
		}
		os.Setenv("VAR2", "slice,slace")

		cfg := &config{}
		err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(foo{}): customParserFunc,
		})

		assert.Empty(t, cfg.Var)
		assert.EqualError(t, err, "env: parse error on field \"Var\" of type \"[]env_test.foo\": something broke")
	})
}

func TestCustomParserBasicType(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_VAL"`
	}

	exp := ConstT(123)
	os.Setenv("CONST_VAL", fmt.Sprintf("%d", exp))

	customParserFunc := func(v string) (interface{}, error) {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		r := ConstT(i)
		return r, nil
	}

	cfg := &config{}
	err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	assert.NoError(t, err)
	assert.Equal(t, exp, cfg.Const)
}

func TestCustomParserUint64Alias(t *testing.T) {
	type T uint64

	var one T = 1

	type config struct {
		Val T `env:"VAL" envDefault:"1x"`
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

	err := env.ParseWithFuncs(&cfg, env.CustomParsers{
		reflect.TypeOf(one): tParser,
	})

	assert.True(t, parserCalled, "tParser should have been called")
	assert.NoError(t, err)
	assert.Equal(t, T(1), cfg.Val)
}

func TestTypeCustomParserBasicInvalid(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_VAL"`
	}

	os.Setenv("CONST_VAL", "foobar")

	customParserFunc := func(_ string) (interface{}, error) {
		return nil, errors.New("random error")
	}

	cfg := &config{}
	err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	assert.Empty(t, cfg.Const)
	assert.EqualError(t, err, "env: parse error on field \"Const\" of type \"env_test.ConstT\": random error")
}

func TestCustomParserNotCalledForNonAlias(t *testing.T) {
	type T uint64
	type U uint64

	type config struct {
		Val   uint64 `env:"VAL" envDefault:"33"`
		Other U      `env:"OTHER" envDefault:"44"`
	}

	tParserCalled := false

	tParser := func(value string) (interface{}, error) {
		tParserCalled = true
		return T(99), nil
	}

	cfg := config{}

	err := env.ParseWithFuncs(&cfg, env.CustomParsers{
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
		Const ConstT `env:"CONST_VAL"`
	}

	os.Setenv("CONST_VAL", "42")

	cfg := &config{}
	err := env.Parse(cfg)

	assert.Zero(t, cfg.Const)
	assert.EqualError(t, err, "env: no parser found for field \"Const\" of type \"env_test.ConstT\"")
}

func TestUnsupportedStructType(t *testing.T) {
	type config struct {
		Foo http.Client `env:"FOO"`
	}

	os.Setenv("FOO", "foo")

	cfg := &config{}
	err := env.Parse(cfg)

	assert.EqualError(t, err, "env: no parser found for field \"Foo\" of type \"http.Client\"")
}

func TestEmptyOption(t *testing.T) {
	type config struct {
		Var string `env:"VAR,"`
	}

	cfg := &config{}

	os.Setenv("VAR", "val")
	defer os.Clearenv()
	assert.NoError(t, env.Parse(cfg))
	assert.Equal(t, "val", cfg.Var)
}

func TestErrorOptionNotRecognized(t *testing.T) {
	type config struct {
		Var string `env:"VAR,not_supported!"`
	}

	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: tag option \"not_supported!\" not supported")
}

func TestTextUnmarshalerError(t *testing.T) {
	type config struct {
		Unmarshaler unmarshaler `env:"UNMARSHALER"`
	}
	os.Setenv("UNMARSHALER", "invalid")
	cfg := &config{}
	assert.EqualError(t, env.Parse(cfg), "env: parse error on field \"Unmarshaler\" of type \"env_test.unmarshaler\": time: invalid duration invalid")
}

func TestParseURL(t *testing.T) {
	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL" envDefault:"https://google.com"`
	}
	var cfg config
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "https://google.com", cfg.ExampleURL.String())
}

func TestParseInvalidURL(t *testing.T) {
	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL_2"`
	}
	var cfg config
	os.Setenv("EXAMPLE_URL_2", "nope://s s/")
	assert.EqualError(t, env.Parse(&cfg), "env: parse error on field \"ExampleURL\" of type \"url.URL\": unable parse URL: parse nope://s s/: invalid character \" \" in host name")
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
	if err := env.Parse(&cfg); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Home:/tmp/fakehome Port:3000 IsProduction:false Inner:{Foo:foobar}}
}

func ExampleParseWithFuncs() {
	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL" envDefault:"https://google.com"`
	}
	var cfg config
	if err := env.ParseWithFuncs(&cfg, env.CustomParsers{
		parsers.URLType: parsers.URLFunc,
	}); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Println(cfg.ExampleURL.String())
	// Output: https://google.com
}
