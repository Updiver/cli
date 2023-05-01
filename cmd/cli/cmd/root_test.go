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
	testCmd.Execute()

	expected := "dumper-cli dumps repositories under user account it is specified via command line arguments"
	if strings.Trim(output.String(), "\n") != expected {
		t.Errorf("expected %s, got %s", expected, output.String())
	}
}
