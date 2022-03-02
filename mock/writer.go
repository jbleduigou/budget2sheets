package mock

import (
	budget "github.com/jbleduigou/budget2sheets"
	"github.com/stretchr/testify/mock"
)

// NewWriter provides a mock instance of a Writer
func NewWriter() *Writer {
	return &Writer{}
}

// Writer is an implementation of the Writer interface with a mock, use for testing not for production
type Writer struct {
	mock.Mock
}

func (_m *Writer) Write(_a0 budget.Transaction, _a1 string) error {
	ret := _m.Called(_a0, _a1)
	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(error)
}
