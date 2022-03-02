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
	m := events.SQSMessage{
		MessageId: "ID",
		Body: `
		{
			"Date": "01/01/2020", 
			"Description": "<description/>",
			"Comment": "<comment/>",
			"Category": "<category/>",
			"Value": 13.37
		}
		`,
	}

	return m
}

func TestProcessWithSuccess(t *testing.T) {
	w := mock.NewWriter()
	w.On("Write", budget.NewTransaction("01/01/2020", "<description/>", "<comment/>", "<category/>", 13.37), "ID").Return(nil)

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
	w.On("Write", budget.NewTransaction("01/01/2020", "<description/>", "<comment/>", "<category/>", 13.37), "ID").Return(fmt.Errorf("error for unit test"))

	cmd := command{r: reader.NewReader(), w: w}

	err := cmd.process(getMessage())
	assert.Equal(t, "error for unit test", err.Error())
	w.AssertExpectations(t)
}
