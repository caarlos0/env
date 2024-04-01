//go:build go1.18
// +build go1.18

package env

import "testing"

type Conf struct {
	Foo string `env:"FOO" envDefault:"bar"`
}

func TestParseAs(t *testing.T) {
	config, err := ParseAs[Conf]()
	isNoErr(t, err)
	isEqual(t, "bar", config.Foo)
}

func TestParseAsWithOptions(t *testing.T) {
	config, err := ParseAsWithOptions[Conf](Options{
		Environment: map[string]string{
			"FOO": "not bar",
		},
	})
	isNoErr(t, err)
	isEqual(t, "not bar", config.Foo)
}

type ConfRequired struct {
	Foo string `env:"FOO,required"`
}

func TestMust(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		defer func() {
			err := recover()
			isErrorWithMessage(t, err.(error), `env: required environment variable "FOO" is not set`)
		}()
		conf := Must(ParseAs[ConfRequired]())
		isEqual(t, "", conf.Foo)
	})
	t.Run("success", func(t *testing.T) {
		t.Setenv("FOO", "bar")
		conf := Must(ParseAs[ConfRequired]())
		isEqual(t, "bar", conf.Foo)
	})
}
