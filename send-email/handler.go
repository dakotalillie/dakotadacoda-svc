package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type handler struct {
	ssmSvc ssmiface.SSMAPI
}

func (h *handler) Run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)
	headers["Access-Control-Allow-Headers"] = "Content-Type"
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Methods"] = "OPTIONS,POST"

	config, err := NewEmailConfig(request.Body, h.ssmSvc)
	if codedError, ok := err.(*codedError); ok {
		return events.APIGatewayProxyResponse{StatusCode: codedError.Code, Headers: headers, Body: codedError.Error()}, codedError
	}

	mailer := Mailer{config}
	if err = mailer.Send(); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 502, Headers: headers, Body: err.Error()}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers}, nil
}
