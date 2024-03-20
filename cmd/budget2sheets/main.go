package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jbleduigou/budget2sheets/config"
	"github.com/jbleduigou/budget2sheets/reader"
	"github.com/jbleduigou/budget2sheets/writer"
	slogawslambda "github.com/jbleduigou/slog-aws-lambda"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getClient(ctx context.Context, cfg config.Configuration) (*http.Client, error) {
	jsonCfg, err := google.JWTConfigFromJSON(cfg.GetGoogleJsonCredentials(), sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}
	return jsonCfg.Client(ctx), nil
}

func handler(ctx context.Context, event events.SQSEvent) {
	initLogger(ctx)

	slog.Info("Reading configuration")
	config, err := config.NewConfiguration(ctx)
	if err != nil {
		slog.Error("Unable to retrieve configuration", "error", err)
		os.Exit(1)
	}
	slog.Info("Creating reader")
	reader := reader.NewReader()
	slog.Info("Getting Google Sheets client")
	client, err := getClient(ctx, config)
	if err != nil {
		slog.Error("Unable to retrieve Google Sheets client", "error", err)
		os.Exit(1)
	}
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		slog.Error("Unable to retrieve Google Sheets client", "error", err)
		os.Exit(1)
	}
	slog.Info("Created Google Sheets client with success")

	writer := writer.NewWriter(srv, config.GetSpreadSheetID(), config.GetWriteRange())

	cmd := command{r: reader, w: writer}

	for _, record := range event.Records {
		err = cmd.process(record)
		if err != nil {
			slog.Error("Error while processing SQS message", "error", err)
			os.Exit(1)
		}
	}
}

func initLogger(ctx context.Context) {
	slog.SetDefault(slog.New(slogawslambda.NewAWSLambdaHandler(ctx, nil)))
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
