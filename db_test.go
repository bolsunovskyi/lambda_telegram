package tglambda

import (
	"errors"
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

type testDB struct {
	scanErr, putErr, queryErr error
	scanOutput                dynamodb.ScanOutput
	queryOutput               dynamodb.QueryOutput
	putItemOutput             dynamodb.PutItemOutput
}

func (t testDB) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return &t.scanOutput, t.scanErr
}

func (t testDB) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return &t.queryOutput, t.queryErr
}

func (t testDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &t.putItemOutput, t.putErr
}

func TestDB_SaveChat(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}
	lmb.db = testDB{queryOutput: dynamodb.QueryOutput{Count: aws.Int64(0)}, putErr: errors.New("db fatal")}

	if err := lmb.saveChat(125, 125, "test"); err == nil || err.Error() != "db fatal" {
		t.Error("no error on db fatal")
		if err != nil {
			t.Error(err)
		}
		return
	}
}

func TestDB_GetChatsByUsersError(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}
	lmb.db = testDB{scanErr: errors.New("fatal")}

	_, err = lmb.getChatsByUsernames([]string{"bolsunovskyi", "Sundarina"})
	if err == nil {
		t.Error("no error on db fatal")
		return
	}

	lmb.db = testDB{scanErr: nil, scanOutput: dynamodb.ScanOutput{Count: aws.Int64(0)}}
	_, err = lmb.getChatsByUsernames([]string{"bolsunovskyi", "Sundarina"})
	if err == nil {
		t.Error("no error on empty usernames")
		return
	}

	lmb.db = testDB{scanErr: nil, scanOutput: dynamodb.ScanOutput{Count: aws.Int64(1), Items: []map[string]*dynamodb.AttributeValue{
		{
			"chat_id": &dynamodb.AttributeValue{
				N: aws.String("gopa"),
			},
		},
	}}}
	_, err = lmb.getChatsByUsernames([]string{"bolsunovskyi", "Sundarina"})
	if err == nil {
		t.Error("no error on empty usernames")
		return
	}
}

func TestDB_GetChatsByUsers(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	chatID, err := lmb.getChatsByUsernames([]string{"bolsunovskyi", "Sundarina"})
	if err != nil {
		t.Error(err)
		return
	}

	if len(chatID) != 2 {
		t.Error("wrong chat id")
		return
	}
}

func TestDB_GetChatByUserFail(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	lmb.db = testDB{scanOutput: dynamodb.ScanOutput{Count: aws.Int64(0)}}

	if _, err := lmb.getChatByUsername("bolsunovskyi"); err == nil || err.Error() != "username not found" {
		t.Error("no error on scan fatal")
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestDB_GetChatByUser(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	chatID, err := lmb.getChatByUsername("bolsunovskyi")
	if err != nil {
		t.Error(err)
		return
	}

	if chatID == 0 {
		t.Error("wrong chat id")
		return
	}
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

	if err := lmb.saveChat(update.Message.Chat.ID, update.Message.From.ID, update.Message.From.Username); err != nil {
		t.Error(err)
		return
	}

	if err := lmb.saveChat(update.Message.Chat.ID, update.Message.From.ID, update.Message.From.Username); err != nil {
		t.Error(err)
		return
	}
}
