package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

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

func HandleRequest(ctx context.Context, sm SubscriptionMessage) (string, error) {
	sDec, err := b64.StdEncoding.DecodeString(sm.Message.Data)
	if err != nil {
		return err.Error(), nil
	}

	repo := "testRepo"     // with test repo name. Real will be obtained from data after the test with real life data
	branch := "testBranch" // with test repo branch. Real will be obtained from data after the test with real life data

	filename := fmt.Sprintf("build/%s-%s.svg", repo, branch)

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new bucket.
	bucketName := "my-new-bucket"

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	if sm.Message.Attributes.Status == "SUCCESS" {
		logrus.Info("Detected build success!")

		f, err := os.Open("build/success.svg")
		if err != nil {
			return "", err
		}
		defer f.Close()

		wc := bucket.Object(filename).NewWriter(ctx)
		if _, err = io.Copy(wc, f); err != nil {
			return "", err
		}
		if err := wc.Close(); err != nil {
			return "", err
		}

		acl := bucket.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return "", err
		}
	}

	if sm.Message.Attributes.Status == "FAILURE" {
		logrus.Info("Detected build failure!")

		f, err := os.Open("build/failure.svg")
		if err != nil {
			return "", err
		}
		defer f.Close()

		wc := bucket.Object(filename).NewWriter(ctx)
		if _, err = io.Copy(wc, f); err != nil {
			return "", err
		}
		if err := wc.Close(); err != nil {
			return "", err
		}

		acl := bucket.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("Build status: %s, data: %s", sm.Message.Attributes.Status, string(sDec)), nil
}

func main() {
	lambda.Start(HandleRequest)
}

//curl http://go-lambda.fedorenkotolik.dev.triggermesh.io -d @testdata.json --header "Content-Type: application/json"
