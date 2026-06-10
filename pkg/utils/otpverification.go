package utils

import (
	"fmt"

	"os"

	"github.com/resend/resend-go/v3"
)

func SendOTPEmail(toEmail, otp string) error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
   params := &resend.SendEmailRequest{
        From:    "Acme <onboarding@resend.dev>",
        To:      []string{toEmail},
        Html: `
			<h2>Email Verification</h2>
			<p>Your OTP is:</p>
			<h1>`+otp+`</h1>
			<p>Expires in 10 minutes.</p>
		`,
        Subject: "Verify your account",

        ReplyTo: "replyto@example.com",
    }

    sent, err := client.Emails.Send(params)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
    fmt.Println(sent.Id)
	return nil
}
