package env

import (
	"net/url"
	"reflect"
	"testing"
)

type scalarSliceRecord struct {
	Value string `env:"VALUE"`
}

type scalarSliceText struct {
	Value string `env:"VALUE"`
	Calls int    `env:"-"`
}

func (s *scalarSliceText) UnmarshalText(text []byte) error {
	s.Value = string(text)
	s.Calls++
	return nil
}

func TestScalarBoundStructSlices(t *testing.T) {
	type explicitURLs struct {
		URLs []url.URL `env:"URLS" envPrefix:"ITEM_"`
	}
	type inferredURLs struct {
		URLs []url.URL
	}
	type defaultURLs struct {
		URLs []url.URL `env:"URLS" envPrefix:"ITEM_" envDefault:"https://one.example/a,https://two.example/b"`
	}
	type records struct {
		Values []scalarSliceRecord `env:"VALUES" envPrefix:"ITEM_"`
	}
	type textRecords struct {
		Values []scalarSliceText `env:"VALUES" envPrefix:"ITEM_"`
	}
	type defaultRecords struct {
		Values []scalarSliceRecord `env:"VALUES" envPrefix:"ITEM_" envDefault:"default"`
	}
	wantURLs := []url.URL{
		{Scheme: "https", Host: "one.example", Path: "/a"},
		{Scheme: "https", Host: "two.example", Path: "/b"},
	}
	tests := []struct {
		name          string
		config        any
		want          any
		inference     bool
		required      bool
		zeroOnly      bool
		customCalls   int
		withoutScalar bool
	}{
		{"explicit", &explicitURLs{}, &explicitURLs{URLs: wantURLs}, false, false, false, 0, false},
		{"explicit-required", &explicitURLs{}, &explicitURLs{URLs: wantURLs}, false, true, false, 0, false},
		{"explicit-inferred", &explicitURLs{}, &explicitURLs{URLs: wantURLs}, true, false, false, 0, false},
		{"explicit-inferred-required", &explicitURLs{}, &explicitURLs{URLs: wantURLs}, true, true, false, 0, false},
		{"inferred", &inferredURLs{}, &inferredURLs{URLs: wantURLs}, true, false, false, 0, false},
		{"inferred-required", &inferredURLs{}, &inferredURLs{URLs: wantURLs}, true, true, false, 0, false},
		{"default", &defaultURLs{}, &defaultURLs{URLs: wantURLs}, true, true, false, 0, true},
		{"custom-parser", &records{}, &records{Values: []scalarSliceRecord{{"one"}, {"two"}}}, true, true, false, 2, false},
		{"text-unmarshaler", &textRecords{}, &textRecords{Values: []scalarSliceText{{"one", 1}, {"two", 1}}}, true, true, false, 0, false},
		{"retained-nonzero-default", &defaultRecords{Values: []scalarSliceRecord{{"initial"}}}, &defaultRecords{Values: []scalarSliceRecord{{"initial"}}}, true, true, true, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			environment := map[string]string{
				"UR_LS":       "https://one.example/a,https://two.example/b",
				"URLS":        "https://one.example/a,https://two.example/b",
				"VALUES":      "one,two",
				"ITEM_0_HOST": "wrong.example", "ITEM_0_PATH": "/wrong",
				"ITEM_0_VALUE": "wrong", "ITEM_1_VALUE": "also-wrong", "ITEM_1_HOST": "also-wrong.example",
				"0_HOST": "wrong.example", "0_PATH": "/wrong",
			}
			if tt.withoutScalar {
				delete(environment, "URLS")
				delete(environment, "UR_LS")
				delete(environment, "VALUES")
			}
			err := ParseWithOptions(tt.config, Options{
				Environment: environment, UseFieldNameByDefault: tt.inference,
				RequiredIfNoDef: tt.required, SetDefaultsForZeroValuesOnly: tt.zeroOnly,
				FuncMap: map[reflect.Type]ParserFunc{
					reflect.TypeOf(scalarSliceRecord{}): func(value string) (interface{}, error) {
						calls++
						return scalarSliceRecord{Value: value}, nil
					},
				},
			})
			if err != nil {
				t.Fatalf("parse scalar-bound slice: %v", err)
			}
			if !reflect.DeepEqual(tt.config, tt.want) {
				t.Errorf("got %#v, want %#v", tt.config, tt.want)
			}
			if calls != tt.customCalls {
				t.Errorf("custom parser called %d times, want %d", calls, tt.customCalls)
			}
		})
	}
}

