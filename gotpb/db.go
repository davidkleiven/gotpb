package gotpb

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const CREATE_STMT = `
CREATE TABLE IF NOT EXISTS songs
(code INTEGER NOT NULL, title TEXT, timestamp TEXT, active INTEGER)
`

const INSERT_STMT = `
INSERT INTO songs VALUES (?, ?, ?, ?)
`

const GET_SONGS = `
SELECT code, title FROM songs WHERE timestamp >= ?
`

func initDb(db *sql.DB) {
	statement, err := db.Prepare(CREATE_STMT)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func getDB(dbName string) *sql.DB {
	if _, err := os.Stat(dbName); err != nil {
		os.Create(dbName)
	}
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	initDb(db)
	return db
}

func insertSong(db *sql.DB, song Song) {
	t := time.Now().UTC().String()
	statement, err := db.Prepare(INSERT_STMT)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(song.Code, song.Title, t, 1)
}

func insertSongs(db *sql.DB, songs []Song) {
	for _, song := range songs {
		insertSong(db, song)
	}
}

func fetchNewerSongs(db *sql.DB, t time.Time) []Song {
	songs := []Song{}
	rows, err := db.Query(GET_SONGS, t.String())
	if err != nil {
		log.Printf("%v", err)
		return songs
	}

	for rows.Next() {
		song := Song{}
		rows.Scan(&song.Code, &song.Title)
		songs = append(songs, song)
	}
	return songs
}

func newSongs(new []Song, old []Song) []Song {
	songs := []Song{}

	// Create map of old songs
	oldMap := map[string]bool{}
	for _, song := range old {
		oldMap[song.Title] = true
	}

	for _, song := range new {
		if _, ok := oldMap[song.Title]; !ok {
			songs = append(songs, song)
		}
	}
	return songs
}
