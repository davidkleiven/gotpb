package gotpb

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

const HOURS_PER_DAY = 24
const DAYS_PER_MONTH = 30

func panicOnErr(e error) {
	if e != nil {
		panic(e)
	}
}

func RunSingleCheck(conf Config) {
	file := make(chan string, len(conf.Groups))

	for _, url := range conf.Groups {
		go fetch(url, file)
	}

	for name := range conf.Groups {
		songs := songsFromZip(<-file)
		newSongs := filterAndInsertSongsInDb(songs, conf)
		users := conf.UsersInGroup(name)
		sendNewSongNotification(newSongs, users, conf)
		sendSongListNotification(songs, users, conf)
	}
}

type Song struct {
	Code  int
	Title string
	Ext   string
}

func songFromFilename(fname string) Song {
	re := regexp.MustCompile(`([0-9]+) ([A-Za-z0-9 -ÆØÅæøå]+)\.([a-z]+)`)
	res := re.FindStringSubmatch(fname)
	if len(res) != 4 {
		log.Printf("Could not extract song information from %s. Num capture groups %d", fname, len(res))
		return Song{Code: -1, Title: fname, Ext: "unknown"}
	}
	code, _ := strconv.Atoi(res[1])
	return Song{
		Code:  code,
		Title: res[2],
		Ext:   res[3],
	}
}

func filterAndInsertSongsInDb(songs []Song, conf Config) []Song {
	t := time.Now().UTC().Add(-time.Hour * HOURS_PER_DAY * DAYS_PER_MONTH * conf.MemoryMonths)

	db := getDB(conf.Db)
	defer db.Close()

	knownSongs := fetchNewerSongs(db, t)
	newSongs := newSongs(songs, knownSongs)
	insertSongs(db, newSongs)
	return newSongs
}

func sendNewSongNotification(songs []Song, users []User, conf Config) {
	if len(songs) == 0 || len(users) == 0 {
		log.Printf("Number of new songs %d. Number of users in group %d. No notifications sent.", len(songs), len(users))
		return
	}

	msg := produceEmail(songs)
	log.Printf("%s", msg)

	db := getDB(conf.Db)
	defer db.Close()
	insertNewSongNotification(db)
}

func produceEmail(songs []Song) string {
	msg := fmt.Sprintf("%d new songs available:\n\n", len(songs))
	for _, song := range songs {
		msg += fmt.Sprintf("%s\n", song.Title)
	}
	return msg
}

func sendSongListNotification(songs []Song, users []User, conf Config) {
	timestamp := time.Now().UTC()

	if timestamp.Weekday() != time.Friday {
		return
	}
	db := getDB(conf.Db)
	defer db.Close()
	latest := getLatestSongListNotification(db)

	if time.Since(latest) < time.Hour*time.Duration(48) {
		return
	}

	msg := produceEmail(songs)
	log.Panicf("%s", msg)
	insertSongListNotification(db)
}