func TestScalarSliceIndexedContainers(t *testing.T) {
	type prefixRecords struct {
		Values []scalarSliceRecord `envPrefix:"ITEM_"`
	}
	type prefixText struct {
		Values []scalarSliceText `envPrefix:"ITEM_"`
	}
	type inferredRecords struct {
		Values []scalarSliceRecord
	}
	type pointerRecords struct {
		Values *[]scalarSliceRecord `envPrefix:"ITEM_"`
	}
	type pointerElements struct {
		Values []*scalarSliceRecord `envPrefix:"ITEM_"`
	}
	wantPointer := []scalarSliceRecord{{"indexed"}, {"tail"}}
	initialPointer := []scalarSliceRecord{{"initial"}, {"tail"}}
	tests := []struct {
		name        string
		config      any
		want        any
		inference   bool
		environment map[string]string
		customCalls int
	}{
		{"prefix-custom", &prefixRecords{}, &prefixRecords{Values: []scalarSliceRecord{{"indexed"}}}, false, map[string]string{"ITEM_0_VALUE": "indexed"}, 0},
		{"prefix-custom-inferred", &prefixRecords{}, &prefixRecords{Values: []scalarSliceRecord{{"indexed"}}}, true, map[string]string{"ITEM_0_VALUE": "indexed"}, 0},
		{"prefix-inferred-value-still-container", &prefixRecords{}, &prefixRecords{Values: []scalarSliceRecord{{"indexed"}}}, true, map[string]string{"VALUES": "scalar", "ITEM_0_VALUE": "indexed"}, 1},
		{"prefix-text", &prefixText{}, &prefixText{Values: []scalarSliceText{{"indexed", 0}}}, true, map[string]string{"ITEM_0_VALUE": "indexed"}, 0},
		{"missing-inferred", &inferredRecords{}, &inferredRecords{Values: []scalarSliceRecord{{"indexed"}}}, true, map[string]string{"0_VALUE": "indexed"}, 0},
		{"empty-inferred", &inferredRecords{}, &inferredRecords{Values: []scalarSliceRecord{{"indexed"}}}, true, map[string]string{"VALUES": "", "0_VALUE": "indexed"}, 0},
		{"pointer-container", &pointerRecords{Values: &initialPointer}, &pointerRecords{Values: &wantPointer}, false, map[string]string{"ITEM_0_VALUE": "indexed"}, 0},
		{"unsupported-pointer-elements-unchanged", &pointerElements{}, &pointerElements{}, false, map[string]string{"ITEM_0_VALUE": "indexed"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			err := ParseWithOptions(tt.config, Options{
				Environment: tt.environment, UseFieldNameByDefault: tt.inference,
				FuncMap: map[reflect.Type]ParserFunc{
					reflect.TypeOf(scalarSliceRecord{}): func(value string) (interface{}, error) {
						calls++
						return scalarSliceRecord{Value: value}, nil
					},
				},
			})
			if err != nil {
				t.Fatalf("parse indexed container: %v", err)
			}
			if !reflect.DeepEqual(tt.config, tt.want) {
				t.Errorf("got %#v, want %#v", tt.config, tt.want)
			}
			if calls != tt.customCalls {
				t.Errorf("custom parser called %d times, want %d", calls, tt.customCalls)
			}
		})
	}
}
