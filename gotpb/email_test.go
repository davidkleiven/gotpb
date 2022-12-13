package gotpb

import "testing"

func TestPrepEmail(t *testing.T) {
	conf := GetConf("test_data/config.yml")
	email := prepareEmail(conf, conf.Users)

	from := "apps.brottem@gmail.com"
	if email.GetFrom() != from {
		t.Errorf("Expected %s got %s", from, email.GetFrom())
	}

	to := email.GetRecipients()
	if len(to) != len(conf.Users) {
		t.Errorf("Expected %d recipents got %d", len(conf.Users), len(to))
	}

	for i := range to {
		if to[i] != conf.Users[i].Email {
			t.Errorf("Recipent #%d: Got %s expected %s", i, to[i], conf.Users[i].Email)
		}
	}

}
