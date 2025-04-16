package cmd

import (
	"log"
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
		log.Fatalln("GITEA_URL is not set")
	}
	if envValues.GiteaToken == "" {
		log.Fatalln("GITEA_TOKEN is not set")
	}
	if envValues.GithubToken == "" {
		log.Fatalln("GITHUB_TOKEN is not set")
	}
}

type InitEnvOptions struct {
	EnvFile string
	Proxy   string
}

// optional proxy string
func InitEnv(options *InitEnvOptions) {
	// Read gitea URL and token from .env
	if options.EnvFile == "" {
		options.EnvFile = ".env"
	}

	err := godotenv.Load(options.EnvFile)
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	if options.Proxy != "" {
		os.Setenv("HTTP_PROXY", options.Proxy)
		os.Setenv("HTTPS_PROXY", options.Proxy)
	}
}
