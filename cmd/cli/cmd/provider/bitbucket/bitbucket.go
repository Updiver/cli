package bitbucket

import (
	"log"
	"os"
	"path"

	bitbucket "github.com/ktrysmt/go-bitbucket"
	"github.com/spf13/cobra"
	"github.com/updiver/dumper"
)

var (
	logger = log.New(os.Stdout, "bitbucket | ", log.Ldate|log.Ltime|log.Lmicroseconds)

	Username          string
	Token             string
	DestinationFolder string
	BitbucketCmd      = &cobra.Command{
		Use:   "bitbucket",
		Short: "bitbucket clones repositories by using user creds passed in",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Println("dumping bitbucket repositories")

			client := bitbucket.NewBasicAuth(Username, Token)
			client.Pagelen = 10
			client.DisableAutoPaging = false

			workspaces, err := GetWorkspaces(client)
			if err != nil {
				logger.Fatalf("get workspaces: %s\n", err)
			}

			workspaceSlugs := GetWorkspaceSlugs(workspaces)
			for _, workspaceSlug := range workspaceSlugs {
				logger.Printf("= workspace: %s\n", workspaceSlug)

				workspaceRepos, err := client.Repositories.ListForAccount(&bitbucket.RepositoriesOptions{
					Owner: workspaceSlug,
					Role:  "member",
				})
				if err != nil {
					log.Printf("get repositories: %s\n", err)
					continue
				}

				for _, repository := range workspaceRepos.Items {
					logger.Printf("== repository: %s\n", repository.Name)
					if cloneLinks, ok := repository.Links["clone"]; ok {
						for _, link := range cloneLinks.([]interface{}) {
							if link.(map[string]interface{})["name"] == "https" {
								logger.Printf("=== clone link: %s\n", link.(map[string]interface{})["href"])
								httpsCloneLink := link.(map[string]interface{})["href"].(string)

								fullDestFolder := path.Join(DestinationFolder, workspaceSlug, repository.Name)
								logger.Printf("=== clone repository to: %s\n", fullDestFolder)

								dpr := dumper.New()
								err = dpr.DumpRepository(httpsCloneLink, fullDestFolder, Username, Token)
								if err != nil {
									logger.Printf("dump repository: %s\n", err)
									continue
								}

								logger.Println("Repo cloned OK")
							}
						}
					} else {
						logger.Printf("=== repository %s has no clone link\n", repository.Name)
						continue
					}
				}
			}
		},
	}
)

// Workspaces

func GetWorkspaces(client *bitbucket.Client) (*bitbucket.WorkspaceList, error) {
	workspaces, err := client.Workspaces.List()
	return workspaces, err
}

func GetWorkspaceNames(workspaces *bitbucket.WorkspaceList) []string {
	wList := []string{}
	for _, workspace := range workspaces.Workspaces {
		wList = append(wList, workspace.Name)
	}
	return wList
}

func GetWorkspaceSlugs(workspaces *bitbucket.WorkspaceList) []string {
	wList := []string{}
	for _, workspace := range workspaces.Workspaces {
		wList = append(wList, workspace.Slug)
	}
	return wList
}

// Projects

func GetProjects(client *bitbucket.Client, workspaceSlug string) ([]bitbucket.Project, error) {
	projects, err := client.Workspaces.Projects(workspaceSlug)
	return projects.Items, err
}

func init() {
	BitbucketCmd.Flags().StringVarP(&Username, "username", "u", "", "username for git hosting account")
	BitbucketCmd.Flags().StringVarP(&Token, "token", "t", "", "token which is given by git provider")
	BitbucketCmd.Flags().StringVarP(&DestinationFolder, "destFolder", "d", "", "destination folder where repositories will be cloned to")

	BitbucketCmd.MarkFlagRequired("username")
	BitbucketCmd.MarkFlagRequired("token")
	BitbucketCmd.MarkFlagRequired("destFolder")
	BitbucketCmd.MarkFlagsRequiredTogether("username", "token", "destFolder")
}
