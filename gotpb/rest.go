package gotpb

import (
	"encoding/json"
	"net/http"
)

func StatusHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Success"
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func InitRestService() {
	statusHandler := http.HandlerFunc(StatusHandlerFunc)
	http.Handle("/status", statusHandler)
}
