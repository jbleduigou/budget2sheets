package reader

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	budget "github.com/jbleduigou/budget2sheets"
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
	fmt.Printf("Processing SQS message with id '%v'\n", m.MessageId)
	var t budget.Transaction
	err := json.Unmarshal([]byte(m.Body), &t)
	return t, err
}
