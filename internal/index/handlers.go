package index

import "net/http"

func HandleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Welcome"))
}
