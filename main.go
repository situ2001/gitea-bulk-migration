package main

import (
	"github/situ2001.com/gitea-bulk-migration/cmd"
)

func main() {
	cmd.Execute()

	// envValues := cmd.GetEnvValues()
	// giteaUrl := envValues.GiteaUrl
	// giteaToken := envValues.GiteaToken
	// githubToken := envValues.GithubToken

	// // Gitea Client
	// giteaClient, err := client.NewGiteaClient(giteaUrl, giteaToken)
	// if err != nil {
	// 	fmt.Println("Error creating Gitea client:", err)
	// 	return
	// }
	// // print Gitea client base URL
	// userInfo, _, err := giteaClient.GetMyUserInfo()
	// if err != nil {
	// 	fmt.Println("Error getting user info:", err)
	// 	return
	// }

	// // print username and email
	// fmt.Println("Hello, Gitea:")
	// fmt.Println("Username:", userInfo.UserName)
	// fmt.Println("Email:", userInfo.Email)

	// // GitHub Client
	// githubClient := client.NewGitHubClient(github.NewClientWithEnvProxy().WithAuthToken(githubToken))
	// println(githubClient.BaseURL.String())

	// repos, err := githubClient.GetAllGitHubRepoByUsername("situ2001")
	// if err != nil {
	// 	fmt.Println("GetAllOwnedGitHubRepoByUsername error", err)
	// 	os.Exit(1)
	// }

	// // get all non-forked repositories
	// nonForkedRepos := make([]*github.Repository, 0)
	// for _, repo := range repos {
	// 	if !*repo.Fork {
	// 		nonForkedRepos = append(nonForkedRepos, repo)
	// 	}
	// }
	// println("Total non-forked repositories:", len(nonForkedRepos))

	// // Serialize the nonForkedRepos to JSON, and save to ./repo-github.json
	// file, err := os.Create("./repo-github.json")
	// if err != nil {
	// 	fmt.Println("Error creating JSON file:", err)
	// 	return
	// }
	// defer file.Close()

	// encoder := json.NewEncoder(file)
	// encoder.SetIndent("", "  ")
	// if err := encoder.Encode(nonForkedRepos); err != nil {
	// 	fmt.Println("Error encoding JSON:", err)
	// 	return
	// }
	// fmt.Println("Non-forked repositories saved to ./repo-github.json")

	// Deserialize the JSON file
	// repos := make([]*github.Repository, 0)
	// file, err := os.Open("./repo-github.json")
	// if err != nil {
	// 	fmt.Println("Error opening JSON file:", err)
	// 	return
	// }
	// defer file.Close()
	// decoder := json.NewDecoder(file)
	// if err := decoder.Decode(&repos); err != nil {
	// 	fmt.Println("Error decoding JSON:", err)
	// 	return
	// }
	// // Print the non-forked repositories
	// for _, repo := range repos {
	// 	fmt.Printf("Name: %s, Forked: %t\n", *repo.Name, *repo.Fork)
	// }

	// // Try to perform migration
	// for _, repo := range repos {
	// 	_, resp, err := giteaClient.GetRepo("situ2001-github-mirror", *repo.Name)

	// 	isGiteaRepoNotFound := resp.StatusCode == 404

	// 	if err != nil && !isGiteaRepoNotFound {
	// 		fmt.Println("Error getting Gitea repository:", err)
	// 		continue
	// 	}

	// 	giteaClient.MirrorGithubRepository(repo, "situ2001-github-mirror", githubToken)
	// }
}
