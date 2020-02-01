package reader

import (
	"fmt"
	"strconv"

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
	euro, err := strconv.ParseFloat(r.getAttributeValue(m, "Value"), 64)
	if err != nil {
		fmt.Printf("Error while parsing transaction value: %v\n", err)
	}
	date := r.getAttributeValue(m, "Date")
	description := r.getAttributeValue(m, "Description")
	comment := r.getAttributeValue(m, "Comment")
	category := r.getAttributeValue(m, "Category")
	fmt.Printf("Description of transaction is '%v'\n", description)
	return budget.NewTransaction(date, description, comment, category, euro), nil
}

func (r *sqsReader) getAttributeValue(m events.SQSMessage, name string) string {
	value := m.MessageAttributes[name].StringValue
	if value == nil {
		return ""
	}
	return *value
}
