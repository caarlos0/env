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
