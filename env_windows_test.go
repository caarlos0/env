package env

import (
	"testing"

	"github.com/matryer/is"
)

// On Windows, environment variables can start with '='. This test verifies this behavior without relying on a Windows environment.
// See env_windows.go in the Go source: https://github.com/golang/go/blob/master/src/syscall/env_windows.go#L58
func TestToMapWindows(t *testing.T) {
	is := is.New(t)
	envVars := []string{"=::=::\\", "=C:=C:\\test", "VAR=REGULARVAR"}
	result := toMap(envVars)
	is.Equal(map[string]string{
		"=::": "::\\",
		"=C:": "C:\\test",
		"VAR": "REGULARVAR",
	}, result)
}
