package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/jbleduigou/budget2sheets/iface"
	"github.com/stretchr/testify/mock"
)

// MockedSecretsManager is an implementation of the SecretsManager interface with a mock, use for testing not for production
type MockedSecretsManager struct {
	iface.SecretsManager
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
