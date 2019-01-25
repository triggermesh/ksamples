package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
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

type Base64EncodedPubSubMessage string

func HandleRequest(ctx context.Context, pubsubMessage Base64EncodedPubSubMessage) (string, error) {

	// message := `{"message": {"attributes": {"buildId": "","status": "SUCCESS"}, "data": "SGVsbG8gQ2xvdWQgUHViL1N1YiEgSGVyZSBpcyBteSBtZXNzYWdlIQ==", "message_id": "136969346945"}, "subscription": "projects/myproject/subscriptions/mysubscription"}`
	// var data SubscriptionMessage
	// sEnc := base64.StdEncoding.EncodeToString([]byte(message))

	googlePubSubMessage, err := base64.StdEncoding.DecodeString(string(pubsubMessage))
	if err != nil {
		logrus.Fatal(err)
	}

	var data SubscriptionMessage
	err = json.Unmarshal([]byte(googlePubSubMessage), &data)
	if err != nil {
		logrus.Fatal(err)
	}

	// Need to check if this context is better used
	//ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		logrus.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new bucket.
	bucketName := "anatoliybucket"

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	repo := "testRepo"     // with test repo name. Real will be obtained from data after the test with real life data
	branch := "testBranch" // with test repo branch. Real will be obtained from data after the test with real life data

	filename := fmt.Sprintf("build/%s-%s.svg", repo, branch)

	if data.Message.Attributes.Status == "SUCCESS" {
		logrus.Info("Detected build success!")

		src := bucket.Object("build/success.svg")
		dst := bucket.Object(filename)

		if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Switched badge to build success")

		acl := bucket.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			logrus.Error(err)
		}
		logrus.Info("Badge set to public")
	}

	if data.Message.Attributes.Status == "FAILURE" {
		logrus.Info("Detected build failure!")

		src := bucket.Object("build/failure.svg")
		dst := bucket.Object(filename)

		if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
			logrus.Error(err)
		}

		logrus.Info("Switched badge to build failure")

		acl := bucket.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			logrus.Error(err)
		}

		logrus.Info("Badge set to public")
	}
}

func main() {
	lambda.Start(HandleRequest)
}
