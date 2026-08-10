package main

import (
	"os"

	"github.com/alexsmedile/spectacular/v2/internal/command"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(3)
	}
	os.Exit(command.Runner{Cwd: cwd, Stdout: os.Stdout, Stderr: os.Stderr}.Run(os.Args[1:]))
}
