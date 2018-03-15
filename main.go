package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(update Update) (interface{}, error) {
	log.Printf("%+v", update)
	return nil, nil
}

func main() {
	lambda.Start(Handler)
}
