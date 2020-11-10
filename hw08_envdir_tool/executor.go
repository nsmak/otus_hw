package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Println("need input command")
		return 1
	}

	var args []string
	if len(cmd) > 1 {
		args = append(args, cmd[1:]...)
	}

	err := prepareOSEnv(env)
	if err != nil {
		fmt.Printf("can't prepare os environment: %v", err)
		return 1
	}

	command := exec.Command(cmd[0], args...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = os.Environ()

	err = command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		fmt.Printf("unknown command error: %v\n", err)
		return 1
	}
	return 0
}

func prepareOSEnv(env Environment) error {
	for key, value := range env {
		if len(value) == 0 {
			err := os.Unsetenv(key)
			if err != nil {
				return fmt.Errorf("cat't unset env for key %v: %w", key, err)
			}
			continue
		}
		err := os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("cat't set env %v for key %v: %w", value, key, err)
		}
	}
	return nil
}
