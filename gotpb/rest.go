package gotpb

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func StatusHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Success"
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
	w.WriteHeader(http.StatusOK)
}

type DbConnector interface {
	DbConnection() *sql.DB
}

type GotbbHandlers struct {
	connector DbConnector
}

func (gh *GotbbHandlers) LatestNotificationHandlerFunc(w http.ResponseWriter, r *http.Request) {
	db := gh.connector.DbConnection()
	defer db.Close()

	numArg := r.URL.Query().Get("num")
	num, err := strconv.Atoi(numArg)
	if err != nil {
		if numArg != "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		num = 10
		log.Printf("Error: %v\nUsing default value %d", err, num)
	}

	notifications, err := getLatestNotifications(db, num)
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonNotifications, err := json.Marshal(notifications)
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonNotifications)
}

func InitRestService(conf Config) {
	handlers := GotbbHandlers{connector: conf}
	statusHandler := http.HandlerFunc(StatusHandlerFunc)
	latestNotificationHandler := http.HandlerFunc(handlers.LatestNotificationHandlerFunc)
	http.Handle("/status", statusHandler)
	http.Handle("/notifications", latestNotificationHandler)
}
