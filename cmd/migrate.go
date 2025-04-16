package cmd

import (
	"fmt"
	"log"
	"os"

	"github/situ2001.com/gitea-bulk-migration/client"
	"github/situ2001.com/gitea-bulk-migration/common"

	"github.com/google/go-github/v71/github"
	"github.com/spf13/cobra"
)

var cliOpts = &common.CliOption{}

var rootCmd = &cobra.Command{
	Use:   "gitea-bulk-migrate",
	Short: "Bulk migrate repositories from GitHub to Gitea or managed the migrated repositories.",
}

func init() {
	// Common options
	rootCmd.PersistentFlags().StringVar(&cliOpts.EnvFilePath, "env-file", ".env", "path to the env file where store GITEA_URL, GITEA_TOKEN, GITHUB_TOKEN")
	rootCmd.PersistentFlags().StringVar(&cliOpts.HttpProxy, "http-proxy", "", "http proxy to access GitHub")

	// Migration options
	rootCmd.PersistentFlags().StringVar(&cliOpts.MigrationCliOption.GiteaOwner, "gitea-owner", "", "the owner name of the repository after migration, can be username or org name")
	rootCmd.PersistentFlags().StringVar(&cliOpts.MigrationCliOption.TypeOfRepoBeingMigrated, "repo-type", "all", "the type of repository being migrated (all, owner, public, private, member)")
	rootCmd.PersistentFlags().BoolVar(&cliOpts.MigrationCliOption.ShouldMigrateForkedRepo, "migrate-fork-repo", true, "whether migrate forked repos from GitHub user")
	rootCmd.PersistentFlags().BoolVar(&cliOpts.MigrationCliOption.ShouldMigrateLFS, "migrate-lfs", false, "whether migrate the LFS of GitHub repo or not")
	rootCmd.PersistentFlags().StringVar(&cliOpts.MigrationCliOption.DuplicationStrategy, "on-duplication", "skip", "strategy to handle the repository that already exists in Gitea and GitHub. (skip, overwrite, abort)")
	rootCmd.PersistentFlags().StringVar(&cliOpts.MigrationCliOption.DeletedRepoStrategy, "on-deletion", "skip", "strategy to handle the repository that already exists in Gitea but not in GitHub. (skip, delete, abort)")

	// load env variables from .env file
	InitEnv(&InitEnvOptions{
		EnvFile: cliOpts.EnvFilePath,
		Proxy:   cliOpts.HttpProxy,
	})

	rootCmd.MarkFlagRequired("gitea-owner")

	EnsureEnvValues()
}

