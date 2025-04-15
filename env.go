package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvValues struct {
	GiteaUrl    string
	GiteaToken  string
	GithubToken string
}

func GetEnvValues() EnvValues {
	giteaUrl := os.Getenv("GITEA_URL")
	giteaToken := os.Getenv("GITEA_TOKEN")
	githubToken := os.Getenv("GITHUB_TOKEN")

	// add to map, and return
	envValues := EnvValues{
		GiteaUrl:    giteaUrl,
		GiteaToken:  giteaToken,
		GithubToken: githubToken,
	}

	return envValues
}

type InitEnvOptions struct {
	EnvFile string
	Proxy   string
}

//optional proxy string
func InitEnv(options *InitEnvOptions) {
	// Read gitea URL and token from .env
	if options.EnvFile == "" {
		options.EnvFile = ".env"
	}

	err := godotenv.Load(options.EnvFile)
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	if options.Proxy != "" {
		os.Setenv("HTTP_PROXY", options.Proxy)
		os.Setenv("HTTPS_PROXY", options.Proxy)
	}
}
