package bitbucket

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/stretchr/testify/require"
	"github.com/updiver/cli/dump"
	"github.com/updiver/dumper"
)

var (
	testRepositoryURL        = "https://bitbucket.org/updiver/test-repository.git"
	destinationRepositoryDir = "repository-clone-example"
)

func TestDumpRepo_DefaultBranchDump(t *testing.T) {
	tempDir := os.TempDir()
	mockRepos := []*bitbucket.Repository{
		{
			Slug: "updiver",
			Name: "test-repository",
			Links: map[string]interface{}{
				"clone": interface{}(
					[]interface{}{
						interface{}(
							map[string]interface{}{
								"name": "https",
								"href": testRepositoryURL,
							},
						),
					},
				),
			},
		},
	}

	for _, repo := range mockRepos {
		prefix, err := dump.GenerateRandomNumber()
		require.NoError(t, err, "generate random prefix")
		fullDestinationPath := path.Join(filepath.Clean(tempDir), destinationRepositoryDir, prefix)
		mockDumpOpts := &DumpOptions{
			Username:    os.Getenv("BITBUCKET_USERNAME"),
			Token:       os.Getenv("BITBUCKET_TOKEN"),
			Destination: fullDestinationPath,
			CloneMode:   "default-branch",
		}
		err = dumpRepo(repo, mockDumpOpts)
		defer os.RemoveAll(fullDestinationPath)

		clonedRepoPath := path.Join(fullDestinationPath, repo.Slug, repo.Name)
		require.NoError(t, err, "expected no error, got %v", err)

		repository, err := dumper.Repository(clonedRepoPath)
		require.NoError(t, err, "expected no error, got %v", err)
		require.NotNil(t, repository, "expected repository to be not nil")

		fileContent, err := os.Open(path.Join(clonedRepoPath, "test-regular-file.txt"))
		require.NoError(t, err, "open file")

		txt, err := io.ReadAll(fileContent)
		require.NoError(t, err, "read file content")

		require.Equal(t, "Test regular file content", string(txt), "expect to have proper file content")

		refIter, err := repository.Branches()
		require.NoError(t, err, "get branches iterator")

		branches := make([]string, 0)
		refIter.ForEach(func(ref *plumbing.Reference) error {
			branches = append(branches, ref.Name().Short())
			return nil
		})
		require.Len(t, branches, 1, "expect to have only one branch")
		require.Equal(t, "main", branches[0], "expect to have proper branch name")
	}
}

func TestDumpRepo_AllBranchesDump(t *testing.T) {
	tempDir := os.TempDir()
	mockRepos := []*bitbucket.Repository{
		{
			Slug: "updiver",
			Name: "test-repository",
			Links: map[string]interface{}{
				"clone": interface{}(
					[]interface{}{
						interface{}(
							map[string]interface{}{
								"name": "https",
								"href": testRepositoryURL,
							},
						),
					},
				),
			},
		},
	}

	for _, repo := range mockRepos {
		prefix, err := dump.GenerateRandomNumber()
		require.NoError(t, err, "generate random prefix")
		fullDestinationPath := path.Join(filepath.Clean(tempDir), destinationRepositoryDir, prefix)
		mockDumpOpts := &DumpOptions{
			Username:    os.Getenv("BITBUCKET_USERNAME"),
			Token:       os.Getenv("BITBUCKET_TOKEN"),
			Destination: fullDestinationPath,
			CloneMode:   "all-branches",
		}
		err = dumpRepo(repo, mockDumpOpts)
		defer os.RemoveAll(fullDestinationPath)

		clonedRepoPath := path.Join(fullDestinationPath, repo.Slug, repo.Name)
		require.NoError(t, err, "expected no error, got %v", err)

		repository, err := dumper.Repository(clonedRepoPath)
		require.NoError(t, err, "expected no error, got %v", err)
		require.NotNil(t, repository, "expected repository to be not nil")

		fileContent, err := os.Open(path.Join(clonedRepoPath, "test-regular-file.txt"))
		require.NoError(t, err, "open file")

		txt, err := io.ReadAll(fileContent)
		require.NoError(t, err, "read file content")

		require.Equal(t, "Test regular file content", string(txt), "expect to have proper file content")

		refIter, err := repository.Branches()
		require.NoError(t, err, "get branches iterator")

		branches := make([]string, 0)
		refIter.ForEach(func(ref *plumbing.Reference) error {
			branches = append(branches, ref.Name().Short())
			return nil
		})
		require.Len(t, branches, 3, "expect to have three branches")

		expectedBranches := []string{
			"feat/test-regular-file-first-change",
			"feat/test-regular-file-second-change",
			"main",
		}
		require.ElementsMatch(t, expectedBranches, branches, "expect to have proper branch names")
	}

}
