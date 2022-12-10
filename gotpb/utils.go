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
		songs := filterAndInsertSongsInDb(<-file, conf)
		sendEmailNotification(songs, conf.UsersInGroup(name))
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

func filterAndInsertSongsInDb(zipFile string, conf Config) []Song {
	songs := songsFromZip(zipFile)
	t := time.Now().UTC().Add(-time.Hour * HOURS_PER_DAY * DAYS_PER_MONTH * conf.MemoryMonths)

	db := getDB(conf.Db)
	defer db.Close()

	knownSongs := fetchNewerSongs(db, t)
	newSongs := newSongs(songs, knownSongs)
	insertSongs(db, newSongs)
	return newSongs
}

func sendEmailNotification(songs []Song, users []User) {
	if len(songs) == 0 || len(users) == 0 {
		log.Printf("Number of new songs %d. Number of users in group %d. No notifications sent.", len(songs), len(users))
		return
	}

	msg := produceEmail(songs)
	log.Printf("%s", msg)
}

func produceEmail(songs []Song) string {
	msg := fmt.Sprintf("%d new songs available:\n\n", len(songs))
	for _, song := range songs {
		msg += fmt.Sprintf("%s\n", song.Title)
	}
	return msg
}
