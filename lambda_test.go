package tglambda

import (
	"log"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}
}

func TestLambda_Handler(t *testing.T) {
	sess := makeTestSession()

	_, err := Make(sess, http.DefaultClient)
	if err != nil {
		t.Error(err)
		return
	}

}
