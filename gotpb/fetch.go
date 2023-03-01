package gotpb

import (
	"archive/zip"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const CACHE_FOLDER = ".gotbbcache"

func cacheFile(url string) string {
	newpath := filepath.Join(".", CACHE_FOLDER)
	err := os.MkdirAll(newpath, os.ModePerm)
	panicOnErr(err)

	return CACHE_FOLDER + "/" + urlHash(url) + ".zip"
}

func urlHash(url string) string {
	h := fnv.New32a()
	h.Write([]byte(url))
	return fmt.Sprintf("%x", h.Sum32())
}

func fetch(url string, c chan string) error {
	log.Print("Downloading zip archive")
	resp, err := http.Get(url)
	if err != nil {
		c <- ""
		return err
	}

	file := cacheFile(url)
	out, err := os.Create(file)
	if err != nil {
		c <- ""
		return err
	}
	defer resp.Body.Close()
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		c <- ""
		return err
	}
	log.Print("Finished downloading zip-archive")
	c <- file
	return nil
}

func songsFromZip(fname string) []Song {
	archive, err := zip.OpenReader(fname)

	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()
	songs := []Song{}
	for _, file := range archive.File {
		newSong := songFromFile(file)
		if newSong.Code > 0 {
			songs = append(songs, newSong)
		}
	}
	return songs
}
