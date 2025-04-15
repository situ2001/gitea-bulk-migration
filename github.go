package main

import (
	"context"
	"fmt"

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

func (c *GitHubClient) GetAllGitHubRepoByUsername(username string) ([]*github.Repository, error) {
	repos := make([]*github.Repository, 0)
	listOption := &github.ListOptions{
		PerPage: 50,
		Page:    1,
	}

	for {
		sourcesRepo, resp, err := c.Repositories.ListByAuthenticatedUser(context.Background(), &github.RepositoryListByAuthenticatedUserOptions{
			Type:        "all",
			Sort:        "full_name",
			ListOptions: *listOption,
		})

		if err != nil {
			fmt.Println("Error getting repositories:", err)
			return nil, err
		}

		if len(sourcesRepo) == 0 {
			break
		}

		fmt.Println("Page:", listOption.Page)

		repos = append(repos, sourcesRepo...)

		if resp.LastPage == listOption.Page {
			break
		}

		listOption.Page++
	}

	fmt.Println("Total repositories:", len(repos))

	return repos, nil
}
