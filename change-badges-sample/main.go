/*
Copyright (c) 2019 TriggerMesh, Inc
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	log "github.com/sirupsen/logrus"
)

type PubSubMessage struct {
	Attributes  map[string]string `json:"attributes"`
	Data        string            `json:"data"`
	ID          int               `json:"ID"`
	PublishTime string            `json:"publishTime"`
}

type PubSubPayload struct {
	Status string
	Source struct {
		RepoSource struct {
			ProjectID  string
			RepoName   string
			BranchName string
		}
	}
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

	payload, err := parseData(message.Data)
	if err != nil {
		return err
	}

	repo := payload.Source.RepoSource.RepoName
	branch := payload.Source.RepoSource.BranchName

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

func parseData(str string) (PubSubPayload, error) {
	var eventPayload PubSubPayload

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return eventPayload, err
	}

	err = json.Unmarshal(data, &eventPayload)
	if err != nil {
		return eventPayload, err
	}

	return eventPayload, nil
}

func main() {
	lambda.Start(Handler)
}
