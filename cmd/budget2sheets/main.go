package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/jbleduigou/budget2sheets/authentication"
	"github.com/jbleduigou/budget2sheets/config"
	"github.com/jbleduigou/budget2sheets/reader"
	"github.com/jbleduigou/budget2sheets/writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	initLogger(ctx)

	zap.S().Info("Reading configuration")
	config := config.NewConfiguration()
	zap.S().Info("Creating reader")
	reader := reader.NewReader()
	zap.S().Info("Getting Google Sheets client")
	client, err := getClient()

	srv, err := sheets.New(client)
	if err != nil {
		zap.S().Error("Unable to retrieve Google Sheets client", zap.Error(err))
		os.Exit(1)
	}
	zap.S().Info("Created Google Sheets client with success")

	writer := writer.NewWriter(srv, config.GetSpreadSheetID(), config.GetWriteRange())

	cmd := command{r: reader, w: writer}

	for _, record := range event.Records {
		err = cmd.process(record)
		if err != nil {
			zap.S().Error("Error while processing SQS message", zap.Error(err))
			os.Exit(1)
		}
	}
}

func initLogger(ctx context.Context) {
	// Retrieve AWS Request ID
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := lc.AwsRequestID
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"request-id": requestID},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()
	defer zap.S().Sync()
	zap.ReplaceGlobals(logger)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
