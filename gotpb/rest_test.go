package gotpb

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	handler := http.HandlerFunc(StatusHandlerFunc)
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d got %d\n", recorder.Code, http.StatusOK)
	}

	expected := `{"message":"Success"}`

	if recorder.Body.String() != expected {
		t.Errorf("Expected body %s got %s\n", expected, recorder.Body.String())
	}

}

type InMemoryConnector struct{}

func (imc *InMemoryConnector) DbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	initDb(db)
	insertNotification(db, "update", "solo")
	insertNotification(db, "summary", "solo")
	return db
}

func request(url string) *http.Request {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	return r
}

func TestLastNotificationHandler(t *testing.T) {
	connector := InMemoryConnector{}
	handlers := GotbbHandlers{connector: &connector}
	for i, test := range []struct {
		code       int
		length     int
		bodySubstr string
		req        *http.Request
	}{
		{
			code:   http.StatusOK,
			length: 2,
			req:    request("/notifications"),
		},
		{
			code:   http.StatusOK,
			length: 1,
			req:    request("/notifications?num=1"),
		},
		{
			code:   http.StatusUnprocessableEntity,
			length: 2,
			req:    request("/notifications?num=one"), // Can not convert to int
		},
	} {

		handler := http.HandlerFunc(handlers.LatestNotificationHandlerFunc)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, test.req)

		if recorder.Code != test.code {
			t.Errorf("Test #%d: Expected %d got %d\n", i, test.code, recorder.Code)
		}

		if recorder.Code != http.StatusOK {
			continue
		}

		// Remaining test only applies when result is success
		var notifications []NotificationInfo
		err := json.Unmarshal(recorder.Body.Bytes(), &notifications)

		if err != nil {
			t.Errorf("Test #%d: Error with during unmarshal:\n%v\n", i, err)
		}

		if len(notifications) != test.length {
			t.Errorf("Test #%d: Expected %d got %d\n", i, test.length, len(notifications))
		}

	}
}
