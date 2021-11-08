package env

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
)

type unmarshaler struct {
	time.Duration
}

// TextUnmarshaler implements encoding.TextUnmarshaler.
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
	String     string    `env:"STRING"`
	StringPtr  *string   `env:"STRING"`
	Strings    []string  `env:"STRINGS"`
	StringPtrs []*string `env:"STRINGS"`

	Bool     bool    `env:"BOOL"`
	BoolPtr  *bool   `env:"BOOL"`
	Bools    []bool  `env:"BOOLS"`
	BoolPtrs []*bool `env:"BOOLS"`

	Int     int    `env:"INT"`
	IntPtr  *int   `env:"INT"`
	Ints    []int  `env:"INTS"`
	IntPtrs []*int `env:"INTS"`

	Int8     int8    `env:"INT8"`
	Int8Ptr  *int8   `env:"INT8"`
	Int8s    []int8  `env:"INT8S"`
	Int8Ptrs []*int8 `env:"INT8S"`

	Int16     int16    `env:"INT16"`
	Int16Ptr  *int16   `env:"INT16"`
	Int16s    []int16  `env:"INT16S"`
	Int16Ptrs []*int16 `env:"INT16S"`

	Int32     int32    `env:"INT32"`
	Int32Ptr  *int32   `env:"INT32"`
	Int32s    []int32  `env:"INT32S"`
	Int32Ptrs []*int32 `env:"INT32S"`

	Int64     int64    `env:"INT64"`
	Int64Ptr  *int64   `env:"INT64"`
	Int64s    []int64  `env:"INT64S"`
	Int64Ptrs []*int64 `env:"INT64S"`

	Uint     uint    `env:"UINT"`
	UintPtr  *uint   `env:"UINT"`
	Uints    []uint  `env:"UINTS"`
	UintPtrs []*uint `env:"UINTS"`

	Uint8     uint8    `env:"UINT8"`
	Uint8Ptr  *uint8   `env:"UINT8"`
	Uint8s    []uint8  `env:"UINT8S"`
	Uint8Ptrs []*uint8 `env:"UINT8S"`

	Uint16     uint16    `env:"UINT16"`
	Uint16Ptr  *uint16   `env:"UINT16"`
	Uint16s    []uint16  `env:"UINT16S"`
	Uint16Ptrs []*uint16 `env:"UINT16S"`

	Uint32     uint32    `env:"UINT32"`
	Uint32Ptr  *uint32   `env:"UINT32"`
	Uint32s    []uint32  `env:"UINT32S"`
	Uint32Ptrs []*uint32 `env:"UINT32S"`

	Uint64     uint64    `env:"UINT64"`
	Uint64Ptr  *uint64   `env:"UINT64"`
	Uint64s    []uint64  `env:"UINT64S"`
	Uint64Ptrs []*uint64 `env:"UINT64S"`

	Float32     float32    `env:"FLOAT32"`
	Float32Ptr  *float32   `env:"FLOAT32"`
	Float32s    []float32  `env:"FLOAT32S"`
	Float32Ptrs []*float32 `env:"FLOAT32S"`

	Float64     float64    `env:"FLOAT64"`
	Float64Ptr  *float64   `env:"FLOAT64"`
	Float64s    []float64  `env:"FLOAT64S"`
	Float64Ptrs []*float64 `env:"FLOAT64S"`

	Duration     time.Duration    `env:"DURATION"`
	Durations    []time.Duration  `env:"DURATIONS"`
	DurationPtr  *time.Duration   `env:"DURATION"`
	DurationPtrs []*time.Duration `env:"DURATIONS"`

	Unmarshaler     unmarshaler    `env:"UNMARSHALER"`
	UnmarshalerPtr  *unmarshaler   `env:"UNMARSHALER"`
	Unmarshalers    []unmarshaler  `env:"UNMARSHALERS"`
	UnmarshalerPtrs []*unmarshaler `env:"UNMARSHALERS"`

	URL     url.URL    `env:"URL"`
	URLPtr  *url.URL   `env:"URL"`
	URLs    []url.URL  `env:"URLS"`
	URLPtrs []*url.URL `env:"URLS"`

	StringWithdefault string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/db"`

	CustomSeparator []string `env:"SEPSTRINGS" envSeparator:":"`

	NonDefined struct {
		String string `env:"NONDEFINED_STR"`
	}

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
	is := is.New(t)

	defer os.Clearenv()

	tos := func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}

	toss := func(v ...interface{}) string {
		ss := []string{}
		for _, s := range v {
			ss = append(ss, tos(s))
		}
		return strings.Join(ss, ",")
	}

	str1 := "str1"
	str2 := "str2"
	os.Setenv("STRING", str1)
	os.Setenv("STRINGS", toss(str1, str2))

	bool1 := true
	bool2 := false
	os.Setenv("BOOL", tos(bool1))
	os.Setenv("BOOLS", toss(bool1, bool2))

	int1 := -1
	int2 := 2
	os.Setenv("INT", tos(int1))
	os.Setenv("INTS", toss(int1, int2))

	var int81 int8 = -2
	var int82 int8 = 5
	os.Setenv("INT8", tos(int81))
	os.Setenv("INT8S", toss(int81, int82))

	var int161 int16 = -24
	var int162 int16 = 15
	os.Setenv("INT16", tos(int161))
	os.Setenv("INT16S", toss(int161, int162))

	var int321 int32 = -14
	var int322 int32 = 154
	os.Setenv("INT32", tos(int321))
	os.Setenv("INT32S", toss(int321, int322))

	var int641 int64 = -12
	var int642 int64 = 150
	os.Setenv("INT64", tos(int641))
	os.Setenv("INT64S", toss(int641, int642))

	var uint1 uint = 1
	var uint2 uint = 2
	os.Setenv("UINT", tos(uint1))
	os.Setenv("UINTS", toss(uint1, uint2))

	var uint81 uint8 = 15
	var uint82 uint8 = 51
	os.Setenv("UINT8", tos(uint81))
	os.Setenv("UINT8S", toss(uint81, uint82))

	var uint161 uint16 = 532
	var uint162 uint16 = 123
	os.Setenv("UINT16", tos(uint161))
	os.Setenv("UINT16S", toss(uint161, uint162))

	var uint321 uint32 = 93
	var uint322 uint32 = 14
	os.Setenv("UINT32", tos(uint321))
	os.Setenv("UINT32S", toss(uint321, uint322))

	var uint641 uint64 = 5
	var uint642 uint64 = 43
	os.Setenv("UINT64", tos(uint641))
	os.Setenv("UINT64S", toss(uint641, uint642))

	var float321 float32 = 9.3
	var float322 float32 = 1.1
	os.Setenv("FLOAT32", tos(float321))
	os.Setenv("FLOAT32S", toss(float321, float322))

	float641 := 1.53
	float642 := 0.5
	os.Setenv("FLOAT64", tos(float641))
	os.Setenv("FLOAT64S", toss(float641, float642))

	duration1 := time.Second
	duration2 := time.Second * 4
	os.Setenv("DURATION", tos(duration1))
	os.Setenv("DURATIONS", toss(duration1, duration2))

	unmarshaler1 := unmarshaler{time.Minute}
	unmarshaler2 := unmarshaler{time.Millisecond * 1232}
	os.Setenv("UNMARSHALER", tos(unmarshaler1.Duration))
	os.Setenv("UNMARSHALERS", toss(unmarshaler1.Duration, unmarshaler2.Duration))

	url1 := "https://goreleaser.com"
	url2 := "https://caarlos0.dev"
	os.Setenv("URL", tos(url1))
	os.Setenv("URLS", toss(url1, url2))

	os.Setenv("SEPSTRINGS", strings.Join([]string{str1, str2}, ":"))

	nonDefinedStr := "nonDefinedStr"
	os.Setenv("NONDEFINED_STR", nonDefinedStr)

	cfg := Config{}
	is.NoErr(Parse(&cfg))

	is.Equal(str1, cfg.String)
	is.Equal(&str1, cfg.StringPtr)
	is.Equal(str1, cfg.Strings[0])
	is.Equal(str2, cfg.Strings[1])
	is.Equal(&str1, cfg.StringPtrs[0])
	is.Equal(&str2, cfg.StringPtrs[1])

	is.Equal(bool1, cfg.Bool)
	is.Equal(&bool1, cfg.BoolPtr)
	is.Equal(bool1, cfg.Bools[0])
	is.Equal(bool2, cfg.Bools[1])
	is.Equal(&bool1, cfg.BoolPtrs[0])
	is.Equal(&bool2, cfg.BoolPtrs[1])

	is.Equal(int1, cfg.Int)
	is.Equal(&int1, cfg.IntPtr)
	is.Equal(int1, cfg.Ints[0])
	is.Equal(int2, cfg.Ints[1])
	is.Equal(&int1, cfg.IntPtrs[0])
	is.Equal(&int2, cfg.IntPtrs[1])

	is.Equal(int81, cfg.Int8)
	is.Equal(&int81, cfg.Int8Ptr)
	is.Equal(int81, cfg.Int8s[0])
	is.Equal(int82, cfg.Int8s[1])
	is.Equal(&int81, cfg.Int8Ptrs[0])
	is.Equal(&int82, cfg.Int8Ptrs[1])

	is.Equal(int161, cfg.Int16)
	is.Equal(&int161, cfg.Int16Ptr)
	is.Equal(int161, cfg.Int16s[0])
	is.Equal(int162, cfg.Int16s[1])
	is.Equal(&int161, cfg.Int16Ptrs[0])
	is.Equal(&int162, cfg.Int16Ptrs[1])

	is.Equal(int321, cfg.Int32)
	is.Equal(&int321, cfg.Int32Ptr)
	is.Equal(int321, cfg.Int32s[0])
	is.Equal(int322, cfg.Int32s[1])
	is.Equal(&int321, cfg.Int32Ptrs[0])
	is.Equal(&int322, cfg.Int32Ptrs[1])

	is.Equal(int641, cfg.Int64)
	is.Equal(&int641, cfg.Int64Ptr)
	is.Equal(int641, cfg.Int64s[0])
	is.Equal(int642, cfg.Int64s[1])
	is.Equal(&int641, cfg.Int64Ptrs[0])
	is.Equal(&int642, cfg.Int64Ptrs[1])

	is.Equal(uint1, cfg.Uint)
	is.Equal(&uint1, cfg.UintPtr)
	is.Equal(uint1, cfg.Uints[0])
	is.Equal(uint2, cfg.Uints[1])
	is.Equal(&uint1, cfg.UintPtrs[0])
	is.Equal(&uint2, cfg.UintPtrs[1])

	is.Equal(uint81, cfg.Uint8)
	is.Equal(&uint81, cfg.Uint8Ptr)
	is.Equal(uint81, cfg.Uint8s[0])
	is.Equal(uint82, cfg.Uint8s[1])
	is.Equal(&uint81, cfg.Uint8Ptrs[0])
	is.Equal(&uint82, cfg.Uint8Ptrs[1])

	is.Equal(uint161, cfg.Uint16)
	is.Equal(&uint161, cfg.Uint16Ptr)
	is.Equal(uint161, cfg.Uint16s[0])
	is.Equal(uint162, cfg.Uint16s[1])
	is.Equal(&uint161, cfg.Uint16Ptrs[0])
	is.Equal(&uint162, cfg.Uint16Ptrs[1])

	is.Equal(uint321, cfg.Uint32)
	is.Equal(&uint321, cfg.Uint32Ptr)
	is.Equal(uint321, cfg.Uint32s[0])
	is.Equal(uint322, cfg.Uint32s[1])
	is.Equal(&uint321, cfg.Uint32Ptrs[0])
	is.Equal(&uint322, cfg.Uint32Ptrs[1])

	is.Equal(uint641, cfg.Uint64)
	is.Equal(&uint641, cfg.Uint64Ptr)
	is.Equal(uint641, cfg.Uint64s[0])
	is.Equal(uint642, cfg.Uint64s[1])
	is.Equal(&uint641, cfg.Uint64Ptrs[0])
	is.Equal(&uint642, cfg.Uint64Ptrs[1])

	is.Equal(float321, cfg.Float32)
	is.Equal(&float321, cfg.Float32Ptr)
	is.Equal(float321, cfg.Float32s[0])
	is.Equal(float322, cfg.Float32s[1])
	is.Equal(&float321, cfg.Float32Ptrs[0])
	is.Equal(&float322, cfg.Float32Ptrs[1])

	is.Equal(float641, cfg.Float64)
	is.Equal(&float641, cfg.Float64Ptr)
	is.Equal(float641, cfg.Float64s[0])
	is.Equal(float642, cfg.Float64s[1])
	is.Equal(&float641, cfg.Float64Ptrs[0])
	is.Equal(&float642, cfg.Float64Ptrs[1])

	is.Equal(duration1, cfg.Duration)
	is.Equal(&duration1, cfg.DurationPtr)
	is.Equal(duration1, cfg.Durations[0])
	is.Equal(duration2, cfg.Durations[1])
	is.Equal(&duration1, cfg.DurationPtrs[0])
	is.Equal(&duration2, cfg.DurationPtrs[1])

	is.Equal(unmarshaler1, cfg.Unmarshaler)
	is.Equal(&unmarshaler1, cfg.UnmarshalerPtr)
	is.Equal(unmarshaler1, cfg.Unmarshalers[0])
	is.Equal(unmarshaler2, cfg.Unmarshalers[1])
	is.Equal(&unmarshaler1, cfg.UnmarshalerPtrs[0])
	is.Equal(&unmarshaler2, cfg.UnmarshalerPtrs[1])

	is.Equal(url1, cfg.URL.String())
	is.Equal(url1, cfg.URLPtr.String())
	is.Equal(url1, cfg.URLs[0].String())
	is.Equal(url2, cfg.URLs[1].String())
	is.Equal(url1, cfg.URLPtrs[0].String())
	is.Equal(url2, cfg.URLPtrs[1].String())

	is.Equal("postgres://localhost:5432/db", cfg.StringWithdefault)
	is.Equal(nonDefinedStr, cfg.NonDefined.String)

	is.Equal(str1, cfg.CustomSeparator[0])
	is.Equal(str2, cfg.CustomSeparator[1])

	is.Equal(cfg.NotAnEnv, "")

	is.Equal(cfg.unexported, "")
}

