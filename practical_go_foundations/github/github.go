package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	name, count, err := githubInfo("4925k")
	if err != nil {
		log.Fatalf("getting github info: %v", err)
	}

	fmt.Printf("%s has %d public repos\n", name, count)
}

// githubInfo returns the name and number of public repos of a github user
func githubInfo(username string) (string, int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s", url.PathEscape(username)))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("getting github info: %s", resp.Status)
	}

	var r GithubResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&r); err != nil {
		return "", 0, err
	}

	return r.Name, r.PublicRepos, nil
}

type GithubResponse struct {
	Name        string `json:"name"`
	PublicRepos int    `json:"public_repos"`
}
