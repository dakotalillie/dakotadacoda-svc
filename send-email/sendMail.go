package main

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

func sendMail(config *EmailConfig) error {
	from := mail.Address{Name: "Dakota Lillie", Address: config.From}
	to := mail.Address{Name: "Dakota Lillie", Address: config.To}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = config.Subject
	headers["Reply-To"] = config.Email

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += fmt.Sprintf(
		"\r\nNew message received from %s via DakotaDaCoda:\n\n%s", config.Name, config.Message,
	)

	servername := fmt.Sprintf("%s:%s", config.Host, config.Port)
	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)
	tlsconfig := &tls.Config{ServerName: config.Host}

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
