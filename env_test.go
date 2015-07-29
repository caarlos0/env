package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const key = "WHATEVER"

func TestGetenv(t *testing.T) {
	value := "bla"
	Set(key, value)
	defer Unset(key)
	assert.Equal(t, value, GetOr(key, "default"))
}

func TestGetenvUnseted(t *testing.T) {
	value := "default"
	assert.Equal(t, value, GetOr(key, value))
}

type Env struct {
	Some  string `env:"somevar"`
	Other bool   `env:"othervar"`
	Port  int    `env:"PORT"`
}

func TestParsesEnv(t *testing.T) {
	Set("somevar", "somevalue")
	Set("othervar", "true")
	Set("PORT", "8080")
	defer Unset("somevar")
	defer Unset("othervar")
	defer Unset("PORT")

	env := Env{}
	err := Parse(&env)

	assert.Nil(t, err)
	assert.Equal(t, "somevalue", env.Some)
	assert.Equal(t, true, env.Other)
	assert.Equal(t, 8080, env.Port)
}

func TestEmptyVars(t *testing.T) {
	env := Env{}
	err := Parse(&env)

	assert.Nil(t, err)
	assert.Equal(t, "", env.Some)
	assert.Equal(t, false, env.Other)
	assert.Equal(t, 0, env.Port)
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	assert.Error(t, Parse(&thisShouldBreak))
}

func TestPassReference(t *testing.T) {
	env := Env{}
	assert.Error(t, Parse(env))
}

func TestInvalidBool(t *testing.T) {
	Set("othervar", "should-be-a-bool")
	defer Unset("othervar")

	env := Env{}
	assert.Error(t, Parse(&env))
}

func TestInvalidInt(t *testing.T) {
	Set("PORT", "should-be-an-int")
	defer Unset("PORT")

	env := Env{}
	assert.Error(t, Parse(&env))
}
