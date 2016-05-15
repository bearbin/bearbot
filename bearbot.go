package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/google/go-github/github"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web/middleware"
	"golang.org/x/oauth2"
)

var (
	config   = configuration{}
	ghClient *github.Client
	dbmap    *gorp.DbMap
)

// Load the API key from the environment variable.
func init() {
	// Seed the RNG. Only needs doing once at startup.
	rand.Seed(time.Now().UTC().UnixNano())

	// Open config file.
	configFile, err := os.Open("conf.json")
	if err != nil {
		log.Fatal("failed to open config: ", err)
	}
	defer configFile.Close()

	// Decode config.
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("failed to decode config: ", err)
	}

	// Check that all required fields in configuration are filled.
	if config.APIToken == "" {
		log.Fatal("api_token not set in configuration file")
	}
	if config.GHUsername == "" {
		log.Fatal("gh_username not set in configuration file")
	}
	if config.CanonicalURL == "" {
		log.Fatal("canonical_url not set in configuration file")
	}
	if config.DatabasePath == "" {
		log.Fatal("database_path not set in configuration file")
	}

	// Create the database object.
	dbmap, err = initDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Create the GitHub client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.APIToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	ghClient = github.NewClient(tc)
}

func main() {
	// Default Middleware
	goji.Use(middleware.EnvInit)
	// The prevents a panic crash, instead logging an error.
	goji.Use(middleware.Recoverer)
	// This GZIPs the output so less bandwidth is used.
	goji.Use(gzipHandler)

	goji.Post("/hooks/:owner/:reponame", handleWebhook)

	goji.Serve()
}
