package gotpb

import (
	"database/sql"
	"gotpb/gotpb/t_utils"
	"testing"

	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailBody struct {
	body        string
	contenttype mail.ContentType
}

// EmailMock is a struct that implements the Email interface
type EmailMock struct {
	from    string
	to      string
	subject string
	body    EmailBody
}

func (em *EmailMock) AddTo(to ...string) *mail.Email {
	em.to = to[0]
	return nil
}

func (em *EmailMock) Send(client *mail.SMTPClient) error {
	return nil
}

func (em *EmailMock) SetBody(t mail.ContentType, body string) *mail.Email {
	em.body = EmailBody{body: body, contenttype: t}
	return nil
}

func (em *EmailMock) SetFrom(f string) *mail.Email {
	em.from = f
	return nil
}

func (em *EmailMock) SetSubject(s string) *mail.Email {
	em.subject = s
	return nil
}

func TestSendSongListNotification(t *testing.T) {
	conf := GetConf("test_data/config.yml")
	conf.AlwaysSend = true
	db, err := sql.Open("sqlite3", t_utils.SqliteInMemResource(t.Name()))

	if err != nil {
		t.Errorf("%v\n", err)
	}
	defer db.Close()
	if err = initDb(db); err != nil {
		t.Errorf("%v\n", err)
	}

	songs := []Song{
		{
			Code:    10,
			Title:   "My song",
			Ext:     "pdf",
			Content: nil,
		},
	}

	group := "solo"
	email := EmailMock{}

	for i := 0; i < 2; i++ {
		sendSongListNotification(songs, group, conf, &email, db)

		// Now one notification should be placed in the database
		notifications, err := getLatestNotifications(db, 1000)

		if err != nil {
			t.Errorf("Attempt %d: %v\n", i, err)
		}

		if len(notifications) != 1 {
			t.Errorf("Attempt %d: Expected 1 notification got %d\n", i, len(notifications))
		}
	}
}
