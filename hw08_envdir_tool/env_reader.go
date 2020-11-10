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
		return nil, fmt.Errorf("can't read directory: %w", err)
	}

	var builder strings.Builder
	env := make(Environment)

	for _, info := range infos {
		if info.IsDir() {
			fmt.Printf("[WARR]: it's not a file: %v\n", info.Name())
			continue
		}
		if strings.Contains(info.Name(), "=") {
			fmt.Printf("[WARR]: invalid variable name: %v", info.Name())
			continue
		}

		name := info.Name()

		file, err := os.Open(path.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("can't open file: %w", err)
		}

		r := bufio.NewReader(file)
		l, _, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				env[name] = ""
				continue
			}
			file.Close()
			return nil, fmt.Errorf("can't read line from file %v: %w", name, err)
		}

		l = bytes.ReplaceAll(l, []byte("\x00"), []byte("\n"))
		l = bytes.TrimRight(l, "\t \n")

		builder.Write(l)
		env[name] = builder.String()
		builder.Reset()
		file.Close()
	}

	return env, nil
}
