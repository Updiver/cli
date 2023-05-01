package main

import (
	"fmt"
	"os"

	cmd "github.com/updiver/cli/cmd/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "dumper-cli error: %v\n", err)
		os.Exit(1)
	}
}
