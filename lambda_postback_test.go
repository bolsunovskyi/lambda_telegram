package tglambda

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestLambda_PostBackHandler_Failure2(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	pwd, ok := lmb.params[postbackPasswordParam]
	if !ok {
		t.Error("wrong postback password")
		return
	}

	lmb.db = testDB{scanErr: errors.New("db fatal")}

	if rsp, _ := lmb.PostBackHandler(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"pwd":      pwd,
			"username": "bolsunovskyi",
		},
		Body: "test",
	}); rsp.StatusCode != 400 || rsp.Body != "db fatal" {
		t.Error("no error on wrong poassword param")
		t.Log(rsp.Body)
		return
	}

}

func TestLambda_PostBackHandler_Failure(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	if rsp, _ := lmb.PostBackHandler(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"pwd": lmb.params[postbackPasswordParam],
		},
		Body: "test",
	}); rsp.StatusCode != 400 || rsp.Body != "wrong username" {
		t.Error("no error on wrong poassword param")
		t.Log(rsp.Body)
		return
	}

	if rsp, _ := lmb.PostBackHandler(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"pwd":      "123",
			"username": "bolsunovskyi",
		},
		Body: "test",
	}); rsp.StatusCode != 400 || rsp.Body != "wrong password" {
		t.Error("no error on wrong poassword param")
		t.Log(rsp.Body)
		return
	}

	delete(lmb.params, postbackPasswordParam)

	if rsp, _ := lmb.PostBackHandler(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"pwd":      "123",
			"username": "bolsunovskyi",
		},
		Body: "test",
	}); rsp.StatusCode != 400 || rsp.Body != "password param not found" {
		t.Error("no error on wrong poassword param")
		t.Log(rsp.Body)
		return
	}
}

func TestLambda_PostBackHandler(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	if err := lmb.loadConfig(); err != nil {
		t.Error(err)
		return
	}

	pwd, ok := lmb.params[postbackPasswordParam]
	if !ok {
		t.Error("wrong postback password")
		return
	}
	lmb.tgClient = testTg{}

	if rsp, _ := lmb.PostBackHandler(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"pwd":      pwd,
			"username": "bolsunovskyi",
		},
		Body: "test",
	}); rsp.StatusCode != 200 {
		t.Error(rsp.Body)
		return
	}
}
