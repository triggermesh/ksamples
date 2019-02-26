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

//PubSubMessage struct taken from https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Attributes  map[string]string `json:"attributes"`
	Data        string            `json:"data"`
	MessageID   int               `json:"messageId"`
	PublishTime string            `json:"publishTime"`
}

//Handler handles events from GpcPubSub source
func Handler(ctx context.Context, message PubSubMessage) error {

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
	log.Info("data: ", message)

	repo := "testRepo" // with test repo name. Real will be obtained from data after the test with real life data
	branch := "master" // with test repo branch. Real will be obtained from data after the test with real life data

	filename := fmt.Sprintf("build/%s-%s.svg", repo, branch)

	if message.Attributes["status"] == "SUCCESS" {
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

	if message.Attributes["status"] == "FAILURE" {
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
