# Populate your Google Spreadsheet with data from Dynamo DB

1. Get your Google Service Account (SA) Credentials 
Open the (Service accounts page)[https://console.developers.google.com/iam-admin/serviceaccounts]If prompted, select a project. 

Copy email address of your SA. In `Actions` tab, select `create key` and choose `json` as your key type.

2. Enable Google Sheets API for your Service Account

Navigate [API Library](https://console.developers.google.com/apis/library/sheets.googleapis.com) and click on `enable` button on `Google Sheets API`

3. Create Google Sheet and Give your Service Account Email ability to edit its content

4. Copy your Google Sheet ID 
`https://docs.google.com/spreadsheets/d/[your_sheet_id_is_here]/edit#gid=0`

5. Pass Credentials & Google Sheet ID as enviromental variables to your functon

    - create a valid secret with ``` kubectl create secret generic gsheets --from-literal=SPREADSHEET_ID=yourSpreadSheetIDhere --from-file=CREDENTIALS=credentials.json ```
    - deploy the function with ``` tm deploy service googlesheets -f . --build-template https://raw.githubusercontent.com/triggermesh/knative-lambda-runtime/master/go-1.x/buildtemplate.yaml --env-secret gsheets --wait ```

6. Wait till your function is ready and then run DynamoDB source to send events for your function 


Enjoy Data Written to Your Google Sheet!