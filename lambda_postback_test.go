package tglambda

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

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

	if _, err := lmb.PostBackHandler(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"pwd":      pwd,
			"username": "bolsunovskyi",
		},
		Body: "test",
	}); err != nil {
		t.Error(err)
		return
	}
}
