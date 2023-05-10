package bitbucket

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	bbDump "github.com/updiver/cli/dump/bitbucket"
	"github.com/updiver/cli/flags"
)

var (
	logger = log.New(os.Stdout, "bitbucket | ", log.Ldate|log.Ltime|log.Lmicroseconds)

	Username          string
	Token             string
	DestinationFolder string
	CloneMode         flags.CloneMode

	BitbucketCmd = &cobra.Command{
		Use:   "bitbucket",
		Short: "bitbucket clones repositories by using user creds passed in",
		Args: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if f.Name == "clone-mode" {
					if err := flags.CloneMode(f.Value.String()).Valid(); err != nil {
						logger.Fatalf("clone mode: %s\n", err)
					}
				}
			})

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Println("dumping bitbucket repositories")

			cloneModeFlag := cmd.Flag("clone-mode")
			CloneMode = flags.CloneMode(cloneModeFlag.Value.String())

			err := bbDump.Dump(&bbDump.DumpOptions{
				Username:    Username,
				Token:       Token,
				Destination: DestinationFolder,
				CloneMode:   CloneMode.String(),
			})
			if err != nil {
				fmt.Printf("dumping repositories: %s\n", err)
			}
		},
	}
)

func init() {
	BitbucketCmd.Flags().StringVarP(&Username, "username", "u", "", "username for git hosting account")
	BitbucketCmd.Flags().StringVarP(&Token, "token", "t", "", "token which is given by git provider")
	BitbucketCmd.Flags().StringVarP(&DestinationFolder, "destFolder", "d", "", "destination folder where repositories will be cloned to")

	BitbucketCmd.MarkFlagRequired("username")
	BitbucketCmd.MarkFlagRequired("token")
	BitbucketCmd.MarkFlagRequired("destFolder")
	BitbucketCmd.MarkFlagsRequiredTogether("username", "token", "destFolder")
}
