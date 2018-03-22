package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bolsunovskyi/lambda_telegram"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	lmb, err := tglambda.Make(sess, &http.Client{Timeout: time.Second * 5})
	if err != nil {
		log.Fatalln(err)
	}

	lambda.Start(lmb.SNSHandler)
}
