package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var reqBodyVals map[string]string

func setup(t *testing.T) func(t *testing.T) {
	os.Setenv("FROM_ADDRESS", "dakota@test.com")
	os.Setenv("TO_ADDRESS", "dakota@test.com")
	os.Setenv("SMTP_HOST", "smtp.mail.com")
	os.Setenv("SMTP_PORT", "587")

	reqBodyVals = make(map[string]string)
	reqBodyVals["name"] = "Jane"
	reqBodyVals["email"] = "jane@test.com"
	reqBodyVals["subject"] = "Test subject"
	reqBodyVals["message"] = "Test message"

	return func(t *testing.T) {
		os.Unsetenv("FROM_ADDRESS")
		os.Unsetenv("TO_ADDRESS")
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")

		reqBodyVals = make(map[string]string)
	}
}

func TestHandler(t *testing.T) {

	t.Run("No from address", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		os.Unsetenv("FROM_ADDRESS")

		res, err := handler(events.APIGatewayProxyRequest{})

		if res.StatusCode != 500 {
			t.Fatalf("status code should be 500, got %v", res.StatusCode)
		} else if res.Body != "Missing from address" {
			t.Fatalf("response body should be \"Missing from address\", got %v", res.Body)
		} else if err == nil {
			t.Fatal("err should not be nil")
		}
	})

	t.Run("No to address", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		os.Unsetenv("TO_ADDRESS")

		res, err := handler(events.APIGatewayProxyRequest{})

		if res.StatusCode != 500 {
			t.Fatalf("status code should be 500, got %v", res.StatusCode)
		} else if res.Body != "Missing to address" {
			t.Fatalf("response body should be \"Missing to address\", got \"%v\"", res.Body)
		} else if err == nil {
			t.Fatal("err should not be nil")
		}
	})

	t.Run("No smtp host", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		os.Unsetenv("SMTP_HOST")

		res, err := handler(events.APIGatewayProxyRequest{})

		if res.StatusCode != 500 {
			t.Fatalf("status code should be 500, got %v", res.StatusCode)
		} else if res.Body != "Missing host" {
			t.Fatalf("response body should be \"Missing host\", got \"%v\"", res.Body)
		} else if err == nil {
			t.Fatal("err should not be nil")
		}
	})

	t.Run("No smtp port", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		os.Unsetenv("SMTP_PORT")

		res, err := handler(events.APIGatewayProxyRequest{})

		if res.StatusCode != 500 {
			t.Fatalf("status code should be 500, got %v", res.StatusCode)
		} else if res.Body != "Missing port" {
			t.Fatalf("response body should be \"Missing port\", got \"%v\"", res.Body)
		} else if err == nil {
			t.Fatal("err should not be nil")
		}
	})

	t.Run("Request body is not valid JSON", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		res, err := handler(events.APIGatewayProxyRequest{Body: "{\"test: wut}"})

		if res.StatusCode != 400 {
			t.Fatalf("status code should be 400, got %v", res.StatusCode)
		} else if res.Body != "Unable to unmarshal request body" {
			t.Fatalf("response body should be \"Unable to unmarshal request body\", got \"%v\"", res.Body)
		} else if err == nil {
			t.Fatal("err should not be nil")
		}
	})

	for _, key := range []string{"name", "email", "subject", "message"} {
		t.Run(fmt.Sprintf("Request body is missing %v field", key), func(t *testing.T) {
			teardown := setup(t)
			defer teardown(t)
			delete(reqBodyVals, key)
			reqBody, err := json.Marshal(reqBodyVals)
			if err != nil {
				t.Fatal("Unable to marshal request body")
			}

			res, err := handler(events.APIGatewayProxyRequest{Body: string(reqBody)})

			expectedResBody := fmt.Sprintf("Missing %s", key)
			if res.StatusCode != 400 {
				t.Fatalf("status code should be 400, got %v", res.StatusCode)
			} else if res.Body != expectedResBody {
				t.Fatalf("response body should be \"%v\", got \"%v\"", expectedResBody, res.Body)
			} else if err == nil {
				t.Fatal("err should not be nil")
			}
		})
	}
}
