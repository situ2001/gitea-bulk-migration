package cmd

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

func GetEnvValues() *EnvValues {
	giteaUrl := os.Getenv("GITEA_URL")
	giteaToken := os.Getenv("GITEA_TOKEN")
	githubToken := os.Getenv("GITHUB_TOKEN")

	// add to map, and return
	envValues := EnvValues{
		GiteaUrl:    giteaUrl,
		GiteaToken:  giteaToken,
		GithubToken: githubToken,
	}

	return &envValues
}

func EnsureEnvValues() {
	envValues := GetEnvValues()
	if envValues.GiteaUrl == "" {
		fmt.Println("GITEA_URL is not set")
		os.Exit(1)
	}
	if envValues.GiteaToken == "" {
		fmt.Println("GITEA_TOKEN is not set")
		os.Exit(1)
	}
	if envValues.GithubToken == "" {
		fmt.Println("GITHUB_TOKEN is not set")
		os.Exit(1)
	}
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
		os.Exit(1)
	}

	if options.Proxy != "" {
		os.Setenv("HTTP_PROXY", options.Proxy)
		os.Setenv("HTTPS_PROXY", options.Proxy)
	}
}