func TestSetEnvAndTagOptsChain(t *testing.T) {
	is := is.New(t)

	defer os.Clearenv()
	type config struct {
		Key1 string `mytag:"KEY1,required"`
		Key2 int    `mytag:"KEY2,required"`
	}
	envs := map[string]string{
		"KEY1": "VALUE1",
		"KEY2": "3",
	}

	cfg := config{}
	is.NoErr(Parse(&cfg, Options{TagName: "mytag"}, Options{Environment: envs}))
	is.Equal("VALUE1", cfg.Key1)
	is.Equal(3, cfg.Key2)
}

func TestJSONTag(t *testing.T) {
	is := is.New(t)

	defer os.Clearenv()
	type config struct {
		Key1 string `json:"KEY1"`
		Key2 int    `json:"KEY2"`
	}

	os.Setenv("KEY1", "VALUE7")
	os.Setenv("KEY2", "5")

	cfg := config{}
	is.NoErr(Parse(&cfg, Options{TagName: "json"}))
	is.Equal("VALUE7", cfg.Key1)
	is.Equal(5, cfg.Key2)
}

func TestParsesEnvInner(t *testing.T) {
	is := is.New(t)

	os.Setenv("innervar", "someinnervalue")
	os.Setenv("innernum", "8")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
		unexported:  &InnerStruct{},
	}
	is.NoErr(Parse(&cfg))
	is.Equal("someinnervalue", cfg.InnerStruct.Inner)
	is.Equal(uint(8), cfg.InnerStruct.Number)
}

