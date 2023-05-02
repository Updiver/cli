package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version    string
	GitCommit  string
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "dumper-cli utility version",
		Long:  "shows dumper-cli utility version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dumper-cli version: %s\n", Version)
			fmt.Printf("build commit: %s\n", GitCommit)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
