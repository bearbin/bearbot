package main

import (
	"net/http"

	"encoding/json"
	"github.com/google/go-github/github"
)

type pullRequestWebhookData struct {
	// What action was performed on the pull request. None of the actions have a
	// special relevance for us.
	Action string `json:"action"`
	// The pull request number.
	Number int `json:"number"`
	PullRequest struct {
		// The state of the pull request (open/closed etc.)
		State string `json:"state"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`
}

func handlePullRequestUpdate(w http.ResponseWriter, r *http.Request, whuc *webhookUpdateContext) error {
	// Unmarshal the body JSON into a usable form.
	pr := &pullRequestWebhookData{}
	err := json.Unmarshal(whuc.Body, pr)
	if err != nil {
		return err
	}

	// There's no point checking anything about an open pull request.
	if pr.PullRequest.State != "open" {
		return nil
	}

	return nil
}

func handlePullRequestSynchronize(pruc *pullRequestUpdateContext) error {
	// Get the newest commit ID for the pull request.
	latestCommit, err := pruc.PullData.Get("head").Get("sha").String()
	if err != nil {
		return err
	}

	// Update the pull request record to reflect the new head commit.
	pruc.Pull.Head = latestCommit
	_, err = dbmap.Update(pruc.Pull)
	if err != nil {
		return err
	}

	// TODO: Is this neccessary? Should the old signoffs just be kept around.
	// Delete signoffs for the older commit.
	signoffs, err := dbmap.Select(
		signoffRecord{},
		"SELECT * FROM signoffs WHERE CommitHash=?",
		latestCommit,
	)
	if err != nil {
		return err
	}
	n, err := dbmap.Delete(signoffs...)
	if err != nil {
		return err
	}
	if n > 0 {
		// Post a comment informing that signoffs have been reset, since a
		// new commit has been pushed.
		// Get the string from the database.
		text, err := getRepoStringByName(pruc.Repo.RepoID, "signoffsremoved")
		if err != nil {
			return err
		}
		_, _, err = ghClient.Issues.CreateComment(
			pruc.Repo.Owner,
			pruc.Repo.Name,
			pruc.Pull.Number,
			&github.IssueComment{Body: &(text.StringText)},
		)
		if err != nil {
			return err
		}
	}

	err = updateInfoComment(pruc.Pull)
	if err != nil {
		return err
	}
	return nil
}
