package main

import (
	"os"

	"github.com/webknife/webknife/internal/adapters/cli"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
