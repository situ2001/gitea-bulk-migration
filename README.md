# gitea-bulk-migration

A command-line tool mainly for bulk migrating repositories from GitHub to Gitea as mirrors.

## Overview

This tool helps you migrate all GitHub repositories in a specified organization or user account to a Gitea instance as mirrors.

- Mirror GitHub repositories to Gitea.
- CLI can be executed in multiple runs, without needing to re-migrate already migrated repositories.
- Handle existing repositories on Gitea with configurable strategies (skip, abort, delete, overwrite) in a robust way.
- Sync existing mirror repositories, if there are any.
- Filter repositories by type (all, public, private, etc.)
- Include or exclude forked repositories

## Prerequisites

- Go 1.24 or higher
- GitHub personal access token with appropriate permissions (repository, read) // TODO 检查一下
- Gitea instance with admin access and API token (repository, read-write) // TODO 检查一下

## Installation

TODO

## Usage

First, create a `.env` file in the project root with the following variables:

```
GITEA_URL=https://your-gitea-instance.com
GITEA_TOKEN=your_gitea_api_token
GITHUB_TOKEN=your_github_personal_access_token
```

## Options

Run the following command to see all available options:

```shell
gitea-bulk-migration --help
```

## Example

TODO

## How it compares GitHub and Gitea repo

TODO

## TODO

- [ ]  Bulk manage the mirrored repositories on Gitea, for example, change mirror sync interval, etc.
- [ ]  Add more git source repositories, such as GitLab, etc.
