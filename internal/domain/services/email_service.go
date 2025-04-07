package services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"clothing-shop-api/internal/config"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{config: cfg}
}

type EmailData struct {
	Subject   string
	To        string
	Template  string
	Variables map[string]interface{}
}

func (s *EmailService) SendEmail(data EmailData) error {
	// Parse template
	tmpl, err := template.New("email").Parse(data.Template)
	if err != nil {
		return err
	}

	// Fill template with variables
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data.Variables); err != nil {
		return err
	}

	// Set up authentication
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	// Set up message
	from := fmt.Sprintf("%s <%s>", s.config.SMTPFromName, s.config.SMTPFromEmail)
	to := []string{data.To}

	header := make(map[string]string)
	header["From"] = from
	header["To"] = data.To
	header["Subject"] = data.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	// Send email
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.SMTPFromEmail,
		to,
		[]byte(message),
	)

	return err
}

func (s *EmailService) SendVerificationEmail(email, token string) error {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.config.AppURL, token)

	template := `
    <html>
    <body>
        <h1>Email Verification</h1>
        <p>Thank you for registering with our Clothing Shop. Please click the link below to verify your email address:</p>
        <p><a href="{{.VerificationURL}}">Verify Email</a></p>
        <p>If you did not register for an account, you can ignore this email.</p>
    </body>
    </html>
    `

	return s.SendEmail(EmailData{
		Subject:  "Verify Your Email Address",
		To:       email,
		Template: template,
		Variables: map[string]interface{}{
			"VerificationURL": verificationURL,
		},
	})
}

func (s *EmailService) SendPasswordResetEmail(email, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.AppURL, token)

	template := `
    <html>
    <body>
        <h1>Password Reset</h1>
        <p>You requested to reset your password. Please click the link below to set a new password:</p>
        <p><a href="{{.ResetURL}}">Reset Password</a></p>
        <p>If you did not request a password reset, please ignore this email.</p>
    </body>
    </html>
    `

	return s.SendEmail(EmailData{
		Subject:  "Reset Your Password",
		To:       email,
		Template: template,
		Variables: map[string]interface{}{
			"ResetURL": resetURL,
		},
	})
}
