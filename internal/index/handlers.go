package index

import (
	"log/slog"
	"net/http"

	"github.com/bansal-arsh/personal-finance-tracker/internal/email"
)

func HandleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Welcome"))
}

func HandleEmail(gd *email.GmailDialer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		receiver, ok := r.URL.Query()["email"]
		var receiverAddr string
		if ok {
			receiverAddr = receiver[0]
		} else {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
		}

		htmlBody := `
			<html>
				<body>
					<h1>This is a Test Email</h1>
					<p><b>Hello!</b> This is a test email with <a href='uwaterloo.ca'>Test link</a>.</p>
					<p>Thanks,<br>Mailtrap</p>
				</body>
			</html>
		`
		confirmEmail, err := email.NewEmail(receiverAddr, "Test", htmlBody, "uwaterloo test link")
		if err != nil {
			slog.Error("Error creating email message", "err", err)
			http.Error(w, "Error sending email", http.StatusInternalServerError)
		}

		err = gd.Send(confirmEmail)
		if err != nil {
			slog.Error("Error sending email message", "err", err)
			http.Error(w, "Error sending email", http.StatusInternalServerError)
		}

		w.Write([]byte("Email sent!"))
	}
}
