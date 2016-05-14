package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

func initDatabase() (*gorp.DbMap, error) {
	// Open the database.
	db, err := sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		return nil, err
	}
	// Create a gorp mapping.
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	// Add tables to the mapping.
	dbmap.AddTableWithName(repoRecord{}, "repositories").SetKeys(false, "RepoID")
	dbmap.AddTableWithName(repoStringsRecord{}, "repostrings").SetKeys(true, "RepostringsID")
	dbmap.AddTableWithName(authorisedUserRecord{}, "authedusers").SetKeys(false, "UserID")
	// Create the tables if they don't exist already.
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}
	return dbmap, nil
}
