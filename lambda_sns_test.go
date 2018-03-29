package tglambda

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestLambda_SNSHandler(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}
	lmb.tgClient = testTg{}

	if _, err := lmb.SNSHandler(map[string]string{"foo": "bar"}); err != nil {
		t.Error(err)
		return
	}
}

func TestLambda_SNSHandler_Failure(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	lmb.db = testDB{scanErr: errors.New("db fatal")}

	if _, err := lmb.SNSHandler(map[string]string{"foo": "bar"}); err == nil || err.Error() != "db fatal" {
		t.Error("no error on db fatal")
		if err != nil {
			t.Error(err)
		}

		return
	}

	lmb.db = testDB{scanErr: nil, scanOutput: dynamodb.ScanOutput{Count: aws.Int64(1), Items: []map[string]*dynamodb.AttributeValue{
		{
			"chat_id": &dynamodb.AttributeValue{N: aws.String("125")},
		},
	}}}
	lmb.tgClient = testTg{err: errors.New("tg fatal")}

	if _, err := lmb.SNSHandler(map[string]string{"foo": "bar"}); err == nil || err.Error() != "tg fatal" {
		t.Error("no error on tg fatal")
		if err != nil {
			t.Error(err)
		}

		return
	}
}
