package df

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	protocolVersion = "20150910"
	baseURL         = "https://api.dialogflow.com/v1"
)

type Request struct {
	Query     string `json:"query"`
	SessionID string `json:"sessionId"`
	Lang      string `json:"lang"`
	Version   string `json:"v"`
}

type Response struct {
	ID        string           `json:"id"`
	Timestamp time.Time        `json:"timestamp"`
	Lang      string           `json:"lang"`
	SessionID string           `json:"session_id"`
	Result    Result           `json:"result"`
	Status    dialogFlowStatus `json:"status"`
}

type Result struct {
	Source        string      `json:"source"`
	ResolvedQuery string      `json:"resolvedQuery"`
	Speech        string      `json:"speech"`
	Action        string      `json:"action"`
	Parameters    interface{} `json:"parameters"`
	MetaData      MetaData    `json:"metadata"`
}

type MetaData struct {
	InputContext              []interface{} `json:"inputContexts"`
	OutputContexts            []interface{} `json:"outputContexts"`
	IntentName                string        `json:"intentName"`
	IntentID                  string        `json:"intentId"`
	WebHookUsed               string        `json:"webhookUsed"`
	WebHookForSlotFillingUsed string        `json:"webhookForSlotFillingUsed"`
	Contexts                  []interface{} `json:"contexts"`
}

type dialogFlowStatus struct {
	Code            int    `json:"code"`
	ErrorType       string `json:"errorType"`
	WebHookTimedOut bool   `json:"webhookTimedOut"`
}

type Client struct {
	httpClient *http.Client
	token      string
	lang       string
}

func Make(token, lang string, httpClient *http.Client) Client {
	return Client{
		token:      token,
		lang:       lang,
		httpClient: httpClient,
	}
}

func (c Client) SendMessage(sessionID string, query string) (*Response, error) {
	bts, err := json.Marshal(Request{
		Query:     query,
		Lang:      c.lang,
		SessionID: sessionID,
		Version:   protocolVersion,
	})
	if err != nil {
		return nil, err
	}

	rq, err := http.NewRequest("POST", fmt.Sprintf(`%s/query`, baseURL),
		bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}
	rq.Header.Add("Content-Type", "application/json; charset=UTF-8")
	rq.Header.Add("Accept", "application/json")
	rq.Header.Add("Authorization", "Bearer "+c.token)

	rsp, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		bts, _ := ioutil.ReadAll(rsp.Body)
		return nil, errors.New(string(bts))
	}

	var dfRsp Response
	if err := json.NewDecoder(rsp.Body).Decode(&dfRsp); err != nil {
		return nil, err
	}

	return &dfRsp, nil
}