func TestParsesEnvInnerFails(t *testing.T) {
	defer os.Clearenv()
	type config struct {
		Foo struct {
			Number int `env:"NUMBER"`
		}
	}
	os.Setenv("NUMBER", "not-a-number")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "Number" of type "int": strconv.ParseInt: parsing "not-a-number": invalid syntax`)
}

func TestParsesEnvInnerNil(t *testing.T) {
	is := is.New(t)

	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{}
	is.NoErr(Parse(&cfg))
}

func TestParsesEnvInnerInvalid(t *testing.T) {
	os.Setenv("innernum", "-547")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
	}
	isErrorWithMessage(t, Parse(&cfg), `env: parse error on field "Number" of type "uint": strconv.ParseUint: parsing "-547": invalid syntax`)
}

func TestParsesEnvNested(t *testing.T) {
	is := is.New(t)

	os.Setenv("nestedvar", "somenestedvalue")
	defer os.Clearenv()
	var cfg ForNestedStruct
	is.NoErr(Parse(&cfg))
	is.Equal("somenestedvalue", cfg.NestedVar)
}

func TestEmptyVars(t *testing.T) {
	is := is.New(t)

	os.Clearenv()
	cfg := Config{}
	is.NoErr(Parse(&cfg))
	is.Equal("", cfg.String)
	is.Equal(false, cfg.Bool)
	is.Equal(0, cfg.Int)
	is.Equal(uint(0), cfg.Uint)
	is.Equal(uint64(0), cfg.Uint64)
	is.Equal(int64(0), cfg.Int64)
	is.Equal(0, len(cfg.Strings))
	is.Equal(0, len(cfg.CustomSeparator))
	is.Equal(0, len(cfg.Ints))
	is.Equal(0, len(cfg.Bools))
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	isErrorWithMessage(t, Parse(&thisShouldBreak), "env: expected a pointer to a Struct")
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	isErrorWithMessage(t, Parse(cfg), "env: expected a pointer to a Struct")
}

func TestInvalidBool(t *testing.T) {
	os.Setenv("BOOL", "should-be-a-bool")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Bool" of type "bool": strconv.ParseBool: parsing "should-be-a-bool": invalid syntax`)
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("INT", "should-be-an-int")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Int" of type "int": strconv.ParseInt: parsing "should-be-an-int": invalid syntax`)
}

