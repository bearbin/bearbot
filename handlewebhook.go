package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"
)

func handleWebhook(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the repository this request is for.
	repository := &repoRecord{}
	err := dbmap.SelectOne(
		repoRecord{},
		"SELECT * FROM repositories WHERE Owner=? AND RepoName=?",
		c.URLParams["owner"],
		c.URLParams["reponame"],
	)
	if err != nil {
		log.Println("handleWebhook: ", err.Error())
		http.Error(w, "database select failed", http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("handleWebhook: ", err.Error())
		http.Error(w, "body read failed", http.StatusInternalServerError)
		return
	}

	ghSignature := r.Header.Get("X-Hub-Signature")

	signatureMatch := verifyGHSignature(body, ghSignature, repository.WebhookSecret)
	if !signatureMatch {
		http.Error(w, "403 Forbidden - HMAC verification failed", http.StatusForbidden)
	}
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
