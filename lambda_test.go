package tglambda

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}
}

func TestLambda_Handler(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("region")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("aws_access_key_id"),
			os.Getenv("aws_secret_access_key"),
			"",
		),
	})

	if err != nil {
		t.Error(err)
		return
	}

	_, err = Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

}
