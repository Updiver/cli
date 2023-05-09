package cmd

import (
	"github.com/spf13/cobra"
	"github.com/updiver/cli/flags"
)

var (
	cloneMode string
	rootCmd   = &cobra.Command{
		Use:   "dumper-cli",
		Short: "dumper-cli is a tool for dumping your repositories",
		Long:  `dumper-cli dumps repositories under user account it is specified via command line arguments`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&cloneMode,
		"clone-mode",
		"m",
		flags.CloneModeDefaultBranch,
		"clone mode (default-branch, all-branches)",
	)
}

func Execute() error {
	return rootCmd.Execute()
}
