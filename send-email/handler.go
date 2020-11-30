package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type handler struct {
	ssmSvc ssmiface.SSMAPI
}

func (h *handler) Run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	params, err := newSendMailParams(request.Body, h.ssmSvc)
	if codedError, ok := err.(*codedError); ok {
		code := codedError.Code
		if code == 0 {
			code = 500
		}
		return events.APIGatewayProxyResponse{StatusCode: code, Body: codedError.Error()}, codedError
	}

	if err = sendMail(&params); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
