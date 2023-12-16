package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/et-codes/songvote"
)

const (
	port             = ":5050"
	dbPath           = "./db/songvote.db"
	userTestDataFile = "./testdata/users.json"
	songTestDataFile = "./testdata/songs.json"
	voteTestDataFile = "./testdata/votes.json"
)

func main() {
	store := songvote.NewSQLiteStore(dbPath)
	if err := store.ClearDB("DELETE_IT_ALL"); err != nil {
		slog.Error("Error clearing DB: %v", err)
		os.Exit(1)
	}

	populateWithUsers(store, userTestDataFile)
	populateWithSongs(store, songTestDataFile)
	populateWithVotes(store, voteTestDataFile)

	slog.Info("Successfully cleared database and seeded with sample data.")
}

func populateWithUsers(store *songvote.SQLiteStore, path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Could not open file %s: %v", path, err)
		os.Exit(1)
	}
	users := songvote.Users{}
	if err = json.Unmarshal(file, &users); err != nil {
		slog.Error("Could not unmarshal JSON from %s: %v", path, err)
		os.Exit(1)
	}
	for _, user := range users {
		_, err := store.AddUser(user)
		if err != nil {
			slog.Error("Could not add user: %v", err)
			os.Exit(1)
		}
	}
}

func populateWithSongs(store *songvote.SQLiteStore, path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Could not open file %s: %v", path, err)
		os.Exit(1)
	}
	songs := songvote.Songs{}
	if err = json.Unmarshal(file, &songs); err != nil {
		slog.Error("Could not unmarshal JSON from %s: %v", path, err)
		os.Exit(1)
	}
	for _, song := range songs {
		_, err := store.AddSong(song)
		if err != nil {
			slog.Error("Could not add song: %v", err)
			os.Exit(1)
		}
	}
}

func populateWithVotes(store *songvote.SQLiteStore, path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Could not open file %s: %v", path, err)
		os.Exit(1)
	}
	votes := songvote.Votes{}
	if err = json.Unmarshal(file, &votes); err != nil {
		slog.Error("Could not unmarshal JSON from %s: %v", path, err)
		os.Exit(1)
	}
	for _, vote := range votes {
		err := store.AddVote(vote)
		if err != nil {
			slog.Error("Could not add song: %v", err)
			os.Exit(1)
		}
	}
}
