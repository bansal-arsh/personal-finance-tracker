package main

import (
	"log"
	"net/http"

	"github.com/bansal-arsh/personal-finance-tracker/internal/index"
)

func main() {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: "0.0.0.0:80", Handler: mux}

	mux.HandleFunc("/{$}", index.HandleIndex)

	log.Fatal(srv.ListenAndServe())
}
