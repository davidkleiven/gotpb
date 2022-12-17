package gotpb

import (
	"gotpb/gotpb/t_utils"
	"strings"
	"testing"
)

func TestSongFromFilename(t *testing.T) {
	fname := "1023 mysong.pdf"
	song := songFromFilename(fname)
	if song.Code != 1023 {
		t.Errorf("Expected 1023 got %d", song.Code)
	}

	if song.Title != "mysong" {
		t.Errorf("Expected mysong got %s", song.Title)
	}

	if song.Ext != "pdf" {
		t.Errorf("Expected pdf got %s", song.Ext)
	}
}

func TestFetch(t *testing.T) {
	server := t_utils.MockDownloadServer()
	defer server.Close()

	c := make(chan string, 1)
	fetch(server.URL, c)
	res := <-c
	h := urlHash(server.URL)
	if !strings.Contains(res, h) {
		t.Errorf("Expeted %s to part of filename. Got %s", h, res)
	}
}
