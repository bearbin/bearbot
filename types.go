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
	RepoName         string
	WebhookSecret    string
	SignoffThreshold int
}

type pullRequestRecord struct {
	PullID int
	RepoID int
	Number int
}

type repoStringsRecord struct {
	RepostringsID int
	RepoID        int
	StringType    int
	StringText    string
}

type authorisedUserRecord struct {
	UserID   int
	RepoID   int
	Username string
}

type signoffRecord struct {
	SignoffID   int
	CommitID    int
	UserID      int
	DateCreated time.Time
}
