package tglambda

import (
	"net/http"
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
