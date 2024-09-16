package mailer

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

var client *resend.Client

func Init(key string) {
	if key == "" {
		panic("RESEND_API_KEY is not set")
	}

	client = resend.NewClient(key)
}

func SendConfirmationEmail(to string, code string) error {
	_, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{to},
		Html:    fmt.Sprintf("<h1>%s</h1>", code),
		Subject: "Confirm your email",
	})

	return err
}

func SendPasswordResetEmail(to string, code string) error {
	_, err := client.Emails.Send(&resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{to},
		Html:    fmt.Sprintf("<h1>%s</h1>", code),
		Subject: "Confirm your email",
	})

	return err
}
