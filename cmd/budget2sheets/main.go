package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	budget "github.com/jbleduigou/budget2sheets"
	"github.com/jbleduigou/budget2sheets/authentication"
	"github.com/jbleduigou/budget2sheets/config"
	"github.com/jbleduigou/budget2sheets/reader"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func getClient() (*http.Client, error) {
	oauth2, err := google.ConfigFromJSON([]byte(authentication.GetCredentials()), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	return oauth2.Client(context.Background(), authentication.GetToken()), nil
}

func handler(ctx context.Context, event events.SQSEvent) {
	config := config.NewConfiguration()
	reader := reader.NewReader()
	client, err := getClient()

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	for _, record := range event.Records {
		t, _ := reader.Read(record)
		vr, _ := asValueRange(t)
		_, err = srv.Spreadsheets.Values.Append(config.GetSpreadSheetID(), config.GetWriteRange(), &vr).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet. %v", err)
		}
	}
}

func asValueRange(t budget.Transaction) (sheets.ValueRange, error) {
	var vr sheets.ValueRange
	myval := []interface{}{t.Date, t.Description, t.Comment, t.Category, t.Value}
	vr.Values = append(vr.Values, myval)
	return vr, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