func TestInvalidUint(t *testing.T) {
	os.Setenv("UINT", "-44")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Uint" of type "uint": strconv.ParseUint: parsing "-44": invalid syntax`)
}

func TestInvalidFloat32(t *testing.T) {
	os.Setenv("FLOAT32", "AAA")
	defer os.Clearenv()

	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Float32" of type "float32": strconv.ParseFloat: parsing "AAA": invalid syntax`)
}

func TestInvalidFloat64(t *testing.T) {
	os.Setenv("FLOAT64", "AAA")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Float64" of type "float64": strconv.ParseFloat: parsing "AAA": invalid syntax`)
}

func TestInvalidUint64(t *testing.T) {
	os.Setenv("UINT64", "AAA")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Uint64" of type "uint64": strconv.ParseUint: parsing "AAA": invalid syntax`)
}

func TestInvalidInt64(t *testing.T) {
	os.Setenv("INT64", "AAA")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Int64" of type "int64": strconv.ParseInt: parsing "AAA": invalid syntax`)
}

func TestInvalidInt64Slice(t *testing.T) {
	os.Setenv("BADINTS", "A,2,3")
	defer os.Clearenv()
	type config struct {
		BadFloats []int64 `env:"BADINTS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadFloats" of type "[]int64": strconv.ParseInt: parsing "A": invalid syntax`)
}

func TestInvalidUInt64Slice(t *testing.T) {
	os.Setenv("BADINTS", "A,2,3")
	defer os.Clearenv()
	type config struct {
		BadFloats []uint64 `env:"BADINTS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadFloats" of type "[]uint64": strconv.ParseUint: parsing "A": invalid syntax`)
}

func TestInvalidFloat32Slice(t *testing.T) {
	os.Setenv("BADFLOATS", "A,2.0,3.0")
	defer os.Clearenv()
	type config struct {
		BadFloats []float32 `env:"BADFLOATS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadFloats" of type "[]float32": strconv.ParseFloat: parsing "A": invalid syntax`)
}

