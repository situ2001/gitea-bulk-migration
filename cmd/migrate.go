package cmd

import (
	"fmt"
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
			fmt.Println("Error creating Gitea client:", err)
			os.Exit(1)
		}

		fmt.Println(githubClient)
		fmt.Println(giteaClient)

		reposUnderGiteaOwner, err := giteaClient.ListUserReposAll(cliOpts.MigrationCliOption.GiteaOwner)
		if err != nil {
			fmt.Println("Error getting Gitea repos:", err)
			os.Exit(1)
		}
		fmt.Println("Length of repos under Gitea owner:", len(reposUnderGiteaOwner))

		return

		repos, err := githubClient.GetAllGitHubRepoByUsername(&cliOpts.MigrationCliOption)
		if err != nil {
			fmt.Println("Error getting GitHub repos:", err)
			os.Exit(1)
		}

		// filtering repos
		if !cliOpts.MigrationCliOption.ShouldMigrateForkedRepo {
			// filter out forked repos
			filteredRepos := make([]*github.Repository, 0)
			for _, repo := range repos {
				if !repo.GetFork() {
					filteredRepos = append(filteredRepos, repo)
				}
			}
			repos = filteredRepos
		}

		// TODO Compare repos with existing Gitea repos, find deleted repos in GitHub-side

	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
