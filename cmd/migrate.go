package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github/situ2001.com/gitea-bulk-migration/client"
	"github/situ2001.com/gitea-bulk-migration/common"

	"github.com/google/go-github/v71/github"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var cliOpts = &common.CliOption{
	MigrationCliOption: common.MigrationCliOption{
		DuplicationStrategy:            common.DuplicationStrategySkip,
		DuplicationOnNonMirrorStrategy: common.DuplicationOnNonMirrorStrategySkip,
		DeletedRepoStrategy:            common.DeletedRepoStrategySkip,
	},
}

var rootCmd = &cobra.Command{
	Use:   "gitea-bulk-migrate",
	Short: "Bulk migrate repositories from GitHub to Gitea, as mirror.",
}

func init() {
	// Common options
	rootCmd.PersistentFlags().StringVar(&cliOpts.EnvFilePath, "env-file", ".env", "path to the env file where store GITEA_URL, GITEA_TOKEN, GITHUB_TOKEN")
	rootCmd.PersistentFlags().StringVar(&cliOpts.HttpProxy, "http-proxy", "", "http proxy to access GitHub")

	// Migration options
	rootCmd.PersistentFlags().StringVar(&cliOpts.MigrationCliOption.GiteaOwner, "gitea-owner", "", "the owner name of the repository after migration, can be username or org name")
	rootCmd.PersistentFlags().StringVar(&cliOpts.MigrationCliOption.TypeOfRepoBeingMigrated, "repo-type", "all", "the type of repository being migrated (all, owner, public, private, member)")
	rootCmd.PersistentFlags().BoolVar(&cliOpts.MigrationCliOption.ShouldMigrateForkedRepo, "migrate-fork-repo", false, "whether migrate forked repos from GitHub user")
	rootCmd.PersistentFlags().BoolVar(&cliOpts.MigrationCliOption.ShouldMigrateLFS, "migrate-lfs", false, "whether migrate the LFS of GitHub repo or not")
	rootCmd.PersistentFlags().BoolVar(&cliOpts.MigrationCliOption.TriggerSyncForExistingMirrorRepo, "sync-mirror-repo", false, "Should sync the repo already exists in Gitea and GitHub, after migration")

	rootCmd.PersistentFlags().Var(&cliOpts.MigrationCliOption.DuplicationStrategy, "on-duplication", "strategy to handle the repository that already exists in Gitea(as mirror) and GitHub. (skip, overwrite, abort)")
	rootCmd.PersistentFlags().Var(&cliOpts.MigrationCliOption.DuplicationOnNonMirrorStrategy, "on-duplication-non-mirror", "strategy to handle the repository that already exists in Gitea(as non-mirror) and GitHub. (skip, overwrite, abort)")
	rootCmd.PersistentFlags().Var(&cliOpts.MigrationCliOption.DeletedRepoStrategy, "on-deletion", "strategy to handle the repository that already exists in Gitea but not in GitHub. (skip, delete, abort)")

	// load env variables from .env file
	InitEnv(&InitEnvOptions{
		EnvFile: cliOpts.EnvFilePath,
		Proxy:   cliOpts.HttpProxy,
	})

	rootCmd.MarkPersistentFlagRequired("gitea-owner")

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
			repoStr += repo.FullName + ", "
		}
		log.Println(repoStr)

		// Start fetching repos from GitHub
		reposUnderGithubOwner, err := githubClient.GetAllGitHubRepoByUsername(&cliOpts.MigrationCliOption)
		if err != nil {
			log.Fatalln("Error getting GitHub repos:", err)
		}
		log.Println("Length of repos under GitHub owner:", len(reposUnderGithubOwner))
		repoStr = ""
		for _, repo := range reposUnderGithubOwner {
			repoStr += repo.GetFullName() + ", "
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
				repoStr += repo.GetFullName() + ", "
			}
			log.Println(repoStr)
		}

		diffSet := common.CompareGitHubAndGitea(reposUnderGithubOwner, reposUnderGiteaOwner)

		printDiffSet(&diffSet, &cliOpts.MigrationCliOption)

		if !promptUserConfirmation("Do you want to proceed with the migration? (yes/No): ") {
			log.Println("Migration aborted by the user.")
			return
		}

		log.Println("Start migrating the Mirror Repos Not on GitHub but on Gitea")
		for idx, giteaRepo := range diffSet.MirrorRepoNotExistOnGithub {
			log.Printf("(%d/%d) Mirror repo %s is not on GitHub", idx+1, len(diffSet.MirrorRepoNotExistOnGithub), giteaRepo.CloneURL)

			switch cliOpts.MigrationCliOption.DeletedRepoStrategy {
			case common.DeletedRepoStrategySkip:
				log.Println("Skip deleting the repo.")
			case common.DeletedRepoStrategyAbort:
				log.Fatalln("Abort migration because the repo already exists on Gitea but not on GitHub.")
			case common.DeletedRepoStrategyDelete:
				log.Println("Start deleting", giteaRepo.FullName, "on Gitea")
				giteaClient.DeleteRepo(cliOpts.MigrationCliOption.GiteaOwner, giteaRepo.Name)
				log.Println("Deleted repo on Gitea:", giteaRepo.FullName)
			}
		}
		log.Println("Completed.")

		log.Println("Start migrating the Repo name are same on both side, but is not mirror repo on Gitea")
		for idx, repos := range diffSet.RepoExistBothSideByNameButNotMirrorRepoOnGitea {
			log.Printf("(%d/%d) Has non-mirrored repo %s on Gitea, with same name as GitHub repo %s", idx+1, len(diffSet.RepoExistBothSideByNameButNotMirrorRepoOnGitea), repos.GiteaRepo.CloneURL, repos.GithubRepo.GetCloneURL())

			switch cliOpts.MigrationCliOption.DuplicationOnNonMirrorStrategy {
			case common.DuplicationOnNonMirrorStrategySkip:
				log.Println("Skip deleting the repo.")
			case common.DuplicationOnNonMirrorStrategyAbort:
				log.Fatalln("Abort migration because the repo already exists on Gitea but not on GitHub.")
			case common.DuplicationOnNonMirrorStrategyOverwrite:
				log.Println("Start deleting", repos.GiteaRepo.FullName, "on Gitea")
				giteaClient.DeleteRepo(cliOpts.MigrationCliOption.GiteaOwner, repos.GiteaRepo.Name)
				log.Println("Deleted repo on Gitea:", repos.GiteaRepo.FullName)

				log.Println("Start migrating to Gitea")
				giteaClient.MirrorGithubRepository(repos.GithubRepo, cliOpts.MigrationCliOption.GiteaOwner, envValues.GithubToken)
				log.Println("Migrated repo on Gitea:", repos.GithubRepo.GetFullName())
			}
		}
		log.Println("Completed.")

		log.Println("Start migrating the GitHub Repos Not Mirrored on Gitea")
		for idx, repo := range diffSet.GithubRepoNotMirroredOnGitea {
			log.Printf("(%d/%d) GitHub repo %s is not mirrored on Gitea", idx+1, len(diffSet.GithubRepoNotMirroredOnGitea), repo.GetCloneURL())

			log.Println("Start migrating to Gitea")
			giteaClient.MirrorGithubRepository(repo, cliOpts.MigrationCliOption.GiteaOwner, envValues.GithubToken)
			log.Println("Migrated repo on Gitea:", repo.GetFullName())
		}

		log.Println("Start migrating the Repos Mirrored on Both Sides")
		for idx, repos := range diffSet.RepoExistBothSideWithSameUrl {
			log.Printf("(%d/%d) GitHub repo %s is already mirrored on Gitea repo %s", idx+1, len(diffSet.RepoExistBothSideWithSameUrl), repos.GithubRepo.GetCloneURL(), repos.GiteaRepo.CloneURL)

			if cliOpts.MigrationCliOption.TriggerSyncForExistingMirrorRepo {
				log.Println("Start triggering sync for repo on Gitea")
				giteaClient.MirrorSync(repos.GiteaRepo.Owner.UserName, repos.GiteaRepo.Name)
				log.Println("Triggered sync for repo on Gitea:", repos.GithubRepo.GetFullName())
			}
		}
		log.Println("Completed.")
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Error executing command:", err)
	}
}

