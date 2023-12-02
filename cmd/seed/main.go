package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/et-codes/songvote"
)

const (
	port             = ":5050"
	dbPath           = "./db/songvote.db"
	userTestDataFile = "./testdata/users.json"
	songTestDataFile = "./testdata/songs.json"
)

func main() {
	store := songvote.NewSQLiteStore(dbPath)
	if err := store.ClearDB("DELETE_IT_ALL"); err != nil {
		log.Fatalf("Error clearing DB: %v", err)
	}

	populateWithUsers(store, userTestDataFile)
	populateWithSongs(store, songTestDataFile)

	log.Printf("Successfully cleared database and seeded with sample data.")
}

func populateWithUsers(store *songvote.SQLiteStore, path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not open file %s: %v", path, err)
	}
	users := songvote.Users{}
	if err = json.Unmarshal(file, &users); err != nil {
		log.Fatalf("Could not unmarshal JSON from %s: %v", path, err)
	}
	for _, user := range users {
		_, err := store.AddUser(user)
		if err != nil {
			log.Fatalf("Could not add user: %v", err)
		}
	}
}

func populateWithSongs(store *songvote.SQLiteStore, path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not open file %s: %v", path, err)
	}
	songs := songvote.Songs{}
	if err = json.Unmarshal(file, &songs); err != nil {
		log.Fatalf("Could not unmarshal JSON from %s: %v", path, err)
	}
	for _, song := range songs {
		_, err := store.AddSong(song)
		if err != nil {
			log.Fatalf("Could not add song: %v", err)
		}
	}
}
