package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
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

func Handler(ctx context.Context, data SubscriptionMessage) error {

	bucketName := os.Getenv("BUCKET")
	log.Info("Current bucket: ", bucketName)

	_, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")

	if !ok {
		fmt.Println("Unable to find GOOGLE_APPLICATION_CREDENTIALS env var. Creating it locally")

		credentials := os.Getenv("CREDENTIALS")

		log.Info(credentials)

		err := ioutil.WriteFile("credentials.json", []byte(credentials), 0644)
		if err != nil {
			return err
		}

		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "credentials.json")
	}

	log.Info("Configured env variables and google credentials!")

	ctx = context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	log.Info("Client Created")

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	log.Info("Bucket: ", bucket)
	log.Info("data: ", data)

	repo := "testRepo" // with test repo name. Real will be obtained from data after the test with real life data
	branch := "master" // with test repo branch. Real will be obtained from data after the test with real life data

	filename := fmt.Sprintf("build/%s-%s.svg", repo, branch)

	if data.Message.Attributes.Status == "SUCCESS" {
		log.Info("Detected build success!")

		src := bucket.Object("build/success.svg")
		dst := bucket.Object(filename)

		if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
			return err
		}

		log.Info("Switched badge to build success")

		acl := bucket.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return err
		}
		log.Info("Badge set to public")
	}

	if data.Message.Attributes.Status == "FAILURE" {
		log.Info("Detected build failure!")

		src := bucket.Object("build/failure.svg")
		dst := bucket.Object(filename)

		if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
			return err
		}

		log.Info("Switched badge to build failure")

		acl := bucket.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return err
		}

		log.Info("Badge set to public")
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
