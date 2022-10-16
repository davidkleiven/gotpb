package gotpb

import "testing"

func TestGetConf(t *testing.T) {
	conf := GetConf("test_data/config.yml")

	if conf.Link != "www.example.com" {
		t.Errorf("Expected 'www.example.com' got %s", conf.Link)
	}

	if len(conf.Users) != 2 {
		t.Errorf("Expected 2 users got %d", len(conf.Users))
	}

	expected := []string{"user1@gmail.com", "user2@gmail.com"}

	for i := 0; i < 2; i++ {
		if conf.Users[i].Email != expected[i] {
			t.Errorf("Test #%d: Expected %s got %s", i, conf.Users[i].Email, expected[i])
		}
	}
}
