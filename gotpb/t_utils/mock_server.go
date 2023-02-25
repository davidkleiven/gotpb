package t_utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("filecontent"))
}

func MockDownloadServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func SqliteInMemResource(name string) string {
	return fmt.Sprintf("file:%s?mode=memory&cache=shared", name)
}
