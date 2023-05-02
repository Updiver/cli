package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "dumper-cli",
		Short: "dumper-cli is a tool for dumping your repositories",
		Long:  `dumper-cli dumps repositories under user account it is specified via command line arguments`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
