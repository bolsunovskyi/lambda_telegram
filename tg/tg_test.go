package tg

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalln(err)
	}
}

func TestClient_SendMessage(t *testing.T) {
	cl := MakeClient(os.Getenv("tg_token"), http.DefaultClient)
	if err := cl.SendMessage(148901293, "hello world"); err != nil {
		t.Error(err)
		return
	}
}
