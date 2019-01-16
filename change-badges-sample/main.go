package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type SubscriptionMessage struct {
	Message      Message `json: "message"`
	Subscription string  `json: "subscription"`
}

type Message struct {
	Attributes Attributes `json: "attributes"`
	Data       string     `json: "data"`
	MessageID  int        `json: "message_id"`
}

type Attributes struct {
	BuildID string `json: "buildId"`
	Status  string `json: "status"`
}

func HandleRequest(ctx context.Context, sm SubscriptionMessage) (string, error) {
	sDec, err := b64.StdEncoding.DecodeString(sm.Message.Data)
	if err != nil {
		return err.Error(), nil
	}

	return fmt.Sprintf("Build status: %s, data: %s", sm.Message.Attributes.Status, string(sDec)), nil
}

func main() {
	lambda.Start(HandleRequest)
}

//curl http://go-lambda.fedorenkotolik.dev.triggermesh.io -d @testdata.json --header "Content-Type: application/json"
