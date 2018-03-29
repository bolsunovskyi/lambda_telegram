package tglambda

import (
	"errors"
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

func TestLambdaMake(t *testing.T) {
	sess := makeTestSession()

	if _, err := Make(sess, http.DefaultClient); err != nil {
		t.Error(err)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("region")),
		Endpoint: aws.String("http://localhost"),
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

	if _, err := Make(sess, http.DefaultClient); err == nil {
		t.Error("no error on wrong aws session")
		return
	}
}

type testParamsClient struct {
	params map[string]string
	err    error
}

func (t testParamsClient) GetParams(names []string) (map[string]string, error) {
	return t.params, t.err
}

func TestMake_loadConfig(t *testing.T) {
	sess := makeTestSession()

	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}
	lmb.paramsClient = testParamsClient{
		err: errors.New("fatal"),
	}

	if err := lmb.loadConfig(); err == nil {
		t.Error("no error on wrong params client")
		return
	}
}
