package tglambda

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bolsunovskyi/lambda_telegram/df"
	"github.com/bolsunovskyi/lambda_telegram/params"
	"github.com/bolsunovskyi/lambda_telegram/tg"
)

const (
	allowedUsernamesParam = "telegram_allowed_usernames"
	telegramTokenParam    = "telegram_bot_token"
	dialogFlowTokenParam  = "dialogflow_token"
	dialogFlowLangParam   = "dialogflow_lang"
	paramRefreshTime      = 300
	postbackPasswordParam = "telegram_postback_password"
	snsTopicUsernames     = "telegram_sns_usernames"
)

type Lambda struct {
	sess             *session.Session
	tgClient         TgClient
	dfClient         DFClient
	paramsClient     ParamsClient
	httpClient       *http.Client
	allowedUsernames []string
	params           map[string]string
}

type TgClient interface {
	SendMessage(chatID int, text string) error
}

type DFClient interface {
	SendMessage(sessionID string, query string) (*df.Response, error)
}

type ParamsClient interface {
	GetParams(names []string) (map[string]string, error)
}

func (l *Lambda) loadConfig() error {
	pms, err := l.paramsClient.GetParams([]string{
		allowedUsernamesParam,
		telegramTokenParam,
		dialogFlowTokenParam,
		dialogFlowLangParam,
		postbackPasswordParam,
		snsTopicUsernames,
	})
	if err != nil {
		return err
	}
	l.params = pms

	l.allowedUsernames = strings.Split(pms[allowedUsernamesParam], ",")
	l.tgClient = tg.MakeClient(pms[telegramTokenParam], l.httpClient)
	l.dfClient = df.Make(pms[dialogFlowTokenParam], pms[dialogFlowLangParam], l.httpClient)

	return nil
}

func Make(sess *session.Session, httpClient *http.Client) (*Lambda, error) {
	lambda := Lambda{
		sess:         sess,
		paramsClient: params.Make(sess, paramRefreshTime),
		httpClient:   httpClient,
	}

	return &lambda, nil
}

func (l *Lambda) WebHookHandler(update tg.Update) (interface{}, error) {
	if err := l.loadConfig(); err != nil {
		return nil, err
	}

	if err := l.checkUsername(&update); err != nil {
		return nil, err
	}

	if update.Message.Text == "/start" {
		return nil, l.saveChat(&update)
	} else if update.Message.Text == "ping" {
		return nil, l.tgClient.SendMessage(update.Message.Chat.ID, "pong")
	} else if update.Message.Text != "" {
		rsp, err := l.dfClient.SendMessage(strconv.Itoa(update.Message.Chat.ID), update.Message.Text)
		if err != nil {
			return nil, err
		}

		if err := l.tgClient.SendMessage(update.Message.Chat.ID, rsp.Result.Speech); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (l Lambda) checkUsername(update *tg.Update) error {
	for _, u := range l.allowedUsernames {
		if update.Message.From.Username == u {
			return nil
		}
	}

	return fmt.Errorf("username [%s] is not in list", update.Message.From.Username)
}
