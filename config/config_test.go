package config

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go/middleware"
	"github.com/jbleduigou/budget2sheets/mock"
	"github.com/stretchr/testify/assert"
)

func Test_retrieveCredentialsFromSecret(t *testing.T) {
	ctx := context.Background()
	os.Setenv("GOOGLE_CREDENTIALS_SECRET_ARN", "arn:aws:secretsmanager:us-east-1:123456789012:secret:budget2sheets/test")
	tests := []struct {
		name       string
		setupMocks func(m *mock.MockedSecretsManager)
		want       []byte
		wantErr    string
	}{
		{
			name: "Should return error if error when calling GetSecretValue",
			setupMocks: func(m *mock.MockedSecretsManager) {
				m.On("GetSecretValue", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("error for unit tests"))
			},
			want:    nil,
			wantErr: "error for unit tests",
		},
		{
			name: "Should return the value of the secret",
			setupMocks: func(m *mock.MockedSecretsManager) {
				m.On("GetSecretValue", ctx, &secretsmanager.GetSecretValueInput{SecretId: aws.String("arn:aws:secretsmanager:us-east-1:123456789012:secret:budget2sheets/test")}, mock.Anything).
					Return(&secretsmanager.GetSecretValueOutput{
						SecretString:   aws.String("test"),
						ResultMetadata: middleware.Metadata{},
					}, nil)
			},
			want: []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mock.MockedSecretsManager{}
			tt.setupMocks(svc)
			got, err := retrieveCredentialsFromSecret(ctx, svc)
			assert.Equal(t, tt.want, got)
			if err != nil {
				if tt.wantErr == "" {
					assert.Fail(t, "Unexpected error", err.Error())
				}
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
