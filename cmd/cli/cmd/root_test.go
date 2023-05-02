package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecute(t *testing.T) {
	output := bytes.NewBufferString("")
	rootCmd.SetOutput(output)

	testCmd := cobra.Command(*rootCmd)
	err := testCmd.Execute()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	expected := `
		dumper-cli dumps repositories under user account it is specified via command line arguments

		Usage:
			dumper-cli [command]

		Available Commands:
			completion  Generate the autocompletion script for the specified shell
			dump        dump clones repositories by using user creds passed in
			help        Help about any command
			version     dumper-cli utility version

		Flags:
			-h, --help   help for dumper-cli

		Use "dumper-cli [command] --help" for more information about a command.
	`
	if strings.EqualFold(output.String(), expected) {
		t.Errorf("expected %s, got %s", expected, output.String())
	}
}
