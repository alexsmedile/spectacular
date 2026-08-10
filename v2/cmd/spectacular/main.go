package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alexsmedile/spectacular/v2/internal/buildinfo"
	"github.com/alexsmedile/spectacular/v2/internal/command"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("spectacular %s\n", buildinfo.Version)
		return
	}
	if len(os.Args) == 3 && os.Args[1] == "--version" && os.Args[2] == "--json" {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]string{
			"schema_version": command.VersionInspection.Schema,
			"name":           "spectacular",
			"version":        buildinfo.Version,
			"commit":         buildinfo.Commit,
		})
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(3)
	}
	os.Exit(command.Runner{Cwd: cwd, Stdout: os.Stdout, Stderr: os.Stderr}.Run(os.Args[1:]))
}
