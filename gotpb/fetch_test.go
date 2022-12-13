package gotpb

import "testing"

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
