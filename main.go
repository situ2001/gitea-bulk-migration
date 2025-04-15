package main

import (
	"encoding/json"
	"fmt"
	"os"

	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/v71/github"
)

func main() {
	InitEnv(&InitEnvOptions{
		EnvFile: ".env",
		Proxy:   "http://192.168.114.3:7890",
	})

	envValues := GetEnvValues()
	giteaUrl := envValues.GiteaUrl
	giteaToken := envValues.GiteaToken
	githubToken := envValues.GithubToken

	// Gitea Client
	giteaClient, err := gitea.NewClient(giteaUrl, gitea.SetToken(giteaToken))
	if err != nil {
		fmt.Println("Error creating Gitea client:", err)
		return
	}
	// print Gitea client base URL
	userInfo, _, err := giteaClient.GetMyUserInfo()
	if err != nil {
		fmt.Println("Error getting user info:", err)
		return
	}

	// print username and email
	fmt.Println("Hello, Gitea:")
	fmt.Println("Username:", userInfo.UserName)
	fmt.Println("Email:", userInfo.Email)

	// GitHub Client
	githubClient := NewGitHubClient(github.NewClientWithEnvProxy().WithAuthToken(githubToken))
	println(githubClient.BaseURL.String())

	// print username and email from GitHub
	fmt.Printf("Hello, GitHub:")

	repos, err := githubClient.GetAllGitHubRepoByUsername("situ2001")
	if err != nil {
		fmt.Println("GetAllOwnedGitHubRepoByUsername error", err)
		os.Exit(1)
	}

	// get all non-forked repositories
	nonForkedRepos := make([]*github.Repository, 0)
	for _, repo := range repos {
		if !*repo.Fork {
			nonForkedRepos = append(nonForkedRepos, repo)
		}
	}
	println("Total non-forked repositories:", len(nonForkedRepos))

	// Serialize the nonForkedRepos to JSON, and save to ./repo-github.json
	file, err := os.Create("./repo-github.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(nonForkedRepos); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Println("Non-forked repositories saved to ./repo-github.json")

}
