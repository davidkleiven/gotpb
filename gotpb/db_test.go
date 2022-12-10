package gotpb

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestInit(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("%v", err)
	}

	initDb(db)

	query := "SELECT * FROM sqlite_master WHERE type='table'"
	rows, err := db.Query(query)

	if err != nil {
		t.Errorf("%v", err)
	}
	count := 0
	for rows.Next() {
		count += 1
	}

	if count != 2 {
		t.Errorf("Got %d tables expected 2", count)
	}
}

func TestInsert(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("%v", err)
	}

	initDb(db)
	songs := []Song{
		{
			Code:  1,
			Title: "My song",
			Ext:   "pdf",
		},
		{
			Code:  2,
			Title: "My song",
			Ext:   "pdf",
		},
	}

	insertSongs(db, songs)

	timestamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	songsFetched := fetchNewerSongs(db, timestamp)

	if len(songsFetched) != len(songs) {
		t.Errorf("Expected 2 songs got %d", len(songsFetched))
		return
	}

	for i := 0; i < 2; i++ {
		s1, s2 := songs[i], songsFetched[i]

		if s1.Code != s2.Code || s1.Title != s2.Title {
			t.Errorf("Expected %v got %v", s1, s2)
		}
	}
}
