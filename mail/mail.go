package mail

import (
	"fmt"
	"net/smtp"
)

func SendMailGmail(from, to, password, subject, body string, isHtml bool) {
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"

	if !isHtml {
		mime = ""
	}

	toList := []string{to}
	msg := []byte(fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		mime +
		fmt.Sprintf("%s\r\n", body))
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, toList, msg)
	if err != nil {
		fmt.Println(err)
	}
}
