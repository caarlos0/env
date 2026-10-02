package env

import "testing"

func TestIgnoredNestedFields(t *testing.T) {
	type child struct {
		Value string `env:"VALUE"`
	}
	cfg := struct {
		Pointer *child `env:"-" envPrefix:"POINTER_"`
		Inline  struct {
			Value string `env:"VALUE"`
		} `env:"-" envPrefix:"INLINE_"`
		OptionPointer *child `env:"SKIP,-" envPrefix:"OPTION_POINTER_"`
		OptionInline  struct {
			Value string `env:"VALUE"`
		} `env:"SKIP,-" envPrefix:"OPTION_INLINE_"`
		NilPointer     *child `env:"-,init"`
		ControlPointer *child `envPrefix:"CONTROL_POINTER_"`
		ControlInline  struct {
			Value string `env:"VALUE"`
		} `envPrefix:"CONTROL_INLINE_"`
	}{Pointer: &child{Value: "keep"}, OptionPointer: &child{Value: "keep"}, ControlPointer: &child{}}
	cfg.Inline.Value = "keep"
	cfg.OptionInline.Value = "keep"
	var setKeys []string
	opts := Options{
		Environment: map[string]string{
			"POINTER_VALUE": "changed", "INLINE_VALUE": "changed",
			"OPTION_POINTER_VALUE": "changed", "OPTION_INLINE_VALUE": "changed",
			"CONTROL_POINTER_VALUE": "parsed", "CONTROL_INLINE_VALUE": "parsed",
		},
		OnSet: func(key string, _ interface{}, _ bool) { setKeys = append(setKeys, key) },
	}
	isNoErr(t, ParseWithOptions(&cfg, opts))
	for name, value := range map[string]string{
		"pointer": cfg.Pointer.Value, "inline": cfg.Inline.Value,
		"option pointer": cfg.OptionPointer.Value, "option inline": cfg.OptionInline.Value,
	} {
		t.Run(name, func(t *testing.T) { isEqual(t, "keep", value) })
	}
	isEqual(t, (*child)(nil), cfg.NilPointer)
	isEqual(t, "parsed", cfg.ControlPointer.Value)
	isEqual(t, "parsed", cfg.ControlInline.Value)
	t.Run("hooks", func(t *testing.T) {
		isEqual(t, []string{"CONTROL_POINTER_VALUE", "CONTROL_INLINE_VALUE"}, setKeys)
	})
	t.Run("field parameters", func(t *testing.T) {
		params, err := GetFieldParamsWithOptions(&cfg, opts)
		isNoErr(t, err)
		isEqual(t, []FieldParams{
			{OwnKey: "VALUE", Key: "CONTROL_POINTER_VALUE"},
			{OwnKey: "VALUE", Key: "CONTROL_INLINE_VALUE"},
		}, params)
	})
}

func TestIgnoredNestedFieldsCustomTag(t *testing.T) {
	type child struct {
		Value string `config:"VALUE"`
	}
	cfg := struct {
		Pointer *child `config:"-"`
		Inline  struct {
			Value string `config:"VALUE"`
		} `config:"SKIP,-"`
		Control string `config:"VALUE"`
	}{Pointer: &child{Value: "keep"}}
	cfg.Inline.Value = "keep"
	opts := Options{TagName: "config", Environment: map[string]string{"VALUE": "parsed"}}
	isNoErr(t, ParseWithOptions(&cfg, opts))
	isEqual(t, "keep", cfg.Pointer.Value)
	isEqual(t, "keep", cfg.Inline.Value)
	isEqual(t, "parsed", cfg.Control)
	params, err := GetFieldParamsWithOptions(&cfg, opts)
	isNoErr(t, err)
	isEqual(t, []FieldParams{{OwnKey: "VALUE", Key: "VALUE"}}, params)
}

func TestIgnoredNestedFieldsSkipRequiredValues(t *testing.T) {
	type child struct {
		Value string `env:"VALUE,required"`
	}
	cfg := struct {
		Pointer *child `env:"-"`
		Inline  struct {
			Value string `env:"VALUE,required"`
		} `env:"-"`
	}{Pointer: &child{}}
	isNoErr(t, ParseWithOptions(&cfg, Options{Environment: map[string]string{}}))
}

func TestNestedFieldTagValidationUnchanged(t *testing.T) {
	type child struct {
		Value string `env:"VALUE"`
	}
	cfg := struct {
		Pointer *child `env:",unsupported"`
		Inline  struct {
			Value string `env:"VALUE"`
		} `env:",unsupported"`
	}{Pointer: &child{}}
	isNoErr(t, ParseWithOptions(&cfg, Options{Environment: map[string]string{"VALUE": "parsed"}}))
	isEqual(t, "parsed", cfg.Pointer.Value)
	isEqual(t, "parsed", cfg.Inline.Value)
}
