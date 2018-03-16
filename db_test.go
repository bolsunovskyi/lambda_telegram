package tglambda

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/bolsunovskyi/lambda_telegram/tg"
)

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

func TestDB_SaveUser(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	update := tg.Update{
		Message: tg.Message{
			From: tg.User{
				Username: "mike_t",
				ID:       66613,
			},
			Chat: tg.Chat{
				ID: 66613999,
			},
		},
	}

	svc := dynamodb.New(sess)
	if _, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"chat_id": {
				N: aws.String(strconv.Itoa(update.Message.Chat.ID)),
			},
		},
		TableName: aws.String("chat"),
	}); err != nil {
		t.Error(err)
		return
	}

	if err := lmb.saveChat(&update); err != nil {
		t.Error(err)
		return
	}

	if err := lmb.saveChat(&update); err != nil {
		t.Error(err)
		return
	}
}
