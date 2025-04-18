package client

import (
	"context"

	"github/situ2001.com/gitea-bulk-migration/common"

	"github.com/google/go-github/v71/github"
)

type GitHubClient struct {
	*github.Client
}

func NewGitHubClient(client *github.Client) *GitHubClient {
	return &GitHubClient{
		Client: client,
	}
}

func (c *GitHubClient) GetAllGitHubRepoByUsername(option *common.MigrationCliOption) ([]*github.Repository, error) {
	repos := make([]*github.Repository, 0)
	listOption := &github.ListOptions{
		PerPage: 50,
		Page:    1,
	}

	for {
		sourcesRepo, resp, err := c.Repositories.ListByAuthenticatedUser(context.Background(), &github.RepositoryListByAuthenticatedUserOptions{
			Type:        "owner", // Here, we only get the repositories that the user owns to avoid same-name repositories conflict
			ListOptions: *listOption,
		})

		if err != nil {
			return nil, err
		}

		if len(sourcesRepo) == 0 {
			break
		}

		repos = append(repos, sourcesRepo...)

		if resp.LastPage == listOption.Page {
			break
		}

		listOption.Page++
	}

	return repos, nil
}
