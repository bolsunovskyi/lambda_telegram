package df

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/h2non/gock.v1"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalln(err)
	}
}

func TestClient_SendMessage(t *testing.T) {
	cl := Make(os.Getenv("df_token"), os.Getenv("df_lang"), http.DefaultClient)
	if _, err := cl.SendMessage(strconv.Itoa(int(time.Now().Unix())), "test"); err != nil {
		t.Error(err)
		return
	}
}

func TestClient_SendMessage_Failure(t *testing.T) {
	cl := Make("", os.Getenv("df_lang"), http.DefaultClient)
	if _, err := cl.SendMessage(strconv.Itoa(int(time.Now().Unix())), "test"); err == nil {
		t.Error("no error on wrong token")
		return
	}

	defer gock.Off()

	gock.New("https://api.telegram.org").
		Post("/bot/" + os.Getenv("tg_token") + "/sendMessage").
		Reply(400).
		JSON(map[string]string{"foo": "bar"})

	if _, err := cl.SendMessage(strconv.Itoa(int(time.Now().Unix())), "test"); err == nil {
		t.Error("no error on wrong http")
		return
	}
}

func TestClient_SendMessage_Failure2(t *testing.T) {
	cl := Make("", os.Getenv("df_lang"), http.DefaultClient)

	defer gock.Off()

	gock.New("https://api.dialogflow.com").
		Post("/v1/query").
		Reply(200).
		BodyString("lol")

	if _, err := cl.SendMessage(strconv.Itoa(int(time.Now().Unix())), "test"); err == nil {
		t.Error("no error on wrong body")
		return
	}
}
