package main

import (
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
	log.Info("Current spreadsheet id: ", spreadsheetID)

	_, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")

	if !ok {
		fmt.Println("Unable to find GOOGLE_APPLICATION_CREDENTIALS env var. Creating it locally")

		credentials := os.Getenv("CREDENTIALS")

		log.Info("Credentials: \n")
		log.Info(credentials)

		err := ioutil.WriteFile("credentials.json", []byte(credentials), 0644)
		if err != nil {
			log.Fatal(err)
		}

		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "credentials.json")
	}

	log.Info("Configured env variables and google credentials!")

	ctx = context.Background()

	client, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return err
	}

	log.Info("Created Google Default client")

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

	log.Info("Adding new value to the spreadsheet: ", myval)

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