func TestInvalidFloat64Slice(t *testing.T) {
	os.Setenv("BADFLOATS", "A,2.0,3.0")
	defer os.Clearenv()
	type config struct {
		BadFloats []float64 `env:"BADFLOATS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadFloats" of type "[]float64": strconv.ParseFloat: parsing "A": invalid syntax`)
}

func TestInvalidBoolsSlice(t *testing.T) {
	os.Setenv("BADBOOLS", "t,f,TRUE,faaaalse")
	defer os.Clearenv()
	type config struct {
		BadBools []bool `env:"BADBOOLS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadBools" of type "[]bool": strconv.ParseBool: parsing "faaaalse": invalid syntax`)
}

func TestInvalidDuration(t *testing.T) {
	os.Setenv("DURATION", "should-be-a-valid-duration")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Duration" of type "time.Duration": unable to parse duration: time: invalid duration "should-be-a-valid-duration"`)
}

func TestInvalidDurations(t *testing.T) {
	os.Setenv("DURATIONS", "1s,contains-an-invalid-duration,3s")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Durations" of type "[]time.Duration": unable to parse duration: time: invalid duration "contains-an-invalid-duration"`)
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	is := is.New(t)

	cfg := Config{}
	is.NoErr(Parse(&cfg))
	is.Equal(cfg.NotAnEnv, "")
}

func TestParseStructWithInvalidFieldKind(t *testing.T) {
	type config struct {
		WontWorkByte byte `env:"BLAH"`
	}
	os.Setenv("BLAH", "a")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "WontWorkByte" of type "uint8": strconv.ParseUint: parsing "a": invalid syntax`)
}

func TestUnsupportedSliceType(t *testing.T) {
	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}

	os.Setenv("WONTWORK", "1,2,3")
	defer os.Clearenv()

	isErrorWithMessage(t, Parse(&config{}), `env: no parser found for field "WontWork" of type "[]map[int]int"`)
}

func TestBadSeparator(t *testing.T) {
	type config struct {
		WontWork []int `env:"WONTWORK" envSeparator:":"`
	}

	os.Setenv("WONTWORK", "1,2,3,4")
	defer os.Clearenv()

	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "WontWork" of type "[]int": strconv.ParseInt: parsing "1,2,3,4": invalid syntax`)
}

func TestNoErrorRequiredSet(t *testing.T) {
	is := is.New(t)

	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}

	cfg := &config{}

	os.Setenv("IS_REQUIRED", "")
	defer os.Clearenv()
	is.NoErr(Parse(cfg))
	is.Equal("", cfg.IsRequired)
}

func TestHook(t *testing.T) {
	is := is.New(t)

	type config struct {
		Something string `env:"SOMETHING" envDefault:"important"`
		Another   string `env:"ANOTHER"`
	}

	cfg := &config{}

	os.Setenv("ANOTHER", "1")
	defer os.Clearenv()

	type onSetArgs struct {
		tag       string
		key       interface{}
		isDefault bool
	}

	var onSetCalled []onSetArgs

	is.NoErr(Parse(cfg, Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			onSetCalled = append(onSetCalled, onSetArgs{tag, value, isDefault})
		},
	}))
	is.Equal("important", cfg.Something)
	is.Equal("1", cfg.Another)
	is.Equal(2, len(onSetCalled))
	is.Equal(onSetArgs{"SOMETHING", "important", true}, onSetCalled[0])
	is.Equal(onSetArgs{"ANOTHER", "1", false}, onSetCalled[1])
}

func TestErrorRequiredWithDefault(t *testing.T) {
	is := is.New(t)

	type config struct {
		IsRequired string `env:"IS_REQUIRED,required" envDefault:"important"`
	}

	cfg := &config{}

	os.Setenv("IS_REQUIRED", "")
	defer os.Clearenv()
	is.NoErr(Parse(cfg))
	is.Equal("", cfg.IsRequired)
}

func TestErrorRequiredNotSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: required environment variable "IS_REQUIRED" is not set`)
}

func TestNoErrorNotEmptySet(t *testing.T) {
	is := is.New(t)
	os.Setenv("IS_REQUIRED", "1")
	defer os.Clearenv()
	type config struct {
		IsRequired string `env:"IS_REQUIRED,notEmpty"`
	}
	is.NoErr(Parse(&config{}))
}

func TestNoErrorRequiredAndNotEmptySet(t *testing.T) {
	is := is.New(t)
	os.Setenv("IS_REQUIRED", "1")
	defer os.Clearenv()
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required,notEmpty"`
	}
	is.NoErr(Parse(&config{}))
}

func TestErrorNotEmptySet(t *testing.T) {
	os.Setenv("IS_REQUIRED", "")
	defer os.Clearenv()
	type config struct {
		IsRequired string `env:"IS_REQUIRED,notEmpty"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: environment variable "IS_REQUIRED" should not be empty`)
}

func TestErrorRequiredAndNotEmptySet(t *testing.T) {
	os.Setenv("IS_REQUIRED", "")
	defer os.Clearenv()
	type config struct {
		IsRequired string `env:"IS_REQUIRED,notEmpty,required"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: environment variable "IS_REQUIRED" should not be empty`)
}

