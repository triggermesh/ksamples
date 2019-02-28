# Populate your Google Spreadsheet with data from Dynamo DB

1. Get your Google Service Account (SA) Credentials 

Open the (Service accounts page)[https://console.developers.google.com/iam-admin/serviceaccounts]. If prompted, select a project. 

* Copy email address of your SA.
* In `Actions` tab, select `create key` and choose `json` as your key type, this will download the json key of your service account.

2. Enable Google Sheets API in your Project

Navigate [API Library](https://console.developers.google.com/apis/library/sheets.googleapis.com) and click on `enable` button on `Google Sheets API`

3. Create Google Sheet and Give your Service Account Email ability to edit its content

4. Copy the Google Sheet ID 

`https://docs.google.com/spreadsheets/d/[your_sheet_id_is_here]/edit#gid=0`

5. Pass Credentials & Google Sheet ID as enviromental variables to your functon

* Create a kubernetes secret with:

```
kubectl create secret generic gsheets --from-literal=SPREADSHEET_ID=yourSpreadSheetIDhere \
                                      --from-file=CREDENTIALS=credentials.json
```

Where `credential.json` is the name of the file containing the JSON key of the service account that has write privileges to your Google Sheet.

* Deploy the function with:

```
tm deploy service googlesheets -f . --build-template https://raw.githubusercontent.com/triggermesh/knative-lambda-runtime/master/go-1.x/buildtemplate.yaml \
                                    --env-secret gsheets \
                                    --wait
```

6. Wait till your function is ready and then run DynamoDB source to send events for your function 


Enjoy Data Written to Your Google Sheet!

## Use case

Write a playlist in DynamoDB and watch the entries being populated in your spreadsheet.

1. Create a DynamodDB table called `playlist`

```
aws dynamodb create-table --table-name playlist --attribute-definitions AttributeName=Artist,AttributeType=S AttributeName=SongTitle,AttributeType=S \
                                                --key-schema AttributeName=Artist,KeyType=HASH AttributeName=SongTitle,KeyType=RANGE \
                                                --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
                                                --stream-specification StreamEnabled=true,StreamViewType=NEW_IMAGE
```

2. Store Songs in your table

See the `songs` directory for example JSON object representing Song items. Store them in DynamoDB using the AWS CLI.

```
 aws dynamodb put-item --table-name playlist --item file://songs/shallow.json
```
