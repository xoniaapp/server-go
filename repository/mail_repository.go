package repository

import (
	"fmt"
	"net/smtp"

	"github.com/xoniaapp/server/model"
)

type mailRepository struct {
	username string
	password string
	origin   string
}

func NewMailRepository(username string, password string, origin string) model.MailRepository {
	return &mailRepository{
		username: username,
		password: password,
		origin:   origin,
	}
}

func (m *mailRepository) SendResetMail(email string, token string) error {

	msg := "From: " + m.username + "\n" +
		"To: " + email + "\n" +
		"Subject: Reset Email\n\n" +
		fmt.Sprintf("<a href=\"%s/reset-password/%s\">Reset Password</a>", m.origin, token)

	err := smtp.SendMail("smtp.eu.sparkpostmail.com:2525",
		smtp.PlainAuth("", m.username, m.password, "smtp.eu.sparkpostmail.com:2525"),
		m.username, []string{email}, []byte(msg))

	return err
}
