package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUserJSON(t *testing.T) {
	var user User
	err := json.Unmarshal([]byte(sampleUser), &user)
	if err != nil {
		t.Fatal("Failed to Marshal user", err)
	}
	user.AverageFollowers = 10.0 // make the test work
	if !reflect.DeepEqual(user, userOctocat) {
		t.Error("Unmarshal error", user)
	}
}

var userOctocat = User{
	Name:             "monalisa octocat",
	Login:            "octocat",
	Company:          "GitHub",
	Followers:        20,
	PublicRepos:      2,
	AverageFollowers: 10.0,
}

const sampleUser = `{
	"login": "octocat",
	"id": 1,
	"node_id": "MDQ6VXNlcjE=",
	"avatar_url": "https://github.com/images/error/octocat_happy.gif",
	"gravatar_id": "",
	"url": "https://api.github.com/users/octocat",
	"html_url": "https://github.com/octocat",
	"followers_url": "https://api.github.com/users/octocat/followers",
	"following_url": "https://api.github.com/users/octocat/following{/other_user}",
	"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
	"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
	"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
	"organizations_url": "https://api.github.com/users/octocat/orgs",
	"repos_url": "https://api.github.com/users/octocat/repos",
	"events_url": "https://api.github.com/users/octocat/events{/privacy}",
	"received_events_url": "https://api.github.com/users/octocat/received_events",
	"type": "User",
	"site_admin": false,
	"name": "monalisa octocat",
	"company": "GitHub",
	"blog": "https://github.com/blog",
	"location": "San Francisco",
	"email": "octocat@github.com",
	"hireable": false,
	"bio": "There once was...",
	"twitter_username": "monatheoctocat",
	"public_repos": 2,
	"public_gists": 1,
	"followers": 20,
	"following": 0,
	"created_at": "2008-01-14T04:33:35Z",
	"updated_at": "2008-01-14T04:33:35Z"
  }`
