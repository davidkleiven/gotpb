package gotpb

import (
	"database/sql"
	"fmt"
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
const TIME_LAYOUT = time.RFC3339

func defaultTime() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
}

func initDb(db *sql.DB) {
	statements := []string{CREATE_STMT, CREATE_NOTIFICATIONS}
	tx, err := db.Begin()
	panicOnErr(err)
	defer tx.Rollback()

	for _, query := range statements {
		_, err = tx.Exec(query)
		panicOnErr(err)
	}
	err = tx.Commit()
	panicOnErr(err)
}

func insertNotification(db *sql.DB, notificationType string, group string) {
	t := time.Now().UTC().Round(time.Second).Format(TIME_LAYOUT)
	_, err := db.Exec(INSERT_NOTIFICATION, t, notificationType, group)
	panicOnErr(err)
}

func insertNewSongNotification(db *sql.DB, group string) {
	insertNotification(db, NEW_SONG, group)
}

func insertSongListNotification(db *sql.DB, group string) {
	insertNotification(db, SONG_LIST, group)
}

func getLatestSongListNotification(db *sql.DB, group string) time.Time {
	rows, err := db.Query(LATEST_NOTIFICATION, SONG_LIST, group)
	panicOnErr(err)
	defer rows.Close()

	if rows.Next() {
		var timestamp string
		rows.Scan(&timestamp)

		if t, err := time.Parse(TIME_LAYOUT, timestamp); err == nil {
			return t
		}
	}
	return defaultTime()
}

func getDB(dbName string) *sql.DB {
	if _, err := os.Stat(dbName); err != nil {
		os.Create(dbName)
	}
	db, err := sql.Open("sqlite3", dbName)
	panicOnErr(err)
	initDb(db)
	return db
}

func insertSong(db *sql.DB, song Song) {
	t := time.Now().UTC().Format(TIME_LAYOUT)
	statement, err := db.Prepare(INSERT_STMT)
	panicOnErr(err)
	defer statement.Close()
	_, err = statement.Exec(song.Code, song.Title, t, 1)
	panicOnErr(err)
}

func insertSongs(db *sql.DB, songs []Song) {
	for _, song := range songs {
		insertSong(db, song)
	}
}

func fetchNewerSongs(db *sql.DB, t time.Time) []Song {
	songs := []Song{}
	rows, err := db.Query(GET_SONGS, t.String())
	panicOnErr(err)
	defer rows.Close()

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
func getLatestNotifications(db *sql.DB, num int) []NotificationInfo {
	query := fmt.Sprintf("SELECT * FROM notifications ORDER BY timestamp LIMIT %d", num)
	rows, err := db.Query(query)
	panicOnErr(err)
	defer rows.Close()

	notificationInfo := []NotificationInfo{}
	for rows.Next() {
		i := NotificationInfo{}
		rows.Scan(&i.Time, &i.Type, &i.Group)
		notificationInfo = append(notificationInfo, i)
	}
	return notificationInfo
}
