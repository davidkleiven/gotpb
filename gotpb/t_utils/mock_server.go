package t_utils

import (
	"net/http"
	"net/http/httptest"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("filecontent"))
}

func MockDownloadServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}
