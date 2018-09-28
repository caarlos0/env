package env_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/caarlos0/env"
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

type Config struct {
	Some            string `env:"somevar"`
	Other           bool   `env:"othervar"`
	Port            int    `env:"PORT"`
	Int64Val        int64  `env:"INT64VAL"`
	UintVal         uint   `env:"UINTVAL"`
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
	os.Setenv("UINT64VAL", "6464")
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
	assert.Equal(t, int64(-7575), cfg.Int64Val)
	assert.Equal(t, uint64(6464), cfg.Uint64Val)
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
	assert.Equal(t, []unmarshaler{unmarshaler{time.Minute * 2}, unmarshaler{time.Minute * 3}}, cfg.Unmarshalers)
	assert.Equal(t, []*unmarshaler{&unmarshaler{time.Minute * 2}, &unmarshaler{time.Minute * 3}}, cfg.UnmarshalerPtrs)
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
	assert.Error(t, env.Parse(&cfg))
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
	assert.Error(t, env.Parse(&thisShouldBreak))
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidBool(t *testing.T) {
	os.Setenv("othervar", "should-be-a-bool")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("PORT", "should-be-an-int")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidUint(t *testing.T) {
	os.Setenv("UINTVAL", "-44")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidFloat32(t *testing.T) {
	os.Setenv("FLOAT32", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidFloat64(t *testing.T) {
	os.Setenv("FLOAT64", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidUint64(t *testing.T) {
	os.Setenv("UINT64VAL", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidInt64(t *testing.T) {
	os.Setenv("INT64VAL", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidInt64Slice(t *testing.T) {
	type config struct {
		BadFloats []int64 `env:"BADINTS"`
	}

	os.Setenv("BADINTS", "A,2,3")
	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidUInt64Slice(t *testing.T) {
	type config struct {
		BadFloats []uint64 `env:"BADINTS"`
	}

	os.Setenv("BADFLOATS", "A,2,3")
	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidFloat32Slice(t *testing.T) {
	type config struct {
		BadFloats []float32 `env:"BADFLOATS"`
	}

	os.Setenv("BADFLOATS", "A,2.0,3.0")
	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidFloat64Slice(t *testing.T) {
	type config struct {
		BadFloats []float64 `env:"BADFLOATS"`
	}

	os.Setenv("BADFLOATS", "A,2.0,3.0")
	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidBoolsSlice(t *testing.T) {
	type config struct {
		BadBools []bool `env:"BADBOOLS"`
	}

	os.Setenv("BADBOOLS", "t,f,TRUE,faaaalse")
	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidDuration(t *testing.T) {
	os.Setenv("DURATION", "should-be-a-valid-duration")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidDurations(t *testing.T) {
	os.Setenv("DURATIONS", "1s,contains-an-invalid-duration,3s")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
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
	assert.Error(t, env.Parse(&cfg))
}

func TestUnsupportedSliceType(t *testing.T) {
	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}

	os.Setenv("WONTWORK", "1,2,3")
	defer os.Clearenv()

	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestBadSeparator(t *testing.T) {
	type config struct {
		WontWork []int `env:"WONTWORK" envSeparator:":"`
	}

	cfg := &config{}
	os.Setenv("WONTWORK", "1,2,3,4")
	defer os.Clearenv()

	assert.Error(t, env.Parse(cfg))
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
	assert.Error(t, env.Parse(cfg))
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
	assert.Error(t, err)
	assert.Equal(t, err, env.ErrNotAStructPtr)
}

func TestParseWithFuncsInvalidType(t *testing.T) {
	var c int
	err := env.ParseWithFuncs(&c, nil)
	assert.Error(t, err)
	assert.Equal(t, err, env.ErrNotAStructPtr)
}

func TestCustomParserError(t *testing.T) {
	type foo struct {
		name string
	}

	type config struct {
		Var foo `env:"VAR"`
	}

	os.Setenv("VAR", "test")

	customParserFunc := func(v string) (interface{}, error) {
		return nil, errors.New("something broke")
	}

	cfg := &config{}
	err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(foo{}): customParserFunc,
	})

	assert.Empty(t, cfg.Var.name, "Var.name should not be filled out when parse errors")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "VAR: custom parser: something broke")
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

func TypeCustomParserBasicInvalid(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_VAL"`
	}

	os.Setenv("CONST_VAL", "foobar")

	expErr := errors.New("Random error")
	customParserFunc := func(_ string) (interface{}, error) {
		return nil, expErr
	}

	cfg := &config{}
	err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	assert.Empty(t, cfg.Const)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func TestCustomParserBasicUnsupported(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_VAL"`
	}

	exp := ConstT(123)
	os.Setenv("CONST_VAL", fmt.Sprintf("%d", exp))

	cfg := &config{}
	err := env.Parse(cfg)

	assert.Zero(t, cfg.Const)
	assert.Error(t, err)
	assert.Equal(t, errors.New("CONST_VAL: " + env.ErrUnsupportedType.Error()), err)
}

func TestUnsupportedStructType(t *testing.T) {
	type config struct {
		Foo http.Client `env:"FOO"`
	}

	os.Setenv("FOO", "foo")

	cfg := &config{}
	err := env.Parse(cfg)

	assert.Error(t, err)
	assert.Equal(t, errors.New("FOO: " + env.ErrUnsupportedType.Error()), err)
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
	assert.Error(t, env.Parse(cfg))

}

func TestTextUnmarshalerError(t *testing.T) {
	type config struct {
		Unmarshaler unmarshaler `env:"UNMARSHALER"`
	}
	os.Setenv("UNMARSHALER", "invalid")
	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestMinMaxBorders(t *testing.T) {
	type config struct {
		СoordinateX         int     `env:"СOORDINATE_X" envMinValue:"-40" envMaxValue:"60"`
		СoordinateZ         uint    `env:"СOORDINATE_Z" envMinValue:"33" envMaxValue:"167"`
		Altitude            int64   `env:"ALTITUDE" envMinValue:"-10994" envMaxValue:"8848"`
		ServerPort          uint64  `env:"SERVER_PORT" envMaxValue:"1023"`
		MercuryTemperature  float32 `env:"MERCURY_TEMPERATURE" envMinValue:"-38.83" envMaxValue:"356.73"`
		NitrogenTemperature float64 `env:"NITROGEN_TEMPERATURE" envMinValue:"-209.86" envMaxValue:"-195.75"`
		Scores              []int64 `env:"SCORES" envMinValue:"-300" envMaxValue:"300"`
	}

	initVars := func() {
		os.Setenv("СOORDINATE_X", "-14")
		os.Setenv("СOORDINATE_Z", "122")
		os.Setenv("ALTITUDE", "300")
		os.Setenv("SERVER_PORT", "88")
		os.Setenv("MERCURY_TEMPERATURE", "20.2")
		os.Setenv("NITROGEN_TEMPERATURE", "-200.5")
		os.Setenv("SCORES", "-300,-200,100,250")
	}

	cfg := &config{}
	initVars()
	err := env.Parse(cfg)
	assert.NoError(t, err)

	initVars()
	os.Setenv("СOORDINATE_X", "-44")

	err = env.Parse(cfg)
	assert.Error(t, err)
	assert.Equal(t, errors.New("СOORDINATE_X: integer value -44 less than -40"), err)

	initVars()
	os.Setenv("СOORDINATE_Z", "667")

	err = env.Parse(cfg)
	assert.Error(t, err)
	assert.Equal(t, errors.New("СOORDINATE_Z: integer value 667 great than 167"), err)

	initVars()
	os.Setenv("ALTITUDE", "10000")

	err = env.Parse(cfg)
	assert.Error(t, err)
	assert.Equal(t, errors.New("ALTITUDE: integer value 10000 great than 8848"), err)

	initVars()
	os.Setenv("SERVER_PORT", "8080")

	err = env.Parse(cfg)
	assert.Error(t, err)
	assert.Equal(t, errors.New("SERVER_PORT: integer value 8080 great than 1023"), err)

	initVars()
	os.Setenv("MERCURY_TEMPERATURE", "-55.5")

	err = env.Parse(cfg)
	assert.Error(t, err)
	assert.Equal(t, errors.New("MERCURY_TEMPERATURE: float value -55.5 less than -38.83"), err)

	initVars()
	os.Setenv("NITROGEN_TEMPERATURE", "-0.5")

	err = env.Parse(cfg)
	assert.Error(t, err)
	assert.Equal(t, errors.New("NITROGEN_TEMPERATURE: float value -0.5 great than -195.75"), err)

	initVars()
	os.Setenv("SCORES", "-250,350,250,0")

	err = env.Parse(cfg)
	assert.Error(t, err)
	assert.Equal(t, errors.New("SCORES: integer value 350 great than 300"), err)
}

func ExampleParse() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	env.Parse(&cfg)
	fmt.Println(cfg)
	// Output: {/tmp/fakehome 3000 false}
}

func ExampleParseRequiredField() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
		SecretKey    string `env:"SECRET_KEY,required"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	err := env.Parse(&cfg)
	fmt.Println(err)
	// Output: required environment variable "SECRET_KEY" is not set
}

func ExampleParseMultipleOptions() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
		SecretKey    string `env:"SECRET_KEY,required,option1"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	err := env.Parse(&cfg)
	fmt.Println(err)
	// Output: env tag option "option1" not supported
}
