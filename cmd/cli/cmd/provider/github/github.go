package provider

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/google/go-github/v52/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/updiver/cli/flags"
	"github.com/updiver/dumper"
	"golang.org/x/oauth2"
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

			ghClient := GetAuthenticatedClient(Token)
			allRepos, err := GetRepositories(ghClient)
			if err != nil {
				logger.Printf("get repositories: %s\n", err)
				return
			}

			for _, repo := range allRepos {
				logger.Printf("org [%s] | repo [%s]\n", *repo.Owner.Login, *repo.Name)
				if repo.CloneURL == nil {
					logger.Printf("skipping repo [%s] as it has no clone url\n", *repo.Name)
					continue
				}

				fullDestFolder := path.Join(DestinationFolder, *repo.Owner.Login, *repo.Name)
				logger.Printf("=== clone repository to: %s\n", fullDestFolder)
				dpr := dumper.New()
				opts := &dumper.DumpRepositoryOptions{
					RepositoryURL: *repo.CloneURL,
					Destination:   fullDestFolder,
					Creds:         dumper.Creds{Username: Username, Password: Token},
				}
				flags.ApplyCloneMode(opts, CloneMode)
				_, err = dpr.DumpRepository(opts)
				if err != nil {
					logger.Printf("dump repository: %s\n", err)
					continue
				}

				logger.Println("Repo cloned OK")
			}
		},
	}
)

func GetAuthenticatedClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetRepositories(ghClient *github.Client) ([]*github.Repository, error) {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := ghClient.Repositories.List(context.Background(), "", opts)
		if err != nil {
			logger.Printf("getting repositories list: %s\n", err)
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func init() {
	GithubCmd.Flags().StringVarP(&Username, "username", "u", "", "username for git hosting account")
	GithubCmd.Flags().StringVarP(&Token, "token", "t", "", "token which is given by git provider")
	GithubCmd.Flags().StringVarP(&DestinationFolder, "destFolder", "d", "", "destination folder where repositories will be cloned to")

	GithubCmd.MarkFlagRequired("username")
	GithubCmd.MarkFlagRequired("token")
	GithubCmd.MarkFlagRequired("destFolder")
	GithubCmd.MarkFlagsRequiredTogether("username", "token", "destFolder")
}
