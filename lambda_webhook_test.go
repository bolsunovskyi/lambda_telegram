package tglambda

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/bolsunovskyi/lambda_telegram/df"
	"github.com/bolsunovskyi/lambda_telegram/tg"
)

type testTg struct {
	err error
}

func (t testTg) SendMessage(chatID int, text string) error {
	return t.err
}

func TestLambda_WebHookHandlerText(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	lmb.tgClient = testTg{}

	if _, err := lmb.WebHookHandler(tg.Update{
		Message: tg.Message{
			From: tg.User{
				Username: "bolsunovskyi",
				ID:       125,
			},
			Text: "ping",
			Chat: tg.Chat{
				ID: 125,
			},
		},
	}); err != nil {
		t.Error(err)
		return
	}

	if _, err := lmb.WebHookHandler(tg.Update{
		Message: tg.Message{
			From: tg.User{
				Username: "bolsunovskyi",
				ID:       125,
			},
			Text: "привет как дела ?",
			Chat: tg.Chat{
				ID: 125,
			},
		},
	}); err != nil {
		t.Error(err)
		return
	}
}

type testDF struct {
	err error
	rsp df.Response
}

func (t testDF) SendMessage(sessionID string, query string) (*df.Response, error) {
	return &t.rsp, t.err
}

func TestLambda_WebHookHandlerTextFailure(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	lmb.tgClient = testTg{}
	lmb.dfClient = testDF{err: errors.New("df fatal")}

	if _, err := lmb.WebHookHandler(tg.Update{
		Message: tg.Message{
			From: tg.User{
				Username: "bolsunovskyi",
				ID:       125,
			},
			Text: "привет как дела ?",
			Chat: tg.Chat{
				ID: 125,
			},
		},
	}); err == nil || err.Error() != "df fatal" {
		t.Error("no error on df fatal")
		if err != nil {
			t.Error(err)
		}
		return
	}

	lmb.tgClient = testTg{err: errors.New("tg fatal")}
	lmb.dfClient = testDF{rsp: df.Response{Result: df.Result{Speech: "hello"}}}

	if _, err := lmb.WebHookHandler(tg.Update{
		Message: tg.Message{
			From: tg.User{
				Username: "bolsunovskyi",
				ID:       125,
			},
			Text: "привет как дела ?",
			Chat: tg.Chat{
				ID: 125,
			},
		},
	}); err == nil || err.Error() != "tg fatal" {
		t.Error("no error on df fatal")
		if err != nil {
			t.Error(err)
		}
		return
	}
}

func TestLambda_WebHookHandlerFailure(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := lmb.WebHookHandler(tg.Update{}); err == nil {
		t.Error("no error on empty webhook")
		return
	}

	if _, err := lmb.WebHookHandler(tg.Update{
		Message: tg.Message{
			From: tg.User{
				Username: "mike_lol",
				ID:       125,
			},
			Text: "/start",
			Chat: tg.Chat{
				ID: 125,
			},
		},
	}); err == nil || err.Error() != "username [mike_lol] is not in list" {
		t.Error("no error on wrong username")
		if err != nil {
			t.Error(err)
		}
		return
	}

	lmb.db = testDB{queryErr: errors.New("db fatal")}

	if _, err := lmb.WebHookHandler(tg.Update{
		Message: tg.Message{
			From: tg.User{
				Username: "bolsunovskyi",
				ID:       125,
			},
			Text: "/start",
			Chat: tg.Chat{
				ID: 125,
			},
		},
	}); err == nil || err.Error() != "db fatal" {
		t.Error("no error on db fatal")
		if err != nil {
			t.Error(err)
		}
		return
	}

	lmb, _ = Make(sess, http.DefaultClient)
	lmb.tgClient = testTg{err: errors.New("tg fatal")}
	lmb.db = testDB{queryOutput: dynamodb.QueryOutput{Count: aws.Int64(1)}}

	if _, err := lmb.WebHookHandler(tg.Update{
		Message: tg.Message{
			From: tg.User{
				Username: "bolsunovskyi",
				ID:       125,
			},
			Text: "/start",
			Chat: tg.Chat{
				ID: 125,
			},
		},
	}); err == nil || err.Error() != "tg fatal" {
		t.Error("no error on tg fatal")
		if err != nil {
			t.Error(err)
		}
		return
	}
}
