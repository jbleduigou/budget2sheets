package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jbleduigou/budget2sheets/authentication"
	"github.com/jbleduigou/budget2sheets/config"
	"github.com/jbleduigou/budget2sheets/reader"
	"github.com/jbleduigou/budget2sheets/writer"
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
		log.Fatalf("Unable to retrieve Google Sheets client: %v", err)
	}

	writer := writer.NewWriter(srv, config.GetSpreadSheetID(), config.GetWriteRange())

	for _, record := range event.Records {
		// Read Transaction from SQS event
		t, err := reader.Read(record)
		// Write Transaction to Google Sheets
		err = writer.Write(t)
		if err != nil {
			log.Fatalf("Unable to send data to Google Sheets. %v", err)
		}
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
