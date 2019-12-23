package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jbleduigou/budget2sheets/authentication"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func getClient() (*http.Client, error) {
	config, err := google.ConfigFromJSON([]byte(authentication.GetCredentials()), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), authentication.GetToken()), nil
}

func handler(ctx context.Context, s3Event events.S3Event) {
	client, err := getClient()

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")
	writeRange := "Suivi Dépenses Janvier!A2"

	for _, record := range s3Event.Records {
		vr, _ := readData(record.S3.Bucket.Name, record.S3.Object.Key)

		_, err = srv.Spreadsheets.Values.Append(spreadsheetID, writeRange, &vr).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet. %v", err)
		}
	}
}

func readData(bucketName string, objectKey string) (sheets.ValueRange, error) {
	var vr sheets.ValueRange

	content, _ := downloadFile(bucketName, objectKey)
	reader := csv.NewReader(bytes.NewReader(content))

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		return vr, err
	}

	for _, each := range rawCSVdata {
		if len(each) == 5 {
			euro, _ := strconv.ParseFloat(string(each[4]), 64)
			myval := []interface{}{each[0], each[1], each[2], each[3], euro}
			vr.Values = append(vr.Values, myval)
		}
	}
	return vr, nil
}

func downloadFile(bucketName string, objectKey string) ([]byte, error) {
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	fmt.Printf("Downloading file '%v' from bucket '%v' \n", objectKey, bucketName)
	buff := &aws.WriteAtBuffer{}
	n, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		fmt.Printf("Failed to download file %v\n, %v", objectKey, err)
		return nil, err
	}
	fmt.Printf("File %v downloaded, read %d bytes\n", objectKey, n)
	return buff.Bytes(), nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
