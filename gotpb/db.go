package gotpb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const CREATE_STMT = `
CREATE TABLE IF NOT EXISTS songs
(code INTEGER NOT NULL, title TEXT, timestamp TEXT, active INTEGER)
`

const CREATE_NOTIFICATIONS = `
CREATE TABLE IF NOT EXISTS notifications (timestamp TEXT, type TEXT, instrument_group TEXT)
`

const INSERT_STMT = `
INSERT INTO songs VALUES (?, ?, ?, ?)
`

const INSERT_NOTIFICATION = `
INSERT INTO notifications VALUES (?, ?, ?)
`

const GET_SONGS = `
SELECT code, title FROM songs WHERE timestamp >= ?
`

const LATEST_NOTIFICATION = `
SELECT MAX(timestamp) FROM notifications WHERE type = ? AND instrument_group = ?
`

const NEW_SONG = "newSong"
const SONG_LIST = "songList"
const TIME_LAYOUT = "2006-01-02T00:00:00Z"

func defaultTime() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
}

func initDb(db *sql.DB) {
	statements := []string{CREATE_STMT, CREATE_NOTIFICATIONS}
	for _, query := range statements {
		statement, err := db.Prepare(query)
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec()
	}
}

func insertNotification(db *sql.DB, notificationType string, group string) {
	statement, err := db.Prepare(INSERT_NOTIFICATION)
	if err != nil {
		log.Printf("insertNotifiction for type %s for group %s failed with %v", notificationType, group, err)
		return
	}
	t := time.Now().UTC().Format(TIME_LAYOUT)
	statement.Exec(t, notificationType, group)
}

func insertNewSongNotification(db *sql.DB, group string) {
	insertNotification(db, NEW_SONG, group)
}

func insertSongListNotification(db *sql.DB, group string) {
	insertNotification(db, SONG_LIST, group)
}

func getLatestSongListNotification(db *sql.DB, group string) time.Time {
	rows, err := db.Query(LATEST_NOTIFICATION, SONG_LIST, group)

	if err != nil {
		log.Fatal(err)
	}

	if rows.Next() {
		var timestamp string
		rows.Scan(&timestamp)

		if t, err := time.Parse(TIME_LAYOUT, timestamp); err == nil {
			return t
		}

		log.Printf("Error while parsing timestamp %s: %v\n", timestamp, err)
	}
	return defaultTime()
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
	t := time.Now().UTC().Format(TIME_LAYOUT)
	statement, err := db.Prepare(INSERT_STMT)
	if err != nil {
		log.Printf("Error while inserting song %v", err)
		return
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

type NotificationInfo struct {
	Time  string `json:"time"`
	Type  string `json:"type"`
	Group string `json:"group"`
}

// Return num latest notifications
func getLatestNotifications(db *sql.DB, num int) ([]NotificationInfo, error) {
	query := fmt.Sprintf("SELECT * FROM notifications ORDER BY timestamp LIMIT %d", num)
	rows, err := db.Query(query)

	if err != nil {
		return []NotificationInfo{}, err
	}

	notificationInfo := []NotificationInfo{}
	for rows.Next() {
		i := NotificationInfo{}
		rows.Scan(&i.Time, &i.Type, &i.Group)
		notificationInfo = append(notificationInfo, i)
	}
	return notificationInfo, nil
}
