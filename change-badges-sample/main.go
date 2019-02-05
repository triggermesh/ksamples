package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
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

func Handler() error {

	creds, err := base64.StdEncoding.DecodeString(os.Getenv("CREDENTIALS"))
	if err != nil {
		return err
	}

	//Set env variable to a file to be created
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "creds.json")

	err = ioutil.WriteFile("creds.json", []byte(creds), 0644)
	if err != nil {
		return err
	}

	var data SubscriptionMessage
	//test message to try create a badge
	message := `{"message": {"attributes": {"buildId": "","status": "SUCCESS"}, "data": "SGVsbG8gQ2xvdWQgUHViL1N1YiEgSGVyZSBpcyBteSBtZXNzYWdlIQ==", "message_id": "136969346945"}, "subscription": "projects/myproject/subscriptions/mysubscription"}`

	err = json.Unmarshal([]byte(message), &data)
	if err != nil {
		return err
	}

	fmt.Println("Message: ", data)

	ctx := context.Background()

	fmt.Println("Context Created")
	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Client Created")

	// Creates a Bucket instance.
	bucket := client.Bucket("tmbadges")

	repo := "testRepo" // with test repo name. Real will be obtained from data after the test with real life data
	branch := "master" // with test repo branch. Real will be obtained from data after the test with real life data

	filename := fmt.Sprintf("build/%s-%s.svg", repo, branch)

	if data.Message.Attributes.Status == "SUCCESS" {
		fmt.Println("Detected build success!")

		src := bucket.Object("build/success.svg")
		dst := bucket.Object(filename)

		if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
			return err
		}

		fmt.Println("Switched badge to build success")

		acl := bucket.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return err
		}
		fmt.Println("Badge set to public")
	}

	if data.Message.Attributes.Status == "FAILURE" {
		fmt.Println("Detected build failure!")

		src := bucket.Object("build/failure.svg")
		dst := bucket.Object(filename)

		if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
			return err
		}

		fmt.Println("Switched badge to build failure")

		acl := bucket.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return err
		}

		fmt.Println("Badge set to public")
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
