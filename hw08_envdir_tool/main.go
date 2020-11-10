package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("directory is not exist")
		os.Exit(1)
	}

	d := os.Args[1]
	env, err := ReadDir(d)
	if err != nil {
		fmt.Printf("can't get environment variables from dir %v: %v\n", d, err)
		os.Exit(1)
	}

	exitCode := RunCmd(os.Args[2:], env)
	os.Exit(exitCode)
}
