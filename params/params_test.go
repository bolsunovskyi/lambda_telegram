package params

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalln(err)
	}
}

func makeTestSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("region")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("aws_access_key_id"),
			os.Getenv("aws_secret_access_key"),
			"",
		),
	})

	if err != nil {
		log.Fatalln(err)
	}

	return sess
}

func TestClient_GetParams(t *testing.T) {
	sess := makeTestSession()
	cl := Make(sess, 300)
	if _, err := cl.GetParams([]string{"telegram_postback_password"}); err != nil {
		t.Error(err)
		return
	}

	if _, err := cl.GetParams([]string{"telegram_postback_password"}); err != nil {
		t.Error(err)
		return
	}
}

func TestClient_GetParamsFailure(t *testing.T) {
	sess := makeTestSession()

	cl := Make(sess, 300)
	if _, err := cl.GetParams([]string{"foo"}); err == nil || err.Error() != "not enough params" {
		t.Error("no error on wrong param")
		if err != nil {
			t.Error(err)
		}
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost"),
		Region:   aws.String(os.Getenv("region")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("aws_access_key_id"),
			os.Getenv("aws_secret_access_key"),
			"",
		),
	})

	if err != nil {
		log.Fatalln(err)
		return
	}

	cl = Make(sess, 300)
	if _, err := cl.GetParams([]string{"telegram_postback_password"}); err == nil {
		t.Error("no error on wrong endpoint")
		return
	}
}