func printDiffSet(diffSet *common.MigrationDifferenceSet, migrationCliOption *common.MigrationCliOption) {
	log.Println("DiffSet Summary:")
	log.Println()

	// Create a table for each category in diffSet
	log.Printf("Strategy: %s\n", migrationCliOption.DeletedRepoStrategy)
	printTable("Mirror Repos Not on GitHub but on Gitea", []string{"Index", "Repo Name", "Clone URL"}, func(table *tablewriter.Table) {
		for idx, repo := range diffSet.MirrorRepoNotExistOnGithub {
			table.Append([]string{fmt.Sprintf("%d", idx+1), repo.FullName, repo.CloneURL})
		}
	})

	log.Printf("Strategy: %s\n", migrationCliOption.DuplicationOnNonMirrorStrategy)
	printTable("Repo name are same on both side, but is not mirror repo on Gitea", []string{"Index", "Gitea Repo", "GitHub Repo"}, func(table *tablewriter.Table) {
		for idx, repos := range diffSet.RepoExistBothSideByNameButNotMirrorRepoOnGitea {
			table.Append([]string{fmt.Sprintf("%d", idx+1), repos.GiteaRepo.FullName, repos.GithubRepo.GetFullName()})
		}
	})

	log.Printf("Will be migrated\n")
	printTable("GitHub Repos Not Mirrored on Gitea", []string{"Index", "Repo Name", "Clone URL"}, func(table *tablewriter.Table) {
		for idx, repo := range diffSet.GithubRepoNotMirroredOnGitea {
			table.Append([]string{fmt.Sprintf("%d", idx+1), repo.GetFullName(), repo.GetCloneURL()})
		}
	})

	log.Printf("Will trigger sync: %v\n", migrationCliOption.TriggerSyncForExistingMirrorRepo)
	printTable("Repos Mirrored on Both Sides", []string{"Index", "GitHub Repo", "Gitea Repo"}, func(table *tablewriter.Table) {
		for idx, repos := range diffSet.RepoExistBothSideWithSameUrl {
			table.Append([]string{fmt.Sprintf("%d", idx+1), repos.GithubRepo.GetFullName(), repos.GiteaRepo.FullName})
		}
	})
}

func printTable(title string, headers []string, fillTable func(table *tablewriter.Table)) {
	fmt.Println(title)
	table := tablewriter.NewWriter(os.Stdout) // TODO can it be logged to file?
	table.SetHeader(headers)
	fillTable(table)
	table.Render()
	fmt.Println()
}

// Function to prompt the user for confirmation
func promptUserConfirmation(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input:", err)
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes"
}
