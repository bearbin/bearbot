package main

import (
	"net/http"

	"github.com/bmatsuo/go-jsontree"
	"github.com/google/go-github/github"
)

type pullRequestUpdateContext struct {
	*webhookUpdateContext
	Pull     *pullRequestRecord
	PullData *jsontree.JsonTree
}

func handlePullRequestUpdate(w http.ResponseWriter, r *http.Request, whuc *webhookUpdateContext) error {
	// Get useful information about the pull request.
	prdata := whuc.JSON.Get("pull_request")
	// What type of event was this update sent for?
	action, err := whuc.JSON.Get("action").String()
	if err != nil {
		return err
	}
	prid, err := prdata.Get("id").Number()
	if err != nil {
		return err
	}
	tmp, err := dbmap.Get(pullRequestRecord{}, int(prid))
	if err != nil {
		return err
	}
	pruc := &pullRequestUpdateContext{
		webhookUpdateContext: whuc,
		Pull:                 tmp.(*pullRequestRecord),
		PullData:             prdata,
	}
	switch action {
	case "synchronize":
		handlePullRequestSynchronize(pruc)
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
