//go:build windows

package env

import "testing"

// On Windows, environment variables can start with '='.
// This test verifies this behavior without relying on a Windows environment.
// See env_windows.go in the Go source: https://github.com/golang/go/blob/master/src/syscall/env_windows.go#L58
func TestToMapWindows(t *testing.T) {
	envVars := []string{"=::=::\\", "=C:=C:\\test", "VAR=REGULARVAR", "FOO=", "BAR"}
	result := ToMap(envVars)
	isEqual(t, map[string]string{
		"=::": "::\\",
		"=C:": "C:\\test",
		"VAR": "REGULARVAR",
		"FOO": "",
	}, result)
}
