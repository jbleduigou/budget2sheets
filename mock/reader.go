package mock

import (
	"github.com/aws/aws-lambda-go/events"
	budget "github.com/jbleduigou/budget2sheets"
	"github.com/stretchr/testify/mock"
)

const (
	// Anything is used in Diff and Assert when the argument being tested
	// shouldn't be taken into consideration.
	Anything = "mock.Anything"
)

// NewReader provides a mock instance of a Reader
func NewReader() *Reader {
	return &Reader{}
}

// Reader is an implementation of the Reader interface with a mock, use for testing not for production
type Reader struct {
	mock.Mock
}

func (_m *Reader) Read(_a0 events.SQSMessage) (budget.Transaction, error) {
	ret := _m.Called(_a0)
	if ret.Get(1) == nil {
		return budget.NewTransaction("", "", "", "", 0.0), nil
	}
	return budget.NewTransaction("", "", "", "", 0.0), ret.Get(1).(error)
}
