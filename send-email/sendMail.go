package main

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

func sendMail(params *sendMailParams) error {
	from := mail.Address{Name: "Dakota Lillie", Address: params.From}
	to := mail.Address{Name: "Dakota Lillie", Address: params.To}
	subj := "This is the email subject"
	body := "This is an example body"

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := fmt.Sprintf("%s:%s", params.Host, params.Port)

	auth := smtp.PlainAuth("", params.From, params.Password, params.Host)

	tlsconfig := &tls.Config{ServerName: params.Host}

	client, err := smtp.Dial(servername)
	if err != nil {
		return err
	}

	client.StartTLS(tlsconfig)

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from.Address); err != nil {
		return err
	}

	if err = client.Rcpt(to.Address); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	client.Quit()

	return nil
}
