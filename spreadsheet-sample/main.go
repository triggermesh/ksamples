package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var spreadsheetID string

type SpreadsheetDumper struct {
	Sheets *sheets.Service
}

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

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (ssd *SpreadsheetDumper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event DynamoDBEvent
	err := decoder.Decode(&event)
	if err != nil {
		panic(err)
	}

	var vr sheets.ValueRange

	myval := []interface{}{}

	resp := ssd.Sheets.Spreadsheets.Values.Get(spreadsheetID, "A1")
	log.Info(resp)

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

	_, err = ssd.Sheets.Spreadsheets.Values.Append(spreadsheetID, "A1", &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

}

func main() {

	spreadsheetID = os.Getenv("SPREADSHEET_ID")

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp := srv.Spreadsheets.Values.Get(spreadsheetID, "A1:A5")
	log.Info(resp)

	log.Fatal(http.ListenAndServe(":8080", &SpreadsheetDumper{srv}))

}
