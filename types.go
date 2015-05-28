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

type pullRequestRecord struct {
	PullID int
	RepoID int
	Number int
	Head   string
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

type signoffRecord struct {
	SignoffID   int
	CommitHash  int
	UserID      int
	DateCreated time.Time
}

type infoCommentContents struct {
	SignoffThreshold int
	Signoffs         []signoffRecord
}
