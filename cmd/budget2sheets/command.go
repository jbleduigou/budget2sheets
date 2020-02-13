package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jbleduigou/budget2sheets/reader"
	"github.com/jbleduigou/budget2sheets/writer"
)

type command struct {
	r reader.Reader
	w writer.Writer
}

func (c *command) process(m events.SQSMessage) error {
	// Read Transaction from SQS message
	t, err := c.r.Read(m)
	if err != nil {
		return err
	}
	// Write Transaction to Google Sheets
	err = c.w.Write(t)
	return err
}