func TestErrorRequiredNotSetWithDefault(t *testing.T) {
	is := is.New(t)

	type config struct {
		IsRequired string `env:"IS_REQUIRED,required" envDefault:"important"`
	}

	cfg := &config{}

	is.NoErr(Parse(cfg))
	is.Equal("important", cfg.IsRequired)
}

func TestParseExpandOption(t *testing.T) {
	is := is.New(t)

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

	is.NoErr(err)
	is.Equal("localhost", cfg.Host)
	is.Equal(3000, cfg.Port)
	is.Equal("qwerty12345", cfg.SecretKey)
	is.Equal("qwerty12345", cfg.ExpandKey)
	is.Equal("localhost:3000", cfg.CompoundKey)
	is.Equal("def1", cfg.Default)
}

func TestParseUnsetRequireOptions(t *testing.T) {
	is := is.New(t)

	type config struct {
		Password string `env:"PASSWORD,unset,required"`
	}
	defer os.Clearenv()
	cfg := config{}

	isErrorWithMessage(t, Parse(&cfg), `env: required environment variable "PASSWORD" is not set`)
	os.Setenv("PASSWORD", "superSecret")
	is.NoErr(Parse(&cfg))

	is.Equal("superSecret", cfg.Password)
	unset, exists := os.LookupEnv("PASSWORD")
	is.Equal("", unset)
	is.Equal(false, exists)
}

func TestCustomParser(t *testing.T) {
	is := is.New(t)

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

	is.NoErr(err)
	is.Equal(cfg.Var.name, "test")
	is.Equal(cfg.Foo.name, "test3")
	is.Equal(cfg.Other.Name, "test2")
	is.Equal(cfg.Other.Foo.name, "test3")
}

func TestParseWithFuncsNoPtr(t *testing.T) {
	type foo struct{}
	isErrorWithMessage(t, ParseWithFuncs(foo{}, nil), "env: expected a pointer to a Struct")
}

func TestParseWithFuncsInvalidType(t *testing.T) {
	var c int
	isErrorWithMessage(t, ParseWithFuncs(&c, nil), "env: expected a pointer to a Struct")
}

func TestCustomParserError(t *testing.T) {
	type foo struct {
		name string
	}

	customParserFunc := func(v string) (interface{}, error) {
		return nil, errors.New("something broke")
	}

	t.Run("single", func(t *testing.T) {
		is := is.New(t)

		type config struct {
			Var foo `env:"VAR"`
		}

		os.Setenv("VAR", "single")
		cfg := &config{}
		err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
			reflect.TypeOf(foo{}): customParserFunc,
		})

		is.Equal(cfg.Var.name, "")
		isErrorWithMessage(t, err, `env: parse error on field "Var" of type "env.foo": something broke`)
	})

	t.Run("slice", func(t *testing.T) {
		is := is.New(t)

		type config struct {
			Var []foo `env:"VAR2"`
		}
		os.Setenv("VAR2", "slice,slace")

		cfg := &config{}
		err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
			reflect.TypeOf(foo{}): customParserFunc,
		})

		is.Equal(cfg.Var, nil)
		isErrorWithMessage(t, err, `env: parse error on field "Var" of type "[]env.foo": something broke`)
	})
}

func TestCustomParserBasicType(t *testing.T) {
	is := is.New(t)

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

	is.NoErr(err)
	is.Equal(exp, cfg.Const)
}

func TestCustomParserUint64Alias(t *testing.T) {
	is := is.New(t)

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

	is.True(parserCalled) // tParser should have been called
	is.NoErr(err)
	is.Equal(T(1), cfg.Val)
}

func TestTypeCustomParserBasicInvalid(t *testing.T) {
	is := is.New(t)

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

	is.Equal(cfg.Const, ConstT(0))
	isErrorWithMessage(t, err, `env: parse error on field "Const" of type "env.ConstT": random error`)
}

func TestCustomParserNotCalledForNonAlias(t *testing.T) {
	is := is.New(t)

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

	is.True(!tParserCalled) // tParser should not have been called
	is.NoErr(err)
	is.Equal(uint64(33), cfg.Val)
	is.Equal(U(44), cfg.Other)
}

func TestCustomParserBasicUnsupported(t *testing.T) {
	is := is.New(t)

	type ConstT struct {
		A int
	}

	type config struct {
		Const ConstT `env:"CONST_"`
	}

	os.Setenv("CONST_", "42")

	cfg := &config{}
	err := Parse(cfg)

	is.Equal(cfg.Const, ConstT{0})
	isErrorWithMessage(t, err, `env: no parser found for field "Const" of type "env.ConstT"`)
}

func TestUnsupportedStructType(t *testing.T) {
	type config struct {
		Foo http.Client `env:"FOO"`
	}
	os.Setenv("FOO", "foo")
	defer os.Clearenv()
	isErrorWithMessage(t, Parse(&config{}), `env: no parser found for field "Foo" of type "http.Client"`)
}

func TestEmptyOption(t *testing.T) {
	is := is.New(t)

	type config struct {
		Var string `env:"VAR,"`
	}

	cfg := &config{}

	os.Setenv("VAR", "")
	defer os.Clearenv()
	is.NoErr(Parse(cfg))
	is.Equal("", cfg.Var)
}

