package reader

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	budget "github.com/jbleduigou/budget2sheets"
	"go.uber.org/zap"
)

// Reader is a interface for reading a message in SQS and converting in to Transaction
type Reader interface {
	Read(m events.SQSMessage) (budget.Transaction, error)
}

// NewReader returns an instances of a reader, actual implementation is not exposed
func NewReader() Reader {
	return &sqsReader{}
}

type sqsReader struct {
}

func (r *sqsReader) Read(m events.SQSMessage) (budget.Transaction, error) {
	zap.L().Info("Processing SQS message",
		zap.String("message-id", m.MessageId))
	var t budget.Transaction
	err := json.Unmarshal([]byte(m.Body), &t)
	zap.L().Info("Successfully processed SQS message",
		zap.String("message-id", m.MessageId))
	return t, err
}
