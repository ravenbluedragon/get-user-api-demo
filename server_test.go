package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestPrintUsage(t *testing.T) {
	server := Server{&MockClient{}}

	ts := httptest.NewServer(http.HandlerFunc(server.printUsage))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal("Failed to get", err)
	}

	out, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Failed to read body", err)
	}
	if string(out) != usage {
		t.Errorf("Expected %s, Received %s", string(out), usage)
	}
}

func TestRetrieveUsersGet(t *testing.T) {
	server := Server{&MockClient{}}

	ts := httptest.NewServer(http.HandlerFunc(server.retrieveUsers))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal("Unexpected Error")
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal("Expected fail code")
	}

	out, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Failed to read body", err)
	}
	if len(out) != 0 {
		t.Error("Expected Empty Response")
	}
}

func TestRetrieveUsersInvalid(t *testing.T) {
	server := Server{&MockClient{}}

	ts := httptest.NewServer(http.HandlerFunc(server.retrieveUsers))
	defer ts.Close()

	res, err := http.Post(ts.URL, "text/json", strings.NewReader("{"))
	if err != nil {
		t.Fatal("Unexpected Error")
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatal("Expected fail code", res.StatusCode)
	}
}

func TestRetrieveUsersClientError(t *testing.T) {
	c := ClientWithError(http.StatusInternalServerError)
	server := Server{&c}

	ts := httptest.NewServer(http.HandlerFunc(server.retrieveUsers))
	defer ts.Close()

	res, err := http.Post(ts.URL, "text/json", strings.NewReader("[]"))
	if err != nil {
		t.Fatal("Unexpected Error")
	}

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatal("Expected fail code", res.StatusCode)
	}
}

func TestRetrieveUsers(t *testing.T) {
	server := Server{&MockClient{}}

	ts := httptest.NewServer(http.HandlerFunc(server.retrieveUsers))
	defer ts.Close()

	type test struct {
		input  []string
		output []User
	}

	table := []test{
		{input: nil, output: []User{}},
		{input: []string{"a"}, output: []User{userA}},
		{input: []string{"b", "b", "a"}, output: []User{userA, userB}},
		{input: []string{"b", "c", "a"}, output: []User{userA, userB, userC}},
		{input: []string{"x", "a", "y", "c", "z"}, output: []User{userA, userC}},
	}

	for _, tc := range table {
		body, err := json.Marshal(tc.input)
		if err != nil {
			t.Fatal("Failed to create body", err)
		}

		res, err := http.Post(ts.URL, "text/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal("Unexpected Error", err)
		}

		var data []User
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			t.Error("Failed to read body", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected Status OK, Received %d", res.StatusCode)
		}
		if !reflect.DeepEqual(data, tc.output) {
			t.Errorf("\nExpected %v\nReceived %v", tc.output, data)
		}
	}
}

func TestSortUniq(t *testing.T) {
	type test struct {
		input  []string
		output []string
	}
	table := []test{
		{input: nil, output: []string{}},
		{input: []string{"a"}, output: []string{"a"}},
		{input: []string{"a", "b", "c", "d"}, output: []string{"a", "b", "c", "d"}},
		{input: []string{"a", "d", "a", "a"}, output: []string{"a", "d"}},
		{input: []string{"e", "e", "d", "a"}, output: []string{"a", "d", "e"}},
		{input: []string{"e", "", "", "a"}, output: []string{"a", "e"}},
	}

	for _, tc := range table {
		out := sortUnique(tc.input)
		if !reflect.DeepEqual(out, tc.output) {
			t.Errorf("\nExpected %v\nReceived %v", tc.output, out)
		}
	}
}

type MockClient struct{}

var userA = User{Name: "a"}
var userB = User{Name: "b"}
var userC = User{Name: "c"}

var userMap = map[string]User{
	"a": userA,
	"b": userB,
	"c": userC,
}

func (*MockClient) get(usernames []string) ([]User, error) {
	users := make([]User, 0, len(usernames))
	for _, name := range usernames {
		if user, ok := userMap[name]; ok {
			users = append(users, user)
		}
	}
	return users, nil
}

type ClientWithError int

func (c *ClientWithError) get(usernames []string) ([]User, error) {
	return nil, fmt.Errorf("STATUS %d", c)
}
