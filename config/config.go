package config

import (
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/jbleduigou/budget2sheets/iface"
	"golang.org/x/net/context"
)

// Configuration provides and interface for the configuration of the software
type Configuration interface {
	GetSpreadSheetID() string
	GetWriteRange() string
	GetGoogleJsonCredentials() []byte
}

// NewConfiguration will provide an instance of a Configuration, implementation is not exposed
func NewConfiguration(ctx context.Context) (Configuration, error) {
	svc, err := initSecretsManagerClient(ctx)
	if err != nil {
		return nil, err
	}

	c, err := retrieveCredentialsFromSecret(ctx, svc)
	if err != nil {
		return nil, err
	}

	return &configuration{
		spreadSheetID:         os.Getenv("GOOGLE_SPREADSHEET_ID"),
		writeRange:            os.Getenv("GOOGLE_SPREADSHEET_RANGE"),
		googleJsonCredentials: c,
	}, nil
}

func initSecretsManagerClient(ctx context.Context) (*secretsmanager.Client, error) {
	slog.Info("Creating Secrets Manager client")
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(awsCfg)
	return svc, nil
}

func retrieveCredentialsFromSecret(ctx context.Context, svc iface.SecretsManager) ([]byte, error) {
	arn := os.Getenv("GOOGLE_CREDENTIALS_SECRET_ARN")
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(arn),
	}
	slog.Info("Retrieving secret",
		slog.String("secret-arn", arn))
	result, err := svc.GetSecretValue(ctx, input)

	if err != nil {
		return nil, err
	}

	return []byte(*result.SecretString), nil
}

type configuration struct {
	spreadSheetID         string
	writeRange            string
	googleJsonCredentials []byte
}

func (c *configuration) GetSpreadSheetID() string {
	return c.spreadSheetID
}

func (c *configuration) GetWriteRange() string {
	return c.writeRange
}

func (c *configuration) GetGoogleJsonCredentials() []byte {
	return c.googleJsonCredentials
}
