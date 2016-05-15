package main

import (
	"time"
)

type configuration struct {
	APIToken     string `json:"api_token"`
	GHUsername   string `json:"gh_username"`
	CanonicalURL string `json:"canonical_url"`
	DatabasePath string `json:"database_path"`
}

type repoRecord struct {
	RepoID           int
	Owner            string
	Name             string
	WebhookSecret    string
	SignoffThreshold int
}

type repoStringsRecord struct {
	RepoStringsID int
	RepoID        int
	StringType    string
	StringText    string
}

type authorisedUserRecord struct {
	UserID   int
	RepoID   int
	Username string
}
