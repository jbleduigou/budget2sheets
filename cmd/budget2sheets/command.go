package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jbleduigou/budget2sheets/reader"
	"github.com/jbleduigou/budget2sheets/writer"
)

type command struct {
	r reader.Reader
	w writer.Writer
}

func (c *command) execute(m events.SQSMessage) {
	// Read Transaction from SQS message
	t, err := c.r.Read(m)
	// Write Transaction to Google Sheets
	err = c.w.Write(t)
	if err != nil {
		log.Fatalf("Unable to send data to Google Sheets. %v", err)
	}
}
