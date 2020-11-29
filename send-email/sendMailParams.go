package main

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type sendMailParams struct {
	From     string
	To       string
	Host     string
	Port     string
	Password string
}

func newSendMailParams() (sendMailParams, error) {
	params := sendMailParams{}

	from := os.Getenv("FROM_ADDRESS")
	if from == "" {
		return params, errors.New("Missing from address")
	}

	to := os.Getenv("TO_ADDRESS")
	if to == "" {
		return params, errors.New("Missing to address")
	}

	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return params, errors.New("Missing host")
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		return params, errors.New("Missing port")
	}

	password, err := getEmailPassword()
	if err != nil {
		return params, err
	}

	params.From = from
	params.To = to
	params.Host = host
	params.Port = port
	params.Password = password

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
