package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedEnv := make(Environment)
		expectedEnv["BAR"] = "bar"
		expectedEnv["FOO"] = "   foo\nwith new line"
		expectedEnv["HELLO"] = "\"hello\""
		expectedEnv["UNSET"] = ""

		env, err := ReadDir("testdata/env")

		require.NoError(t, err)
		require.Equal(t, expectedEnv, env)
	})

	t.Run("invalid file name", func(t *testing.T) {
		tmpDir, err := ioutil.TempDir("", "env")
		if err != nil {
			t.Fatal("can't create temp dir:", err)
		}
		defer os.RemoveAll(tmpDir)

		f, err := os.Create(path.Join(tmpDir, "LOG="))
		if err != nil {
			t.Fatal("can't create file", err)
		}
		f.Close()

		env, err := ReadDir(tmpDir)

		require.Nil(t, env)
		require.Error(t, err)
	})

	t.Run("directory instead file", func(t *testing.T) {
		tmpDir, err := ioutil.TempDir("", "env")
		if err != nil {
			t.Fatal("can't create temp dir:", err)
		}
		defer os.RemoveAll(tmpDir)

		err = os.Mkdir(path.Join(tmpDir, "dir"), os.ModePerm)
		if err != nil {
			t.Fatal("can't create file", err)
		}

		env, err := ReadDir(tmpDir)

		require.Nil(t, env)
		require.Error(t, err)
	})
}
