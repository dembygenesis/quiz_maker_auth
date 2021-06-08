package email

import (
	"gopkg.in/gomail.v2"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	"strconv"
)


func SendMail(
	to string,
	subject string,
	message string,
) error {
	m := gomail.NewMessage()

	m.SetHeader("From", config.From)
	m.SetHeader("To", "dembygenesis@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	parsedEmailPort, _ := strconv.Atoi(config.EmailPort)

	d := gomail.NewDialer(config.EmailHost, parsedEmailPort, config.EmailUsername, config.EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	} else {
		return err
	}

	return nil
}