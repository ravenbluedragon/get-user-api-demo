package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
)

type Server struct {
	client Client
}

// runServer inits and runs the http server
func runServer(client Client) error {
	server := Server{client}
	http.HandleFunc("/", server.printUsage)
	http.HandleFunc("/retrieveUsers", server.retrieveUsers)
	return http.ListenAndServe(":8080", nil)
}

const usage = `Available Endpoints:
- POST /retrieveUsers : Return User Details; Expects Array of Strings
`

// printUsage prints available endpoints
func (server *Server) printUsage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, usage)
}

// retreiveUsers expects a POST request with a JSON array of usernames
// returns a JSON array of User Details
func (server *Server) retrieveUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var usernames []string
	err := json.NewDecoder(r.Body).Decode(&usernames)
	if err != nil {
		// log.Println("Invalid Body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := server.client.get(sortUnique(usernames))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// sortUnique sorts and dedups a slice of strings and removes the empty string
func sortUnique(strings []string) []string {
	sort.StringSlice(strings).Sort()
	uniq := make([]string, 0, len(strings))
	prev := ""
	for _, s := range strings {
		if s != prev {
			prev = s
			uniq = append(uniq, s)
		}
	}
	return uniq
}
