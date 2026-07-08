package main

import (
	"os"

	"github.com/martyn/apidiff/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
