package tglambda

import (
	"net/http"
	"testing"
)

func TestLambda_SNSHandler(t *testing.T) {
	sess := makeTestSession()
	lmb, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := lmb.SNSHandler(map[string]string{"foo": "bar"}); err != nil {
		t.Error(err)
		return
	}
}