func TestErrorOptionNotRecognized(t *testing.T) {
	type config struct {
		Var string `env:"VAR,not_supported!"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: tag option "not_supported!" not supported`)
}

func TestTextUnmarshalerError(t *testing.T) {
	type config struct {
		Unmarshaler unmarshaler `env:"UNMARSHALER"`
	}
	os.Setenv("UNMARSHALER", "invalid")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "Unmarshaler" of type "env.unmarshaler": time: invalid duration "invalid"`)
}

func TestTextUnmarshalersError(t *testing.T) {
	type config struct {
		Unmarshalers []unmarshaler `env:"UNMARSHALERS"`
	}
	os.Setenv("UNMARSHALERS", "1s,invalid")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "Unmarshalers" of type "[]env.unmarshaler": time: invalid duration "invalid"`)
}

func TestParseURL(t *testing.T) {
	is := is.New(t)

	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL" envDefault:"https://google.com"`
	}
	var cfg config
	is.NoErr(Parse(&cfg))
	is.Equal("https://google.com", cfg.ExampleURL.String())
}

func TestParseInvalidURL(t *testing.T) {
	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL_2"`
	}
	os.Setenv("EXAMPLE_URL_2", "nope://s s/")

	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "ExampleURL" of type "url.URL": unable to parse URL: parse "nope://s s/": invalid character " " in host name`)
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

func ExampleParse_onSet() {
	type config struct {
		Home         string `env:"HOME,required"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	var cfg config
	if err := Parse(&cfg, Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			fmt.Printf("Set %s to %v (default? %v)\n", tag, value, isDefault)
		},
	}); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: Set HOME to /tmp/fakehome (default? false)
	// Set PORT to 3000 (default? true)
	// Set PRODUCTION to  (default? false)
	// {Home:/tmp/fakehome Port:3000 IsProduction:false}
}

func ExampleParse_defaults() {
	type config struct {
		A string `env:"FOO" envDefault:"foo"`
		B string `env:"FOO"`
	}

	// env FOO is not set

	cfg := config{
		A: "A",
		B: "B",
	}
	if err := Parse(&cfg); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {A:foo B:B}
}

func TestIgnoresUnexported(t *testing.T) {
	is := is.New(t)

	type unexportedConfig struct {
		home  string `env:"HOME"`
		Home2 string `env:"HOME"`
	}
	cfg := unexportedConfig{}

	os.Setenv("HOME", "/tmp/fakehome")
	is.NoErr(Parse(&cfg))
	is.Equal(cfg.home, "")
	is.Equal("/tmp/fakehome", cfg.Home2)
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
	is := is.New(t)

	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_LEVELS", "debug,info")
	defer os.Unsetenv("LOG_LEVEL")
	defer os.Unsetenv("LOG_LEVELS")

	type config struct {
		LogLevel  LogLevel   `env:"LOG_LEVEL"`
		LogLevels []LogLevel `env:"LOG_LEVELS"`
	}
	var cfg config

	is.NoErr(Parse(&cfg))
	is.Equal(DebugLevel, cfg.LogLevel)
	is.Equal([]LogLevel{DebugLevel, InfoLevel}, cfg.LogLevels)
}

