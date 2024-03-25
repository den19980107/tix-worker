package mail

import (
	"fmt"
	"net/smtp"
)

type Mail struct {
	from     string
	smtpUrl  string
	smtpAuth smtp.Auth
}

func New(from string, password string, smtpHost string, smtpPort string) Mail {
	auth := smtp.PlainAuth("", from, password, smtpHost)
	smtpUrl := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return Mail{
		from:     from,
		smtpUrl:  smtpUrl,
		smtpAuth: auth,
	}
}

func (m *Mail) Send(to string, message string) error {
	msg := m.composeEmail(to, message)
	return smtp.SendMail(m.smtpUrl, m.smtpAuth, m.from, []string{to}, msg)
}

func (m *Mail) composeEmail(to string, message string) []byte {
	msg := []byte(fmt.Sprintf("To: %s\r\n"+

		"Subject: 驗證碼填寫通知\r\n"+

		"\r\n"+

		"%s\r\n", to, message))

	return msg
}
