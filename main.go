package main

import (
	"log"
	"os"
)

func main() {
	githubEndpoint, err := os.ReadFile("env/endpoint")
	if err != nil {
		log.Fatal("Failed to read endpoint from env")
	}

	client := GitHubApiClient{endpoint: string(githubEndpoint)}
	err = runServer(&client)
	if err != nil {
		log.Fatal(err)
	}
}
