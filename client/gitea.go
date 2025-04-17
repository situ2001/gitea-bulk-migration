package client

import (
	"github/situ2001.com/gitea-bulk-migration/common"

	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/v71/github"
)

type GiteaClient struct {
	*gitea.Client

	migrationOption *common.MigrationCliOption
}

func NewGiteaClient(url, token string, migrationOption *common.MigrationCliOption) (*GiteaClient, error) {
	client, err := gitea.NewClient(url, gitea.SetToken(token))
	if err != nil {
		return nil, err
	}
	return &GiteaClient{client, migrationOption}, nil
}

func (c *GiteaClient) ListUserReposAll(giteaUser string) ([]*gitea.Repository, error) {
	listOptions := gitea.ListOptions{
		Page: 1,
	}
	reposUnderGiteaOwner := make([]*gitea.Repository, 0)

	for {
		repos, _, err := c.ListUserRepos(giteaUser, gitea.ListReposOptions{
			ListOptions: listOptions,
		})

		if err != nil {
			return nil, err
		}

		if len(repos) == 0 {
			break
		}

		// TODO seems not working... maybe a bug?
		// if resp.LastPage == listOptions.Page {
		// 	break
		// }

		reposUnderGiteaOwner = append(reposUnderGiteaOwner, repos...)

		listOptions.Page++
	}

	return reposUnderGiteaOwner, nil
}

// It mirrors a GitHub repository to Gitea.
func (c *GiteaClient) MirrorGithubRepository(repo *github.Repository, giteaUser string, githubToken string) (*gitea.Repository, error) {
	migratedRepo, _, err := c.MigrateRepo(gitea.MigrateRepoOption{
		Service:   gitea.GitServiceGithub,
		RepoName:  repo.GetName(),
		RepoOwner: giteaUser,
		CloneAddr: repo.GetCloneURL(),
		AuthToken: githubToken,
		Mirror:    true,
		Private:   *repo.Private,
		LFS:       c.migrationOption.ShouldMigrateLFS,
	})

	if err != nil {
		return nil, err
	}

	return migratedRepo, nil
}
