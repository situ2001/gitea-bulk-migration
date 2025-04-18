# gitea-bulk-migration

A command-line tool for bulk migrating GitHub repositories to Gitea as mirrors, with support for incremental syncs and configurable migration strategies

## Overview

This tool helps you migrate all GitHub repositories in a specified organization or user account to a Gitea instance as mirrors.

- Mirror GitHub repositories to Gitea.
- CLI can be executed in multiple runs to migrate repositories in batches and reflect the latest state of the GitHub repositories.
- Handle existing repositories on Gitea with configurable strategies (skip, abort, delete, overwrite) in a robust way.
- Sync existing mirror repositories, if there are any.
- Include or exclude forked repositories
- HTTP proxy support for GitHub API requests.

## Prerequisites

- Go 1.24 or higher
- GitHub personal access token with appropriate permissions: *All repositories + metadata read-only access*
- Gitea API token with appropriate permissions: *All repositories + repo write access + user read-only access*

## Installation

You can install the tool by cloning the repository and building it with Go:

> For Linux and Windows, you can also download the pre-built binary from the [releases page](https://github.com/situ2001/gitea-bulk-migration/releases).

```shell
# Clone the repository
git clone https://github.com/situ2001/gitea-bulk-migration.git
cd gitea-bulk-migration
# Build the project
go build -o gitea-bulk-migration
# Execute the binary
./gitea-bulk-migration 
```

## Usage

First, create a `.env` file in the project root with the following variables:

```
GITEA_URL=https://your-gitea-instance.com
GITEA_TOKEN=your_gitea_api_token
GITHUB_TOKEN=your_github_personal_access_token
```

Then, run the tool with the following command

```shell
# Migrate all repositories from a GitHub user account to Gitea user or organization. Suitable for first-time migration.
# You can specify --migrate-fork-repo to include forked repositories.
# You can also specify --http-proxy to set a proxy for the HTTP requests to GitHub.
./gitea-bulk-migration --gitea-owner <user or org name>

# If you want a second run, you can use the additional flags to handle existing repositories on Gitea. 
# For more detail of the options, please refer to the "How it compares GitHub and Gitea repo" section.
./gitea-bulk-migration --gitea-owner <user or org name> \
  --on-duplication-non-mirror overwrite \
  --on-deletion delete \
  --on-duplication overwrite \
  --sync-mirror-repo
```

## Options

Run the following command to see all available options:

```shell
gitea-bulk-migration --help
```

## How it compares GitHub and Gitea repo

This tool compares repositories between GitHub and Gitea. For each repo under a owner of Gitea:

1. Check if the repo is a mirror
   1. If not, check if GitHub-side has a repo with the same name
      1. If has, append the repo to the list of `RepoExistBothSideByNameButNotMirrorRepoOnGitea`
   2. If yes, check if GitHub-side has a repo with the same name OriginURL
      1. If has, append the repo to the list of `RepoExistBothSideWithSameUrl`
      2. If not, append the repo to the list of `MirrorRepoNotExistOnGithub`
2. Finally, push remaining GitHub repos to the list of `GithubRepoNotMirroredOnGitea`

> To avoid conflict, should notice that GitHub repos being migrated should be owned by only one owner, the Gitea's too. Since the tool will detect the repo by only the repo name without owner name.

You can specify the strategy to handle the cases below:

| Repo category after comparison                   | cmd arg                       | Value                     |
| ------------------------------------------------ | ----------------------------- | ------------------------- |
| `RepoExistBothSideByNameButNotMirrorRepoOnGitea` | `--on-duplication-non-mirror` | skip, overwrite, abort    |
| `RepoExistBothSideWithSameUrl`                   | `--on-duplication`            | skip, overwrite, abort    |
| `RepoExistBothSideWithSameUrl`                   | `--sync-mirror-repo`          | true, false               |
| `MirrorRepoNotExistOnGithub`                     | `--on-deletion`               | skip, delete, abort       |
| `GithubRepoNotMirroredOnGitea`                   | none                          | None (migrate by Default) |

## TODO

- [ ] Support migrating repositories for Organization on GitHub.
- [ ] Selectively migrate repositories based on the repository name.
- [ ] Bulk manage the mirrored repositories on Gitea, for example, change mirror sync interval, etc.
- [ ] Add more git source repositories, such as GitLab, etc.
