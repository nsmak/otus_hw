package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cmd := []string{"testdata/echo.sh", "arg1", "arg2"}
		env := make(Environment)
		env["FOO"] = "123"
		env["BAR"] = "value"

		code := RunCmd(cmd, env)

		require.Equal(t, 0, code)
	})

	t.Run("command is nil", func(t *testing.T) {
		code := RunCmd(nil, nil)

		require.Equal(t, 1, code)
	})
}
