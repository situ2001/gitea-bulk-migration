package main

import (
	"fmt"
	"os"

	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/v71/github"
	"github.com/joho/godotenv"
)

func main() {
	// Read gitea URL and token from .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	giteaUrl := os.Getenv("GITEA_URL")
	giteaToken := os.Getenv("GITEA_TOKEN")

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
	fmt.Println("Hello:")
	fmt.Println("Username:", userInfo.UserName)
	fmt.Println("Email:", userInfo.Email)

	// GitHub Client
	client := github.NewClient(nil)
	println(client.BaseURL.String())
}
