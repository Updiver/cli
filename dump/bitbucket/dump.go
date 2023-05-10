package bitbucket

import (
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/updiver/cli/flags"
	"github.com/updiver/dumper"
)

var ErrNoHttpsCloneLink = errors.New("no https clone link")

type DumpOptions struct {
	Username    string
	Token       string
	Destination string
	CloneMode   string
}

type WorkspaceRepositories bitbucket.RepositoriesRes

func (wr *WorkspaceRepositories) Iter(f func(repository *bitbucket.Repository) error) error {
	for _, repository := range wr.Items {
		if err := f(&repository); err != nil {
			return err
		}
	}

	return nil
}

type WorkspaceSlug string

func (wslug WorkspaceSlug) WorkspaceRepositories(client *bitbucket.Client) (*WorkspaceRepositories, error) {
	workspaceRepos, err := client.Repositories.ListForAccount(&bitbucket.RepositoriesOptions{
		Owner: string(wslug),
		Role:  "member",
	})

	w := WorkspaceRepositories(*workspaceRepos)
	return &w, err
}

type WorkspaceSlugs []WorkspaceSlug

func (wslugs WorkspaceSlugs) Iter(f func(workspaceSlug WorkspaceSlug) error) error {
	for _, workspaceSlug := range wslugs {
		if err := f(workspaceSlug); err != nil {
			return err
		}
	}

	return nil
}

func Dump(opts *DumpOptions) error {
	client := bitbucket.NewBasicAuth(opts.Username, opts.Token)
	client.Pagelen = 10
	client.DisableAutoPaging = false

	workspaces, err := getWorkspaces(client)
	if err != nil {
		return fmt.Errorf("get workspaces: %w", err)
	}

	workspaceSlugs := getWorkspaceSlugs(workspaces)
	err = workspaceSlugs.Iter(func(workspaceSlug WorkspaceSlug) error {
		repositories, err := workspaceSlug.WorkspaceRepositories(client)
		if err != nil {
			return err
		}

		err = repositories.Iter(wrapper(opts))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("traversing workspace repositories: %w", err)
	}

	return err
}

func wrapper(opts *DumpOptions) func(repository *bitbucket.Repository) error {
	return func(repository *bitbucket.Repository) error {
		return dumpRepo(repository, opts)
	}
}

func dumpRepo(repository *bitbucket.Repository, opts *DumpOptions) error {
	logger.Printf("workspace [%s] | repo [%s]\n", repository.Slug, repository.Name)
	if _, ok := repository.Links["clone"]; !ok {
		logger.Printf("repository %s has no clone link\n", repository.Name)
		return nil
	}

	httpsCloneLink, err := httpsCloneLink(repository.Links["clone"])
	if err != nil {
		return err
	}

	fullDestFolder := path.Join(opts.Destination, repository.Slug, repository.Name)
	dpr := dumper.New()
	dumpRepoOpts := &dumper.DumpRepositoryOptions{
		RepositoryURL: httpsCloneLink,
		Destination:   fullDestFolder,
		Creds: dumper.Creds{
			Username: opts.Username,
			Password: opts.Token,
		},
		Output: &dumper.Output{
			GitOutput: io.Discard,
		},
	}
	flags.ApplyCloneMode(dumpRepoOpts, flags.CloneMode(opts.CloneMode))
	_, err = dpr.DumpRepository(dumpRepoOpts)
	if err != nil {
		return fmt.Errorf("dump repository: %w", err)
	}

	// all-branches is a mirror clone, it's required to convert bare into non-bare
	// so user can properly use that repo
	if flags.CloneMode(opts.CloneMode) == flags.CloneModeAllBranches {
		err = dumper.Convert(fullDestFolder, dumper.RepositoryTypeNonBare)
		if err != nil {
			return fmt.Errorf("convert repository: %w", err)
		}
	}

	logger.Printf("repo [%s] dumped", repository.Name)
	return nil
}

func getWorkspaces(client *bitbucket.Client) (*bitbucket.WorkspaceList, error) {
	workspaces, err := client.Workspaces.List()
	return workspaces, err
}

func getWorkspaceSlugs(workspaces *bitbucket.WorkspaceList) WorkspaceSlugs {
	wList := make(WorkspaceSlugs, len(workspaces.Workspaces))
	for _, workspace := range workspaces.Workspaces {
		wList = append(wList, WorkspaceSlug(workspace.Slug))
	}
	return wList
}

func httpsCloneLink(cloneLinks interface{}) (string, error) {
	for _, link := range cloneLinks.([]interface{}) {
		if link.(map[string]interface{})["name"] == "https" {
			return link.(map[string]interface{})["href"].(string), nil
		}
	}

	return "", ErrNoHttpsCloneLink
}
