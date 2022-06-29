package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name             string  `json:"name"`
	Login            string  `json:"login"`
	Company          string  `json:"company"`
	Followers        int     `json:"followers"`
	PublicRepos      int     `json:"public_repos"`
	AverageFollowers float32 `json:"average_followers"`
}

// A Client retrieves users from a list of usernames
type Client interface {
	get(usernames []string) ([]User, error)
}

// GitHubApi Client retreives users from the GitHub API
type GitHubApiClient struct {
	// baseUrl string
	endpoint string
}

// gets Users from a slice of usernames
func (c *GitHubApiClient) get(usernames []string) ([]User, error) {
	users := make([]User, 0, len(usernames))
	for _, name := range usernames {
		user, _ := c.getUser(name) // ignore errors (from spec)
		users = append(users, user)
	}
	return users, nil
}

var errNotOK error = fmt.Errorf("status not ok")

// gets a user from a username
func (c *GitHubApiClient) getUser(username string) (User, error) {
	var user User
	res, err := http.Get(c.endpoint + username)
	if err != nil {
		return user, err
	}
	if res.StatusCode != http.StatusOK {
		return user, errNotOK
	}
	err = json.NewDecoder(res.Body).Decode(&user)
	if user.PublicRepos != 0 {
		user.AverageFollowers = float32(user.Followers) / float32(user.PublicRepos)
	}
	return user, err
}
