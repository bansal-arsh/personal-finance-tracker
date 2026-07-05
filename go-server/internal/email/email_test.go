package email

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"reflect"
	"strings"
	"testing"
)

func TestNewEmail(t *testing.T) {
	type inputParam struct {
		receiver string
		subject  string
		htmlBody string
		textBody string
	}
	type expectedResult struct {
		email *Email
		err   error
	}
	type testCases struct {
		name     string
		input    inputParam
		expected expectedResult
	}

	testHtmlBody := `
		<html>
			<body>
				<p>Hello!</p>
			</body>
		</html>
	`

	returnErr := func(address string) error {
		_, err := mail.ParseAddress(address)
		return err
	}
	returnAddress := func(address string) *mail.Address {
		parsedAddress, _ := mail.ParseAddress(address)
		return parsedAddress
	}

	tests := []testCases{
		{
			name: "Normal request",
			input: inputParam{
				receiver: "abc@example.com",
				subject:  "Hello",
				htmlBody: testHtmlBody,
				textBody: "Hello!",
			},
			expected: expectedResult{
				email: &Email{
					receiver: *returnAddress("abc@example.com"),
					subject:  "Hello",
					htmlBody: testHtmlBody,
					textBody: "Hello!",
				},
				err: nil,
			},
		},
		{
			name: "Empty receiver address",
			input: inputParam{
				receiver: "",
				subject:  "Hello",
				htmlBody: testHtmlBody,
				textBody: "Hello!",
			},
			expected: expectedResult{
				email: nil,
				err:   ErrNoRecevier,
			},
		},
		{
			name: "Invalid receiver address",
			input: inputParam{
				receiver: "abc.123",
				subject:  "Hello",
				htmlBody: testHtmlBody,
				textBody: "Hello!",
			},
			expected: expectedResult{
				email: nil,
				err:   returnErr("abc.123"),
			},
		},
		{
			name: "Empty subject",
			input: inputParam{
				receiver: "abc@example.com",
				subject:  "",
				htmlBody: testHtmlBody,
				textBody: "Hello!",
			},
			expected: expectedResult{
				email: &Email{
					receiver: *returnAddress("abc@example.com"),
					subject:  "",
					htmlBody: testHtmlBody,
					textBody: "Hello!",
				},
				err: nil,
			},
		},
		{
			name: "Empty Html body",
			input: inputParam{
				receiver: "abc@example.com",
				subject:  "Hello",
				htmlBody: "",
				textBody: "Hello!",
			},
			expected: expectedResult{
				email: &Email{
					receiver: *returnAddress("abc@example.com"),
					subject:  "Hello",
					htmlBody: "",
					textBody: "Hello!",
				},
				err: nil,
			},
		},
		{
			name: "Empty Text body",
			input: inputParam{
				receiver: "abc@example.com",
				subject:  "Hello",
				htmlBody: testHtmlBody,
				textBody: "",
			},
			expected: expectedResult{
				email: nil,
				err:   ErrNoTextBody,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.input
			email, err := NewEmail(input.receiver, input.subject, input.htmlBody, input.textBody)

			expected := test.expected
			if !reflect.DeepEqual(email, expected.email) {
				t.Errorf("Wrong email\nGot: %+v\nWant: %+v", email, expected.email)
			}
			if !errors.Is(err, expected.err) && err.Error() != expected.err.Error() {
				t.Errorf("Wrong error\nGot: %q\nWant: %q", err, expected.err)
			}
		})
	}
}

func TestNewGmailDialer(t *testing.T) {
	type inputParam struct {
		sender      string
		appPassword string
	}
	type expectedResult struct {
		gd  *GmailDialer
		err error
	}
	type testCases struct {
		name     string
		input    inputParam
		expected expectedResult
	}

	returnErr := func(address string) error {
		_, err := mail.ParseAddress(address)
		return err
	}
	returnAddress := func(address string) *mail.Address {
		parsedAddress, _ := mail.ParseAddress(address)
		return parsedAddress
	}

	tests := []testCases{
		{
			name: "Normal request",
			input: inputParam{
				sender:      "abc@example.com",
				appPassword: "dsafghbvfruikjhg",
			},
			expected: expectedResult{
				gd: &GmailDialer{
					sender:      *returnAddress("abc@example.com"),
					appPassword: "dsafghbvfruikjhg",
				},
				err: nil,
			},
		},
		{
			name: "Empty sender address",
			input: inputParam{
				sender:      "",
				appPassword: "dsafghbvfruikjhg",
			},
			expected: expectedResult{
				gd:  nil,
				err: ErrNoSender,
			},
		},
		{
			name: "Invalid sender address",
			input: inputParam{
				sender:      "abc.123",
				appPassword: "dsafghbvfruikjhg",
			},
			expected: expectedResult{
				gd:  nil,
				err: returnErr("abc.123"),
			},
		},
		{
			name: "Empty app password",
			input: inputParam{
				sender:      "abc@example.com",
				appPassword: "",
			},
			expected: expectedResult{
				gd:  nil,
				err: ErrNoGmailAppPassword,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.input
			gmailDialer, err := NewGmailDialer(input.sender, input.appPassword)

			expected := test.expected
			if !reflect.DeepEqual(gmailDialer, expected.gd) {
				t.Errorf("Wrong email\nGot: %+v\nWant: %+v", gmailDialer, expected.gd)
			}
			if !errors.Is(err, expected.err) && err.Error() != expected.err.Error() {
				t.Errorf("Wrong error\nGot: %q\nWant: %q", err, expected.err)
			}
		})
	}
}

