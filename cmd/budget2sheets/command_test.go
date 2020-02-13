package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	budget "github.com/jbleduigou/budget2sheets"
	"github.com/jbleduigou/budget2sheets/mock"
	"github.com/jbleduigou/budget2sheets/reader"
	"github.com/stretchr/testify/assert"
)

func getMessage() events.SQSMessage {
	date := "01/01/2020"
	description := "<description/>"
	comment := "<comment/>"
	category := "<category/>"
	value := "13.37"

	m := events.SQSMessage{
		MessageId:         "ID",
		MessageAttributes: make(map[string]events.SQSMessageAttribute),
	}
	m.MessageAttributes["Date"] = events.SQSMessageAttribute{StringValue: &date}
	m.MessageAttributes["Description"] = events.SQSMessageAttribute{StringValue: &description}
	m.MessageAttributes["Comment"] = events.SQSMessageAttribute{StringValue: &comment}
	m.MessageAttributes["Category"] = events.SQSMessageAttribute{StringValue: &category}
	m.MessageAttributes["Value"] = events.SQSMessageAttribute{StringValue: &value}

	return m
}

func TestProcessWithSuccess(t *testing.T) {
	w := mock.NewWriter()
	w.On("Write", budget.NewTransaction("01/01/2020", "<description/>", "<comment/>", "<category/>", 13.37)).Return(nil)

	cmd := command{r: reader.NewReader(), w: w}

	err := cmd.process(getMessage())
	assert.Nil(t, err)
	w.AssertExpectations(t)
}

func TestProcessWithReaderError(t *testing.T) {
	r := mock.NewReader()
	w := mock.NewWriter()
	r.On("Read", getMessage()).Return(&budget.Transaction{}, fmt.Errorf("error for unit test"))

	cmd := command{r: r, w: w}

	err := cmd.process(getMessage())
	assert.Equal(t, "error for unit test", err.Error())
	r.AssertExpectations(t)
	w.AssertExpectations(t)
}

func TestProcessWithWriterError(t *testing.T) {
	w := mock.NewWriter()
	w.On("Write", budget.NewTransaction("01/01/2020", "<description/>", "<comment/>", "<category/>", 13.37)).Return(fmt.Errorf("error for unit test"))

	cmd := command{r: reader.NewReader(), w: w}

	err := cmd.process(getMessage())
	assert.Equal(t, "error for unit test", err.Error())
	w.AssertExpectations(t)
}
