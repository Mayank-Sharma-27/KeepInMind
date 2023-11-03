package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Reminder struct {
	RecipientEmail  string `json:"recipientEmail"`
	SentBy          string `json:"sentBy"`
	Ics             string `json:"ics"`
	ReminderMessage string `json:"reminderMessage"`
}

const (
	SenderEmail = "mayanksharma.sharma77@gmail.com"
	Charset     = "UTF-8"
)

func sendReminder(w http.ResponseWriter, r *http.Request) {
	log.Println("sendReminder: request received")

	var reminder Reminder
	err := json.NewDecoder(r.Body).Decode(&reminder)
	if err != nil {
		log.Printf("sendReminder: error decoding reminder: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Json decoded : ", reminder)

	// Initialize a session with AWS SDK
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	log.Println("AWS session created")

	// Create SES client
	svc := ses.New(sess)

	// Email body and subject
	subject := fmt.Sprintf("Reminder from %s", reminder.SentBy)
	bodyText := reminder.ReminderMessage

	// Create MIME email with attached calendar invite
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", "text/calendar; charset=UTF-8; method=REQUEST")
	header.Set("Content-Disposition", "attachment; filename=invite.ics")
	mime := &bytes.Buffer{}
	wr := multipart.NewWriter(mime)
	part, _ := wr.CreatePart(header)
	part.Write([]byte(reminder.Ics))
	wr.Close()

	rawEmail := fmt.Sprintf(
		`From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="%s"

--%s
Content-Type: text/plain; charset=UTF-8

%s

--%s
%s
`,
		SenderEmail, reminder.RecipientEmail, subject, wr.Boundary(), wr.Boundary(), bodyText, wr.Boundary(), mime.String(),
	)

	input := &ses.SendRawEmailInput{
		Destinations: []*string{aws.String(reminder.RecipientEmail)},
		RawMessage:   &ses.RawMessage{Data: []byte(rawEmail)},
		Source:       aws.String(SenderEmail),
	}

	result, err := svc.SendRawEmail(input)
	if err != nil {
		log.Printf("sendReminder: error sending email: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("sendReminder: email sent successfully with Message ID: %s", *result.MessageId)

	// Send a response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Email sent with Message ID: %s", *result.MessageId)))
}
