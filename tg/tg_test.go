package tg

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gopkg.in/h2non/gock.v1"
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

func TestClient_SendMessageFail(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.telegram.org").
		Post("/bot/" + os.Getenv("tg_token") + "/sendMessage").
		Reply(400).
		JSON(map[string]string{"foo": "bar"})

	cl := MakeClient(os.Getenv("tg_token"), http.DefaultClient)
	if err := cl.SendMessage(148901293, "hello world"); err == nil {
		t.Error("no error on wrong response from telegram")
		return
	}
}

func TestClient_SendMessageFail2(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.telegram.org").
		Post("/bot" + os.Getenv("tg_token") + "/sendMessage").
		Reply(400).
		JSON(map[string]string{"foo": "bar"})

	cl := MakeClient(os.Getenv("tg_token"), http.DefaultClient)
	if err := cl.SendMessage(148901293, "hello world"); err == nil {
		t.Error("no error on wrong response from telegram")
		return
	}
}
