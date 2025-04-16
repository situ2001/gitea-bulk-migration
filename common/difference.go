package common

import (
	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/v71/github"
)

// Base on github repo. We use GitHub as the source of truth.
type MigrationDifferenceSet struct {
	// The repo should be deleted on gitea, since it is not on github
	MirrorRepoNotExistOnGithub []*gitea.Repository

	// The repo should be migrated to gitea, since it is not on gitea
	GithubRepoNotMirroredOnGitea []*github.Repository

	// TODO should perform sync for existing mirror repos after migration!!!
	// Note: we compare existence of mirror repos by URL
	// Corner case: github A -> B, B -> C, gitea: A B
	// But since B share the same URL with B renamed from A on github. We just need to perform a sync on B
	// No need to delete B on gitea.
	// So after migration, we will have B(not changed, just synced again) & C(newly created) on gitea
	RepoExistBothSideWithSameUrl []struct {
		GiteaRepo  *gitea.Repository
		GithubRepo *github.Repository
	}

	// The repo exists on gitea when compared with github by name
	// But it is not a mirror repo, so it may not a migrated repo
	//
	// IMPORTANT to add it to set, since there may be a repo that is not migrated from github
	RepoExistBothSideByNameButNotMirrorRepoOnGitea []struct {
		GiteaRepo  *gitea.Repository
		GithubRepo *github.Repository
	}
}

func CompareGitHubAndGitea(githubRepos []*github.Repository, giteaRepos []*gitea.Repository) MigrationDifferenceSet {
	diffSet := MigrationDifferenceSet{}

	githubReposMapByRepoName := make(map[string]*github.Repository)
	githubReposMapByOriginalURL := make(map[string]*github.Repository)
	for _, repo := range githubRepos {
		githubReposMapByRepoName[repo.GetName()] = repo
		githubReposMapByOriginalURL[repo.GetCloneURL()] = repo
	}

	for _, giteaRepo := range giteaRepos {
		if !giteaRepo.Mirror {
			githubRepoByName, existOnGithubByName := githubReposMapByRepoName[giteaRepo.Name]

			// if the repo exists on gitea when compared with github by name and it is not a mirror repo
			if existOnGithubByName {
				diffSet.RepoExistBothSideByNameButNotMirrorRepoOnGitea = append(diffSet.RepoExistBothSideByNameButNotMirrorRepoOnGitea, struct {
					GiteaRepo  *gitea.Repository
					GithubRepo *github.Repository
				}{
					GiteaRepo:  giteaRepo,
					GithubRepo: githubRepoByName,
				})

				// delete the repo from github repos from url map, since we already processed it
				delete(githubReposMapByOriginalURL, *githubRepoByName.CloneURL)
			}

			continue
		}

		// From now on, we only care about mirror repos
		githubRepo, existsOnGithubByOriginURL := githubReposMapByOriginalURL[giteaRepo.OriginalURL]

		if !existsOnGithubByOriginURL {
			diffSet.MirrorRepoNotExistOnGithub = append(diffSet.MirrorRepoNotExistOnGithub, giteaRepo)
			continue
		}

		diffSet.RepoExistBothSideWithSameUrl = append(diffSet.RepoExistBothSideWithSameUrl, struct {
			GiteaRepo  *gitea.Repository
			GithubRepo *github.Repository
		}{
			GiteaRepo:  giteaRepo,
			GithubRepo: githubRepo,
		})

		delete(githubReposMapByOriginalURL, giteaRepo.OriginalURL)
	}

	for _, githubRepo := range githubReposMapByOriginalURL {
		diffSet.GithubRepoNotMirroredOnGitea = append(diffSet.GithubRepoNotMirroredOnGitea, githubRepo)
	}

	return diffSet
}