func Execute() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		envValues := GetEnvValues()

		githubClient := client.NewGitHubClient(github.NewClientWithEnvProxy().WithAuthToken(envValues.GithubToken))

		giteaClient, err := client.NewGiteaClient(envValues.GiteaUrl, envValues.GiteaToken, &cliOpts.MigrationCliOption)
		if err != nil {
			log.Fatalln("Error creating Gitea client:", err)
		}

		// Start fetching repos from Gitea
		reposUnderGiteaOwner, err := giteaClient.ListUserReposAll(cliOpts.MigrationCliOption.GiteaOwner)
		if err != nil {
			log.Fatalln("Error getting Gitea repos:", err)
		}
		log.Println("Length of repos under Gitea owner:", len(reposUnderGiteaOwner))
		repoStr := ""
		for _, repo := range reposUnderGiteaOwner {
			repoStr += repo.Name + ", "
		}
		log.Println(repoStr)

		// Start fetching repos from GitHub
		reposUnderGithubOwner, err := githubClient.GetAllGitHubRepoByUsername(&cliOpts.MigrationCliOption)
		if err != nil {
			log.Println("Error getting GitHub repos:", err)
		}
		log.Println("Length of repos under GitHub owner:", len(reposUnderGithubOwner))
		repoStr = ""
		for _, repo := range reposUnderGithubOwner {
			repoStr += repo.GetName() + ", "
		}
		log.Println(repoStr)

		// filtering repos
		if !cliOpts.MigrationCliOption.ShouldMigrateForkedRepo {
			// filter out forked repos
			filteredRepos := make([]*github.Repository, 0)
			for _, repo := range reposUnderGithubOwner {
				if !repo.GetFork() {
					filteredRepos = append(filteredRepos, repo)
				}
			}
			reposUnderGithubOwner = filteredRepos

			log.Println("Length of repos under GitHub owner after filtering forked repos:", len(reposUnderGithubOwner))
			repoStr := ""
			for _, repo := range reposUnderGithubOwner {
				repoStr += repo.GetName() + ", "
			}
			log.Println(repoStr)
		}

		diffSet := common.CompareGitHubAndGitea(reposUnderGithubOwner, reposUnderGiteaOwner)

		// TODO print the diffSet in a more readable way
		// TODO before starting, prompt the user to confirm the migration

		log.Println("Start migrating the repos that exist on Gitea but not on GitHub")
		for idx, giteaRepo := range diffSet.MirrorRepoNotExistOnGithub {
			log.Printf("(%d/%d) Mirror repo %s is not on GitHub", idx+1, len(diffSet.MirrorRepoNotExistOnGithub), giteaRepo.CloneURL)

			// TODO handle strategy for this case
			log.Println("Start deleting", giteaRepo.Name, "on Gitea")
			giteaClient.DeleteRepo(cliOpts.MigrationCliOption.GiteaOwner, giteaRepo.Name)
			log.Println("Deleted repo on Gitea:", giteaRepo.Name)
		}
		log.Println("Completed.")

		log.Println("Start migrating the RepoExistBothSideByNameButNotMirrorRepoOnGitea")
		for idx, repos := range diffSet.RepoExistBothSideByNameButNotMirrorRepoOnGitea {
			log.Printf("(%d/%d) Has non-mirrored repo %s on Gitea, with same name as GitHub repo %s", idx+1, len(diffSet.RepoExistBothSideByNameButNotMirrorRepoOnGitea), repos.GiteaRepo.CloneURL, repos.GithubRepo.GetCloneURL())

			// TODO handle strategy for this case
			log.Println("Start deleting", repos.GiteaRepo.Name, "on Gitea")
			giteaClient.DeleteRepo(cliOpts.MigrationCliOption.GiteaOwner, repos.GiteaRepo.Name)
			log.Println("Deleted repo on Gitea:", repos.GiteaRepo.Name)

			log.Println("Start migrating to Gitea")
			giteaClient.MirrorGithubRepository(repos.GithubRepo, cliOpts.MigrationCliOption.GiteaOwner, envValues.GithubToken)
			log.Println("Migrated repo on Gitea.")
		}
		log.Println("Completed.")

		log.Println("Start migrating the RepoExistBothSideWithSameUrl")
		for idx, repos := range diffSet.RepoExistBothSideWithSameUrl {
			log.Printf("(%d/%d) GitHub repo %s is already mirrored on Gitea repo %s", idx+1, len(diffSet.RepoExistBothSideWithSameUrl), repos.GithubRepo.GetCloneURL(), repos.GiteaRepo.CloneURL)

			// TODO handle strategy for this case
			// log.Println("Start syncing", repos.GithubRepo.GetName(), "to Gitea")
			// giteaClient.MirrorSync(repos.GiteaRepo.Owner.UserName, repos.GiteaRepo.Name)
			// log.Println("Synced repo on Gitea:", repos.GithubRepo.GetName())
		}
		log.Println("Completed.")

		log.Println("Start migrating the GitHub repo that is not mirrored on Gitea")
		for idx, repo := range diffSet.GithubRepoNotMirroredOnGitea {
			log.Printf("(%d/%d) GitHub repo %s is not mirrored on Gitea", idx+1, len(diffSet.GithubRepoNotMirroredOnGitea), repo.GetCloneURL())

			log.Println("Start migrating to Gitea")
			giteaClient.MirrorGithubRepository(repo, cliOpts.MigrationCliOption.GiteaOwner, envValues.GithubToken)
			log.Println("Migrated repo on Gitea:", repo.GetName())
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
