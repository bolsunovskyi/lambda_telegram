package tg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const baseURL = "https://api.telegram.org/bot"

type Client struct {
	httpClient *http.Client
	token      string
}

func MakeClient(token string, client *http.Client) Client {
	return Client{
		httpClient: client,
		token:      token,
	}
}

func (c Client) SendMessage(chatID int, text string) error {
	sendRequest := url.Values{
		"text":    {text},
		"chat_id": {strconv.Itoa(chatID)},
	}

	rq, err := http.NewRequest("POST", fmt.Sprintf("%s%s/%s", baseURL, c.token, "sendMessage"),
		strings.NewReader(sendRequest.Encode()))
	rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	rsp, err := c.httpClient.Do(rq)
	if err != nil {
		return err
	}

	if rsp.StatusCode != http.StatusOK {
		bts, _ := ioutil.ReadAll(rsp.Body)
		return errors.New(string(bts))
	}

	return nil
}
