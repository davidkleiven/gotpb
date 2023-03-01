package gotpb

import (
	"database/sql"
	"gotpb/gotpb/t_utils"
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

	query := "SELECT name FROM sqlite_master WHERE type='table'"
	rows, err := db.Query(query)

	if err != nil {
		t.Errorf("%v", err)
	}

	names := []string{}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	expect := []string{"songs", "notifications"}
	if len(names) != len(expect) {
		t.Errorf("Expectd %v got %v\n", expect, names)
		return
	}

	for i, name := range expect {
		if name != names[i] {
			t.Errorf("Expected %v got %v\n", expect, names)
			return
		}
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

func TestInsertFetchNotifications(t *testing.T) {
	db, err := sql.Open("sqlite3", t_utils.SqliteInMemResource(t.Name()))
	if err != nil {
		t.Errorf("%v\n", err)
	}
	initDb(db)

	for _, group := range []string{"group1", "group2"} {
		insertSongListNotification(db, group)
	}

	notification1 := getLatestSongListNotification(db, "group1")
	notification2 := getLatestSongListNotification(db, "group2")

	if notification2.Before(notification1) {
		t.Errorf("Notification 2 was inserted after notification end. t1: %v, t2: %v\n", notification1, notification2)
	}
}

func TestNewSongs(t *testing.T) {
	songs := []Song{
		{
			Code:  1,
			Title: "One",
			Ext:   "pdf",
		},
		{
			Code:  2,
			Title: "Two",
			Ext:   "pdf",
		},
	}

	old := []Song{
		{
			Code:  1,
			Title: "One",
			Ext:   "pdf",
		},
	}

	new := newSongs(songs, old)
	expect := []Song{
		{
			Code:  2,
			Title: "Two",
			Ext:   "pdf",
		},
	}

	if len(expect) != len(new) {
		t.Errorf("Expectd 1 new song. God %d", len(new))
	}

	for i := range expect {
		if expect[i].Title != new[i].Title {
			t.Errorf("Expected %s got %s", expect[i].Title, new[i].Title)
		}
	}
}

func TestGetLatestSongNotificationNoDBContent(t *testing.T) {
	db, err := sql.Open("sqlite3", t_utils.SqliteInMemResource(t.Name()))
	if err != nil {
		t.Errorf("%v\n", err)
	}
	initDb(db)
	latest := getLatestSongListNotification(db, "group")

	if latest != defaultTime() {
		t.Errorf("Expected %v got %v", defaultTime(), latest)
	}
}

func TestTimeFormatting(t *testing.T) {
	now := time.Now().UTC().Round(time.Second)
	now_formatted := now.Format(TIME_LAYOUT)
	now_unformatted, err := time.Parse(TIME_LAYOUT, now_formatted)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if now != now_unformatted {
		t.Errorf("Expected %v got %v\n", now, now_unformatted)
	}
}
