package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	// Create all collaborators for command
	sess := session.Must(session.NewSession())
	fmt.Println(sess)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	// lambda.Start(handler)

	fmt.Println("Hello gophers!")
}
