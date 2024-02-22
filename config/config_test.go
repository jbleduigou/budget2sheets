package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

// MockedSecretsManager is an implementation of the SecretsManager interface with a mock, use for testing not for production
type MockedSecretsManager struct {
	SecretsManager
	mock.Mock
}

func (_m *MockedSecretsManager) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	ret := _m.Called(ctx, params, optFns)
	if ret.Get(0) == nil && ret.Get(1) == nil {
		return nil, nil
	}
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(error)
	}
	return ret.Get(0).(*secretsmanager.GetSecretValueOutput), nil
}

func Test_retrieveCredentialsFromSecret(t *testing.T) {
	ctx := context.Background()
	os.Setenv("GOOGLE_CREDENTIALS_SECRET_ARN", "arn:aws:secretsmanager:us-east-1:123456789012:secret:budget2sheets/test")
	tests := []struct {
		name       string
		setupMocks func(m *MockedSecretsManager)
		want       []byte
		wantErr    string
	}{
		{
			name: "Should return error if error when calling GetSecretValue",
			setupMocks: func(m *MockedSecretsManager) {
				m.On("GetSecretValue", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("error for unit tests"))
			},
			want:    nil,
			wantErr: "error for unit tests",
		},
		{
			name: "Should return the value of the secret",
			setupMocks: func(m *MockedSecretsManager) {
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
			svc := &MockedSecretsManager{}
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
