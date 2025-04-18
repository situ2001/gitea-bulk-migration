package common

type CliOption struct {
	// path where .env file is located, the default is ".env"
	// and ensure there is url to Gitea server, token to access Gitea server
	// and token to access GitHub
	EnvFilePath string

	// the URL to proxy server, which is used to access GitHub
	HttpProxy string

	// TODO not supported yet, currently, log will be appended to the file migration.log
	// outputLogFilePath string

	// Options for migration
	MigrationCliOption MigrationCliOption
}

type MigrationCliOption struct {
	// The owner name of the repository after migration
	GiteaOwner string

	ShouldMirrorRepo bool

	ShouldMigrateForkedRepo bool
	ShouldMigrateLFS        bool

	// Choose the strategy to handle the repository that already exists in Gitea (mirror repo) and GitHub
	// string will be: "skip", "overwrite", "abort"
	// Default: skip
	DuplicationStrategy DuplicationStrategyType

	// TODO 因为尽可能只处理 mirror 的 repo，但如果有同名非 mirror 的 repo 在 Gitea 上，那么在 migrate 的时候会报错同名，因此我们必须要考虑 Gitea 上的这类 repo
	// TODO command line options
	// Choose the strategy to handle the repository that already exists in Gitea (non-mirror repo) and GitHub
	// string will be: "skip", "overwrite", "abort"
	DuplicationOnNonMirrorStrategy DuplicationOnNonMirrorStrategyType

	// Choose the strategy to handle the repository that already exists in Gitea but not in GitHub
	// string will be: "skip", "delete", "abort"
	// Default: skip
	DeletedRepoStrategy DeletedRepoStrategyType

	// Default: false
	TriggerSyncForExistingMirrorRepo bool
}
