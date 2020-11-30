package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type handler struct {
	ssmSvc ssmiface.SSMAPI
}

func (h *handler) Run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	config, err := NewEmailConfig(request.Body, h.ssmSvc)
	if codedError, ok := err.(*codedError); ok {
		return events.APIGatewayProxyResponse{StatusCode: codedError.Code, Body: codedError.Error()}, codedError
	}

	if err = sendMail(&config); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 502, Body: err.Error()}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
