package gotpb

import (
	"reflect"
	"testing"
)

func TestGetConf(t *testing.T) {
	conf := GetConf("test_data/config.yml")

	expect := map[string]string{"solo": "www.example.com"}

	if !reflect.DeepEqual(conf.Groups, expect) {
		t.Errorf("Expected %s got %s", expect, conf.Groups)
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

func invalidConf() Config {
	return Config{
		Groups: map[string]string{"solo": "www.example.com"},
		Users:  []User{{Email: "dk@hotmail.com", Group: "secondCornet"}}}
}

func TestValidateConf(t *testing.T) {
	for i, test := range []struct {
		conf   Config
		expect bool
	}{
		{
			conf:   GetConf("test_data/config.yml"),
			expect: true,
		},
		{
			conf:   invalidConf(),
			expect: false,
		},
	} {
		res := ValidateConf(test.conf)
		if res != test.expect {
			t.Errorf("Test #%d: %v got %v", i, res, test.expect)
		}
	}
}
