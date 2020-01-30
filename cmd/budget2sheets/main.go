package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jbleduigou/budget2sheets/authentication"
	"github.com/jbleduigou/budget2sheets/config"
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
	client, err := getClient()

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	for _, record := range event.Records {
		vr, _ := extractData(record)
		_, err = srv.Spreadsheets.Values.Append(config.GetSpreadSheetID(), config.GetWriteRange(), &vr).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet. %v", err)
		}
	}
}

func extractData(m events.SQSMessage) (sheets.ValueRange, error) {
	fmt.Printf("Processing SQS message with id '%v'\n", m.MessageId)
	var vr sheets.ValueRange
	euro, _ := strconv.ParseFloat(*m.MessageAttributes["Value"].StringValue, 64)
	description := *m.MessageAttributes["Description"].StringValue
	myval := []interface{}{*m.MessageAttributes["Date"].StringValue, description, "", *m.MessageAttributes["Category"].StringValue, euro}
	vr.Values = append(vr.Values, myval)
	fmt.Printf("Description of transaction is '%v'\n", description)
	return vr, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
