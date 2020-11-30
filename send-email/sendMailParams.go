package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type sendMailParams struct {
	From     string // The email address the email will be sent from
	To       string // The email address the email will be sent to
	Host     string // The SMTP host for the email server
	Port     string // The SMTP port for the email server
	Password string // The password for the sending account
	Name     string // The name of the person who is sending the email
	Email    string // The email of the person who is sending the email, used for the reply to
	Subject  string // The subject of the email
	Message  string // The message body of the email
}

func newSendMailParams(reqBody string) (sendMailParams, error) {
	params := sendMailParams{}

	from := os.Getenv("FROM_ADDRESS")
	if from == "" {
		return params, &codedError{Code: 500, Message: "Missing from address"}
	}

	to := os.Getenv("TO_ADDRESS")
	if to == "" {
		return params, &codedError{Code: 500, Message: "Missing to address"}
	}

	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return params, &codedError{Code: 500, Message: "Missing host"}
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		return params, &codedError{Code: 500, Message: "Missing port"}
	}

	parsedBody := make(map[string]string)
	err := json.Unmarshal([]byte(reqBody), &parsedBody)
	if err != nil {
		return params, &codedError{Code: 400, Message: "Unable to unmarshal request body"}
	}

	for _, key := range []string{"name", "email", "subject", "message"} {
		if parsedBody[key] == "" {
			return params, &codedError{Code: 400, Message: "Missing " + key}
		}
	}

	password, err := getEmailPassword()
	if err != nil {
		return params, &codedError{Code: 500, Message: "Unable to get email password"}
	}

	params.From = from
	params.To = to
	params.Host = host
	params.Port = port
	params.Password = password
	params.Name = parsedBody["name"]
	params.Email = parsedBody["email"]
	params.Subject = parsedBody["subject"]
	params.Message = parsedBody["message"]

	return params, nil
}

func getEmailPassword() (string, error) {
	sess := session.Must(session.NewSession())
	svc := ssm.New(sess)
	res, err := svc.GetParameter(&ssm.GetParameterInput{Name: aws.String("/DakotaDaCoda/EMAIL_PASSWORD"), WithDecryption: aws.Bool(true)})
	if err != nil {
		return "", err
	}
	return *res.Parameter.Value, nil
}
