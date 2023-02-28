package gotpb

import (
	"archive/zip"
	"database/sql"
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
	db, err := conf.DbConnection()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	defer db.Close()

	for _, url := range conf.Groups {
		go fetch(url, file)
	}

	for name := range conf.Groups {
		songs := songsFromZip(<-file)
		newSongs := filterAndInsertSongsInDb(songs, conf, db)

		email := mail.NewMSG()
		res := sendNewSongNotification(newSongs, name, conf, email, db)
		log.Printf("%s\n", res.String())
		email = mail.NewMSG()
		res = sendSongListNotification(songs, name, conf, email, db)
		log.Printf("%s\n", res.String())
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

func filterAndInsertSongsInDb(songs []Song, conf Config, db *sql.DB) []Song {
	t := time.Now().UTC().Add(-time.Hour * HOURS_PER_DAY * DAYS_PER_MONTH * conf.MemoryMonths)

	knownSongs := fetchNewerSongs(db, t)
	newSongs := newSongs(songs, knownSongs)
	insertSongs(db, newSongs)
	return newSongs
}

func sendNewSongNotification(songs []Song, group string, conf Config, email Email, db *sql.DB) ActionResult {
	users := conf.UsersInGroup(group)
	result := ActionResult{header: "sendNewSongNotification"}
	if len(songs) == 0 || len(users) == 0 {
		result.message = fmt.Sprintf("Number of new songs %d. Number of users in group %d. No notifications sent.", len(songs), len(users))
		return result
	}

	prepareEmail(email, users)
	email.SetBody(mail.TextPlain, produceEmail(songs))
	email.SetSubject("New songs")
	sendEmail(email, conf)

	insertNewSongNotification(db, group)
	result.message = "New songs notification sent"
	return result
}

func produceEmail(songs []Song) string {
	msg := fmt.Sprintf("%d new songs available:\n\n", len(songs))
	for _, song := range songs {
		msg += fmt.Sprintf("%s\n", song.Title)
	}
	return msg
}

type ActionResult struct {
	header  string
	message string
	err     error
}

func (a *ActionResult) String() string {
	str := fmt.Sprintf("%s: %s", a.header, a.message)
	if a.err != nil {
		str += fmt.Sprintf("\n%v", a.err)
	}
	return str
}

func sendSongListNotification(songs []Song, group string, conf Config, email Email, db *sql.DB) ActionResult {
	timestamp := time.Now().UTC()
	users := conf.UsersInGroup(group)
	result := ActionResult{header: "sendSongListNotification"}

	// conf.AlwaysSend is primarly used to disable the weekday check in unit tests
	if !conf.AlwaysSend {
		if timestamp.Weekday() != time.Wednesday {
			result.message = fmt.Sprintf("Today is %v. No email sent. (Sends only on %v)", timestamp.Weekday(), time.Wednesday)
			return result
		}
	}
	latest, err := getLatestSongListNotification(db, group)
	if err != nil {
		result.err = err
		return result
	}

	if time.Since(latest) < time.Hour*time.Duration(48) {
		result.message = "Less than 48 hours since last song list notification. No notification sent"
		return result
	}

	prepareEmail(email, users)
	email.SetSubject("Summary")
	email.SetBody(mail.TextPlain, produceEmail(songs))
	sendEmail(email, conf)
	if err = insertSongListNotification(db, group); err != nil {
		result.err = err
		return result
	}
	result.message = "Song list notification sent"
	return result
}
