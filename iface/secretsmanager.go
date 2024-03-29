package iface

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManager interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}
