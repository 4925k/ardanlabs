package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	name, count, err := githubInfo(ctx, "4925k")
	if err != nil {
		log.Fatalf("getting github info: %v", err)
	}

	fmt.Printf("%s has %d public repos\n", name, count)
}

// githubInfo returns the name and number of public repos of a github user
func githubInfo(ctx context.Context, username string) (string, int, error) {
	// resp := http.Get("https://api.github.com/users/" + url.PathEscape(username)) // WITHOUT CONTEXT
	url := fmt.Sprintf("https://api.github.com/users/%s", url.PathEscape(username))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	resp, err := http.DefaultClient.Do(req)
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
