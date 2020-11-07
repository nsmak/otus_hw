package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, info := range infos {
		if info.IsDir() {
			return nil, fmt.Errorf("it's not a file: %v", info.Name())
		}
		if strings.Contains(info.Name(), "=") {
			return nil, fmt.Errorf("invalid variable name: %v", info.Name())
		}

		name := info.Name()

		file, err := os.Open(path.Join(dir, name))
		if err != nil {
			return nil, err
		}

		r := bufio.NewReader(file)
		l, _, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				env[name] = ""
				continue
			}
			return nil, err
		}

		l = bytes.ReplaceAll(l, []byte("\x00"), []byte("\n"))
		l = bytes.TrimRight(l, "\t \n")

		var builder strings.Builder
		builder.Write(l)

		env[name] = builder.String()
	}

	return env, nil
}
