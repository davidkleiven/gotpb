package gotpb

import (
	"archive/zip"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
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
		sendNewSongNotification(newSongs, name, conf)
		sendSongListNotification(songs, name, conf)
	}
}

type Song struct {
	Code    int
	Title   string
	Ext     string
	Content *zip.File
}

func songFromFile(file *zip.File) Song {
	re := regexp.MustCompile(`([0-9]+) ([A-Za-z0-9 -ÆØÅæøå]+)\.([a-z]+)`)
	res := re.FindStringSubmatch(file.Name)
	if len(res) != 4 {
		log.Printf("Could not extract song information from %s. Num capture groups %d", file.Name, len(res))
		return Song{Code: -1, Title: file.Name, Ext: "unknown", Content: file}
	}
	code, _ := strconv.Atoi(res[1])
	return Song{
		Code:    code,
		Title:   res[2],
		Ext:     res[3],
		Content: file,
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

func sendNewSongNotification(songs []Song, group string, conf Config) {
	users := conf.UsersInGroup(group)
	if len(songs) == 0 || len(users) == 0 {
		log.Printf("Number of new songs %d. Number of users in group %d. No notifications sent.", len(songs), len(users))
		return
	}

	email := prepareEmail(conf, users)
	email.SetBody(mail.TextPlain, produceEmail(songs))
	email.SetSubject("New songs")
	sendEmail(email, conf)

	db := getDB(conf.Db)
	defer db.Close()
	insertNewSongNotification(db, group)
	log.Printf("New songs notification sent")
}

func produceEmail(songs []Song) string {
	msg := fmt.Sprintf("%d new songs available:\n\n", len(songs))
	for _, song := range songs {
		msg += fmt.Sprintf("%s\n", song.Title)
	}
	return msg
}

func sendSongListNotification(songs []Song, group string, conf Config) {
	timestamp := time.Now().UTC()
	users := conf.UsersInGroup(group)
	if timestamp.Weekday() != time.Wednesday {
		log.Printf("Today is %v. No email sent. (Sends only on %v)", timestamp.Weekday(), time.Wednesday)
		return
	}
	db := getDB(conf.Db)
	defer db.Close()
	latest := getLatestSongListNotification(db, group)

	if time.Since(latest) < time.Hour*time.Duration(48) {
		log.Printf("Less than 48 hours since last song list notification. No notification sent\n")
		return
	}

	email := prepareEmail(conf, users)
	email.SetSubject("Summary")
	email.SetBody(mail.TextPlain, produceEmail(songs))
	sendEmail(email, conf)
	insertSongListNotification(db, group)
	log.Printf("Song list notification sent")
}
