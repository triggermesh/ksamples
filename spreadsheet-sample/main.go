package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

//DynamoDBEvent represents AWS Dynamo DB payload
type DynamoDBEvent struct {
	AwsRegion    *string                       `locationName:"awsRegion" type:"string"`
	Dynamodb     *dynamodbstreams.StreamRecord `locationName:"dynamodb" type:"structure"`
	EventID      *string                       `locationName:"eventID" type:"string"`
	EventName    *string                       `locationName:"eventName" type:"string" enum:"OperationType"`
	EventSource  *string                       `locationName:"eventSource" type:"string"`
	EventVersion *string                       `locationName:"eventVersion" type:"string"`
	UserIdentity *dynamodbstreams.Identity     `locationName:"userIdentity" type:"structure"`
}

func Handler(ctx context.Context, event DynamoDBEvent) error {

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	_, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")

	if !ok {
		env := os.Getenv("ENV")
		fmt.Println(env)

		creds, err := base64.StdEncoding.DecodeString(os.Getenv("CREDENTIALS"))
		if err != nil {
			log.Error(err)
		}

		err = ioutil.WriteFile("credentials.json", []byte(creds), 0644)
		if err != nil {
			log.Error(err)
		}

		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "credentials.json")
	}

	ctx = context.Background()

	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return err
	}

	srv, err := sheets.New(client)
	if err != nil {
		return err
	}

	log.Info("Authorized via service account. Listen for DynamoDB events!")

	var vr sheets.ValueRange

	myval := []interface{}{}

	for _, v := range event.Dynamodb.NewImage {
		log.Info(v.String())

		if v.B != nil {
			myval = append(myval, v.B)
		} else if v.BOOL != nil {
			myval = append(myval, *v.BOOL)
		} else if v.BS != nil {
			myval = append(myval, v.BS)
		} else if v.L != nil {
			myval = append(myval, v.L)
		} else if v.M != nil {
			myval = append(myval, v.M)
		} else if v.N != nil {
			myval = append(myval, *v.N)
		} else if v.NS != nil {
			myval = append(myval, v.NS)
		} else if v.NULL != nil {
			myval = append(myval, *v.NULL)
		} else if v.S != nil {
			myval = append(myval, *v.S)
		} else if v.SS != nil {
			myval = append(myval, v.BS)
		}
	}

	log.Info(myval)

	vr.Values = append(vr.Values, myval)

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, "A1", &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
