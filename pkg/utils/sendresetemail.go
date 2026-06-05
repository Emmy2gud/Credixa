package utils

import (
	"fmt"

	"os"

	"github.com/resend/resend-go/v3"
)

func SendResetEmail(toEmail, resetLink string) error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
   params := &resend.SendEmailRequest{
        From:    "Acme <onboarding@resend.dev>",
        To:      []string{toEmail},
        Html: "<div style='font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e9e9e9; border-radius: 10px;'>"+
                "<h2 style='color: #333;'>Password Reset Request</h2>"+
                "<p style='color: #555; line-height: 1.6;'>"+
                    "You are receiving this email because we received a password reset request for your account."+ 
                "</p>"+
                "<div style='text-align: center; margin: 30px 0;'>"+
                    "<a href='" + resetLink + "' style='background-color: #000666; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; font-weight: bold;'>"+
                        "Reset Password"+
                    "</a>"+
                "</div>"+
                "<p style='color: #555; line-height: 1.6;'>"+
                    "This password reset link will expire in 60 minutes."+
                "</p>"+
                "<p style='color: #777; font-size: 12px; margin-top: 40px; border-top: 1px solid #eee; padding-top: 20px;'>"+
                    "If you did not request a password reset, no further action is required."+
                "</p>"+
            "</div>",
        Subject: "Password Reset",

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
