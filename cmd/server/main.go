package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/bansal-arsh/personal-finance-tracker/internal/email"
	"github.com/bansal-arsh/personal-finance-tracker/internal/index"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	switch env {
	case "DEV":
		slog.Info("DEV environment detected. Loding environemnt variables...")

		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}
	case "PROD":
	default:
		log.Fatal("ENV variable not set")
	}

	slog.Info("Setting up gmail dialer...")
	gmailDialer, err := email.NewGmailDialer("arshbansal1111@gmail.com", os.Getenv("GMAIL_APP_PASSWORD"))
	if err != nil {
		log.Fatalf("Error creating gmail dialer: %s", err)
	}

	slog.Info("Setting up HTML server...")
	mux := http.NewServeMux()
	srv := &http.Server{Addr: "0.0.0.0:80", Handler: mux}

	mux.HandleFunc("/{$}", index.HandleIndex)
	mux.Handle("/send", index.HandleEmail(gmailDialer))

	slog.Info("Starting server...")
	log.Fatal(srv.ListenAndServe())
}
