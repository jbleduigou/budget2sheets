package config

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// Configuration provides and interface for the configuration of the software
type Configuration interface {
	GetSpreadSheetID() string
	GetWriteRange() string
	GetGoogleJsonCredentials() []byte
}

type SecretsManager interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
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

func initSecretsManagerClient(ctx context.Context) (SecretsManager, error) {
	zap.L().Info("Creating Secrets Manager client")
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(awsCfg)
	return svc, nil
}

func retrieveCredentialsFromSecret(ctx context.Context, svc SecretsManager) ([]byte, error) {
	arn := os.Getenv("GOOGLE_CREDENTIALS_SECRET_ARN")
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(arn),
	}
	zap.L().Info("Retrieving secret",
		zap.String("secret-arn", arn))
	result, err := svc.GetSecretValue(ctx, input)

	return []byte(*result.SecretString), err
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