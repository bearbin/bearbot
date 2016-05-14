package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/bmatsuo/go-jsontree"
	"github.com/google/go-github/github"

	"github.com/zenazn/goji/web"
)

type webhookUpdateContext struct {
	Repo *repoRecord
	Body []byte
	JSON *jsontree.JsonTree
}

// Function handleWebhook is called when an event is produced and GitHub sends
// off a webhook.
func handleWebhook(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the repository this request is for.
	repo, err := getRepoByName(c.URLParams["owner"], c.URLParams["reponame"])
	if err != nil {
		log.Println("handleWebhook: ", err.Error())
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	// Read all the body data.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("handleWebhook: ", err.Error())
		http.Error(w, "body read failed", http.StatusInternalServerError)
		return
	}
	// JSON-decode the data.
	hd := jsontree.New()
	err = hd.UnmarshalJSON(body)
	if err != nil {
		log.Println("handleWebhook: ", err.Error())
		http.Error(w, "json decode failed", http.StatusInternalServerError)
	}

	whuc := &webhookUpdateContext{
		repo,
		body,
		hd,
	}

	// Verify the GitHub signature.
	ghSignature := r.Header.Get("X-Hub-Signature")
	signatureMatch := verifyGHSignature(body, ghSignature, repo.WebhookSecret)
	if !signatureMatch {
		http.Error(w, "403 Forbidden - HMAC verification failed", http.StatusForbidden)
	}

	// Select the correct method to call depening on the
	switch r.Header.Get("X-Github-Event") {
	case "pull_request":
		err = handlePullRequestUpdate(w, r, whuc)
	case "issue_comment":
		err = handleIssueCommentUpdate(w, r, whuc)
	case "status":
		err = handleCommitStatusUpdate(w, r, whuc)
	}
	if err != nil {
		log.Println("handleWebhook: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Function updateInfoComment updates the info comment on a pull request, and
// creates one if it does not already exist.
func updateInfoComment(pr *pullRequestRecord) error {
	// First, get the pull request record.
	tmp, err := dbmap.Get(&repoRecord{}, pr.RepoID)
	repo := tmp.(*repoRecord)
	if err != nil {
		return err
	}
	// Get the strings.
	stringRecord, err := getRepoStringByName(pr.RepoID, "infocomment")
	if err != nil {
		return err
	}
	// Get the signoffs.
	signoffs, err := getSignoffsByPullRequest(pr.PullID)
	if err != nil {
		return err
	}
	// Create the template.
	icc := &infoCommentContents{
		repo.SignoffThreshold,
		signoffs,
	}
	tpl := template.New("1")
	tpl.Parse(stringRecord.StringText)
	var completedTemplate bytes.Buffer
	tpl.Execute(&completedTemplate, icc)
	ctp := completedTemplate.String()
	ilco := &github.IssueListCommentsOptions{
		Sort:      "created",
		Direction: "asc",
	}
	issueComments, _, err := ghClient.Issues.ListComments(repo.Owner, repo.Name, pr.Number, ilco)
	if err != nil {
		return err
	}

	return nil // not implemented yet
}

func handleIssueCommentUpdate(w http.ResponseWriter, r *http.Request, whuc *webhookUpdateContext) error {
	return nil // Not implemented.
}

func handleCommitStatusUpdate(w http.ResponseWriter, r *http.Request, whuc *webhookUpdateContext) error {
	return nil // Not implemented.
}

// Function getRepoByName gets the associated repoRecord for a given Owner and
// RepoName.
func getRepoByName(owner string, reponame string) (*repoRecord, error) {
	repo := &repoRecord{}
	err := dbmap.SelectOne(
		repo,
		"SELECT * FROM repositories WHERE Owner= ? AND RepoName= ?",
		owner,
		reponame,
	)
	return repo, err
}

// Function getRepoStringByName gets a repoString by its name using the specified
// repository ID.
func getRepoStringByName(repoID int, stringName string) (*repoStringsRecord, error) {
	repoString := &repoStringsRecord{}
	err := dbmap.SelectOne(
		repoString,
		"SELECT * FROM repostrings WHERE RepoID = ? AND StringType = ?",
		repoID,
		stringName,
	)
	return repoString, err
}

func getSignoffsByPullRequest(prid int) ([]signoffRecord, error) {
	return nil, nil
}

// Verifies the HMAC signature provided by GitHub for Webhooks. If the signature
// does not equal the computed signature for the provided text and secret this
// function returns false.
func verifyGHSignature(text []byte, sig string, secret string) bool {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(text)
	expectedMAC := mac.Sum(nil)
	expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(expectedSig), []byte(sig))
}
