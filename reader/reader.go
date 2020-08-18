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
func NewReader(logger *zap.SugaredLogger) Reader {
	return &sqsReader{logger: logger}
}

type sqsReader struct {
	logger *zap.SugaredLogger
}

func (r *sqsReader) Read(m events.SQSMessage) (budget.Transaction, error) {
	r.logger.Infof("Processing SQS message with id '%v'", m.MessageId)
	var t budget.Transaction
	err := json.Unmarshal([]byte(m.Body), &t)
	r.logger.Infof("Processed SQS message with id '%v'", m.MessageId)
	return t, err
}