func ExampleParseWithFuncs() {
	type thing struct {
		desc string
	}

	type conf struct {
		Thing thing `env:"THING"`
	}

	os.Setenv("THING", "my thing")

	c := conf{}

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

func TestFile(t *testing.T) {
	is := is.New(t)

	type config struct {
		SecretKey string `env:"SECRET_KEY,file"`
	}

	dir := t.TempDir()
	file := filepath.Join(dir, "sec_key")
	is.NoErr(os.WriteFile(file, []byte("secret"), 0o660))

	defer os.Clearenv()
	os.Setenv("SECRET_KEY", file)

	cfg := config{}
	is.NoErr(Parse(&cfg))
	is.Equal("secret", cfg.SecretKey)
}

func TestFileNoParam(t *testing.T) {
	is := is.New(t)

	type config struct {
		SecretKey string `env:"SECRET_KEY,file"`
	}
	defer os.Clearenv()

	cfg := config{}
	is.NoErr(Parse(&cfg))
}

func TestFileNoParamRequired(t *testing.T) {
	type config struct {
		SecretKey string `env:"SECRET_KEY,file,required"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: required environment variable "SECRET_KEY" is not set`)
}

func TestFileBadFile(t *testing.T) {
	type config struct {
		SecretKey string `env:"SECRET_KEY,file"`
	}

	filename := "not-a-real-file"
	defer os.Clearenv()
	os.Setenv("SECRET_KEY", filename)

	oserr := "no such file or directory"
	if runtime.GOOS == "windows" {
		oserr = "The system cannot find the file specified."
	}
	isErrorWithMessage(t, Parse(&config{}), fmt.Sprintf(`env: could not load content of file "%s" from variable SECRET_KEY: open %s: %s`, filename, filename, oserr))
}

func TestFileWithDefault(t *testing.T) {
	is := is.New(t)

	type config struct {
		SecretKey string `env:"SECRET_KEY,file" envDefault:"${FILE}" envExpand:"true"`
	}
	defer os.Clearenv()

	dir := t.TempDir()
	file := filepath.Join(dir, "sec_key")
	is.NoErr(os.WriteFile(file, []byte("secret"), 0o660))

	defer os.Clearenv()
	os.Setenv("FILE", file)

	cfg := config{}
	is.NoErr(Parse(&cfg))
	is.Equal("secret", cfg.SecretKey)
}

func TestCustomSliceType(t *testing.T) {
	is := is.New(t)

	type customslice []byte

	type config struct {
		SecretKey customslice `env:"SECRET_KEY"`
	}

	parsecustomsclice := func(value string) (interface{}, error) {
		return customslice(value), nil
	}

	defer os.Clearenv()
	os.Setenv("SECRET_KEY", "somesecretkey")

	var cfg config
	is.NoErr(ParseWithFuncs(&cfg, map[reflect.Type]ParserFunc{reflect.TypeOf(customslice{}): parsecustomsclice}))
}

func TestBlankKey(t *testing.T) {
	is := is.New(t)

	type testStruct struct {
		Blank        string
		BlankWithTag string `env:""`
	}

	val := testStruct{}

	defer os.Clearenv()
	os.Setenv("", "You should not see this")

	is.NoErr(Parse(&val))
	is.Equal("", val.Blank)
	is.Equal("", val.BlankWithTag)
}

type MyTime time.Time

func (t *MyTime) UnmarshalText(text []byte) error {
	tt, err := time.Parse("2006-01-02", string(text))
	*t = MyTime(tt)
	return err
}

func TestCustomTimeParser(t *testing.T) {
	is := is.New(t)

	type config struct {
		SomeTime MyTime `env:"SOME_TIME"`
	}

	os.Setenv("SOME_TIME", "2021-05-06")
	defer os.Unsetenv("SOME_TIME")

	var cfg config
	is.NoErr(Parse(&cfg))
	is.Equal(2021, time.Time(cfg.SomeTime).Year())
	is.Equal(time.Month(5), time.Time(cfg.SomeTime).Month())
	is.Equal(6, time.Time(cfg.SomeTime).Day())
}

func TestRequiredIfNoDefOption(t *testing.T) {
	type Tree struct {
		Fruit string `env:"FRUIT"`
	}
	type config struct {
		Name  string `env:"NAME"`
		Genre string `env:"GENRE" envDefault:"Unknown"`
		Tree
	}
	var cfg config

	t.Run("missing", func(t *testing.T) {
		isErrorWithMessage(t, Parse(&cfg, Options{RequiredIfNoDef: true}), `env: required environment variable "NAME" is not set`)
		os.Setenv("NAME", "John")
		t.Cleanup(os.Clearenv)
		isErrorWithMessage(t, Parse(&cfg, Options{RequiredIfNoDef: true}), `env: required environment variable "FRUIT" is not set`)
	})

	t.Run("all set", func(t *testing.T) {
		os.Setenv("NAME", "John")
		os.Setenv("FRUIT", "Apple")
		t.Cleanup(os.Clearenv)

		// should not trigger an error for the missing 'GENRE' env because it has a default value.
		is.New(t).NoErr(Parse(&cfg, Options{RequiredIfNoDef: true}))
	})
}

func TestPrefix(t *testing.T) {
	is := is.New(t)
	type Config struct {
		Home string `env:"HOME"`
	}
	type ComplexConfig struct {
		Foo   Config `envPrefix:"FOO_"`
		Bar   Config `envPrefix:"BAR_"`
		Clean Config
	}
	cfg := ComplexConfig{}
	err := Parse(&cfg, Options{Environment: map[string]string{"FOO_HOME": "/foo", "BAR_HOME": "/bar", "HOME": "/clean"}})
	is.NoErr(err)
	is.Equal("/foo", cfg.Foo.Home)
	is.Equal("/bar", cfg.Bar.Home)
	is.Equal("/clean", cfg.Clean.Home)
}

func TestComplePrefix(t *testing.T) {
	is := is.New(t)
	type Config struct {
		Home string `env:"HOME"`
	}
	type ComplexConfig struct {
		Foo   Config `envPrefix:"FOO_"`
		Clean Config
		Bar   Config `envPrefix:"BAR_"`
		Blah  string `env:"BLAH"`
	}
	cfg := ComplexConfig{}
	err := Parse(&cfg, Options{
		Prefix: "T_",
		Environment: map[string]string{
			"T_FOO_HOME": "/foo",
			"T_BAR_HOME": "/bar",
			"T_BLAH":     "blahhh",
			"T_HOME":     "/clean",
		},
	})
	is.NoErr(err)
	is.Equal("/foo", cfg.Foo.Home)
	is.Equal("/bar", cfg.Bar.Home)
	is.Equal("/clean", cfg.Clean.Home)
	is.Equal("blahhh", cfg.Blah)
}

func isErrorWithMessage(tb testing.TB, err error, msg string) {
	tb.Helper()

	is := is.New(tb)
	is.True(err != nil)        // should have failed
	is.Equal(err.Error(), msg) // should have the expected message
}
