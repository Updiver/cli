package cmd

import (
	"github.com/spf13/cobra"
	bbProvider "github.com/updiver/cli/cmd/cli/cmd/provider/bitbucket"
	ghProvider "github.com/updiver/cli/cmd/cli/cmd/provider/github"
)

var (
	Username          string
	Token             string
	DestinationFolder string

	dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "dump clones repositories by using user creds passed in",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func init() {
	dumpCmd.AddCommand(bbProvider.BitbucketCmd)
	dumpCmd.AddCommand(ghProvider.GithubCmd)

	rootCmd.AddCommand(dumpCmd)
}
