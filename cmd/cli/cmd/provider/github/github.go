package provider

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	ghDump "github.com/updiver/cli/dump/github"
	"github.com/updiver/cli/flags"
)

var (
	logger = log.New(os.Stdout, "github | ", log.Ldate|log.Ltime|log.Lmicroseconds)

	Username          string
	Token             string
	DestinationFolder string
	CloneMode         flags.CloneMode

	GithubCmd = &cobra.Command{
		Use:   "github",
		Short: "github clones repositories by using user creds passed in",
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
			logger.Println("dumping github repositories")

			cloneModeFlag := cmd.Flag("clone-mode")
			CloneMode = flags.CloneMode(cloneModeFlag.Value.String())

			err := ghDump.Dump(&ghDump.DumpOptions{
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
	GithubCmd.Flags().StringVarP(&Username, "username", "u", "", "username for git hosting account")
	GithubCmd.Flags().StringVarP(&Token, "token", "t", "", "token which is given by git provider")
	GithubCmd.Flags().StringVarP(&DestinationFolder, "destFolder", "d", "", "destination folder where repositories will be cloned to")

	GithubCmd.MarkFlagRequired("username")
	GithubCmd.MarkFlagRequired("token")
	GithubCmd.MarkFlagRequired("destFolder")
	GithubCmd.MarkFlagsRequiredTogether("username", "token", "destFolder")
}