func TestBuildMessageWithDialer_WithHtmlBody(t *testing.T) {
	testReceiverAddress := "receive@example.com"
	testSubject := "Test subject"
	testTextBody := "Hello!"

	testHtmlBody := `
	<html>
		<body>
			<p>Hello!</p>
		</body>
	</html>
	`
	for _, removalChar := range []string{"\n", "\t", "\r"} {
		testHtmlBody = strings.ReplaceAll(testHtmlBody, removalChar, "")
	}

	testEmail, err := NewEmail(testReceiverAddress, testSubject, testHtmlBody, testTextBody)
	if err != nil {
		t.Fatalf("Error creating test email: %s", err)
	}

	testSenderAddress := "sender@example.com"
	testAppPassword := "ghkexpynkdewjhbf"
	testGmailDialer, err := NewGmailDialer(testSenderAddress, testAppPassword)
	if err != nil {
		t.Fatalf("Error creating test gmail dialer: %s", err)
	}

	message := testEmail.buildMessageWithDialer(testGmailDialer)

	var buf bytes.Buffer
	if _, err := message.WriteTo(&buf); err != nil {
		t.Fatalf("Error writing message to bytes buffer: %s", err)
	}

	parsedMsg, err := mail.ReadMessage(&buf)
	if err != nil {
		t.Fatalf("Error parsing message: %s", err)
	}

	if got := parsedMsg.Header.Get("From"); got != testSenderAddress {
		t.Errorf("Incorrect sender address\nGot: %q\nWant: %q", got, testSenderAddress)
	}
	if got := parsedMsg.Header.Get("To"); got != testReceiverAddress {
		t.Errorf("Incorrect sender address\nGot: %q\nWant: %q", got, testReceiverAddress)
	}
	if got := parsedMsg.Header.Get("Subject"); got != testSubject {
		t.Errorf("Incorrect sender address\nGot: %q\nWant: %q", got, testSubject)
	}

	// Parse the multipart body and check each part's exact content.
	mediaType, params, err := mime.ParseMediaType(parsedMsg.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("Bad content type: %v", err)
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		t.Fatalf("Expected multipart body, got %q", mediaType)
	}

	parts := map[string]string{}
	mr := multipart.NewReader(parsedMsg.Body, params["boundary"])
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("multipart read error: %v", err)
		}
		partType, _, _ := mime.ParseMediaType(p.Header.Get("Content-Type"))
		body, _ := io.ReadAll(p)
		parts[partType] = strings.TrimSpace(string(body))
	}

	if len(parts) != 2 {
		t.Fatalf("Expected exactly 2 parts, got %d parts: %v", len(parts), parts)
	}
	if got := parts["text/plain"]; got != testTextBody {
		t.Errorf("Wrong text/plain part\n Got: %q\nWant: %q", got, testTextBody)
	}
	if got := parts["text/html"]; got != testHtmlBody {
		t.Errorf("Wrong text/html part\n Got: %q\nWant: %q", got, testHtmlBody)
	}
}

func TestBuildMessageWithDialer_NoHtmlBody(t *testing.T) {
	testReceiverAddress := "receive@example.com"
	testSubject := "Test subject"
	testTextBody := "Hello!"

	testEmail, err := NewEmail(testReceiverAddress, testSubject, "", testTextBody)
	if err != nil {
		t.Fatalf("Error creating test email: %s", err)
	}

	testSenderAddress := "sender@example.com"
	testAppPassword := "ghkexpynkdewjhbf"
	testGmailDialer, err := NewGmailDialer(testSenderAddress, testAppPassword)
	if err != nil {
		t.Fatalf("Error creating test gmail dialer: %s", err)
	}

	message := testEmail.buildMessageWithDialer(testGmailDialer)

	var buf bytes.Buffer
	if _, err := message.WriteTo(&buf); err != nil {
		t.Fatalf("Error writing message to bytes buffer: %s", err)
	}

	parsedMsg, err := mail.ReadMessage(&buf)
	if err != nil {
		t.Fatalf("Error parsing message: %s", err)
	}

	if got := parsedMsg.Header.Get("From"); got != testSenderAddress {
		t.Errorf("Incorrect sender address\nGot: %q\nWant: %q", got, testSenderAddress)
	}
	if got := parsedMsg.Header.Get("To"); got != testReceiverAddress {
		t.Errorf("Incorrect sender address\nGot: %q\nWant: %q", got, testReceiverAddress)
	}
	if got := parsedMsg.Header.Get("Subject"); got != testSubject {
		t.Errorf("Incorrect sender address\nGot: %q\nWant: %q", got, testSubject)
	}

	// Parse the multipart body and check each part's exact content.
	mediaType, _, err := mime.ParseMediaType(parsedMsg.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("Bad content type: %v", err)
	}
	if mediaType != "text/plain" {
		t.Fatalf("Expected text/plain body, got %q", mediaType)
	}

	body, err := io.ReadAll(parsedMsg.Body)
	if err != nil {
		t.Fatalf("Error reading body: %s", err)
	}
	if !bytes.Equal(body, []byte(testTextBody)) {
		t.Errorf("Wrong body\nGot: %q\nWant: %q", string(body), testTextBody)
	}
}
