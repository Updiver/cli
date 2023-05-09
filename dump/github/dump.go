package github

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/google/go-github/v52/github"
	"github.com/updiver/cli/flags"
	"github.com/updiver/dumper"
	"golang.org/x/oauth2"
)

type Repositories []*github.Repository

func (r Repositories) Iter(f func(repository *github.Repository) error) error {
	for _, repository := range r {
		if err := f(repository); err != nil {
			return err
		}
	}

	return nil
}

type DumpOptions struct {
	Username    string
	Token       string
	Destination string
	CloneMode   string
}

func Dump(opts *DumpOptions) error {
	ghClient := AuthenticatedClient(opts.Token)
	repos, err := UserRepositories(ghClient)
	if err != nil {
		return fmt.Errorf("get repositories: %w", err)
	}
	err = repos.Iter(func(repo *github.Repository) error {
		return dumpRepo(repo, opts)
	})

	return err
}

func dumpRepo(repo *github.Repository, opts *DumpOptions) error {
	logger.Printf("org [%s] | repo [%s]\n", *repo.Owner.Login, *repo.Name)
	if repo.CloneURL == nil {
		logger.Printf("skipping repo [%s] as it has no clone url\n", *repo.Name)
		return nil
	}

	fullDestFolder := path.Join(opts.Destination, *repo.Owner.Login, *repo.Name)
	dpr := dumper.New()
	dumpRepoOpts := &dumper.DumpRepositoryOptions{
		RepositoryURL: *repo.CloneURL,
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
	_, err := dpr.DumpRepository(dumpRepoOpts)
	if err != nil {
		return fmt.Errorf("dump repository: %w", err)
	}

	logger.Printf("repo [%s] dumped", *repo.Name)
	return nil
}

func UserRepositories(ghClient *github.Client) (Repositories, error) {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := ghClient.Repositories.List(context.Background(), "", opts)
		if err != nil {
			return nil, fmt.Errorf("get repositories list: %w", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func AuthenticatedClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
