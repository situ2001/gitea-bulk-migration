package common

type CliOption struct {
	// path where .env file is located, the default is ".env"
	// and ensure there is url to Gitea server, token to access Gitea server
	// and token to access GitHub
	EnvFilePath string

	// the URL to proxy server, which is used to access GitHub
	HttpProxy string

	// TODO the path to the file where the output log will be saved
	// outputLogFilePath string

	// Options for migration
	MigrationCliOption MigrationCliOption
}

type MigrationCliOption struct {
	// The owner name of the repository after migration
	GiteaOwner string

	// same as value of `(github.RepositoryListByAuthenticatedUserOptions).Type`
	TypeOfRepoBeingMigrated string

	ShouldMirrorRepo bool

	ShouldMigrateForkedRepo bool
	ShouldMigrateLFS        bool

	// Choose the strategy to handle the repository that already exists in Gitea and GitHub
	// string will be: "skip", "overwrite", "abort"
	// Default: skip
	DuplicationStrategy string

	// Choose the strategy to handle the repository that already exists in Gitea but not in GitHub
	// string will be: "skip", "delete", "abort"
	// Default: skip
	DeletedRepoStrategy string
}
