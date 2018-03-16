package tglambda

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/bolsunovskyi/lambda_telegram/df"
	"github.com/bolsunovskyi/lambda_telegram/tg"
)

const (
	allowedUsernamesParam = "telegram_allowed_usernames"
	telegramTokenParam    = "telegram_bot_token"
	dialogFlowTokenParam  = "dialogflow_token"
	dialogFlowLangParam   = "dialogflow_lang"
)

type Lambda struct {
	sess             *session.Session
	tgClient         TgClient
	dfClient         DFClient
	httpClient       *http.Client
	allowedUsernames []string
}

type TgClient interface {
	SendMessage(chatID int, text string) error
}

type DFClient interface {
	SendMessage(sessionID string, query string) (*df.Response, error)
}

func (l Lambda) loadParams(names []*string) (map[string]string, error) {
	ssmClient := ssm.New(l.sess)

	params, err := ssmClient.GetParameters(&ssm.GetParametersInput{
		Names: names,
	})
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for _, v := range params.Parameters {
		res[*v.Name] = *v.Value
	}

	if len(res) != len(names) {
		return nil, errors.New("not enough params")
	}

	return res, nil
}

func Make(sess *session.Session, httpClient *http.Client) (*Lambda, error) {
	lambda := Lambda{
		sess:       sess,
		httpClient: httpClient,
	}

	params, err := lambda.loadParams([]*string{
		aws.String(allowedUsernamesParam),
		aws.String(telegramTokenParam),
		aws.String(dialogFlowTokenParam),
		aws.String(dialogFlowLangParam),
	})
	if err != nil {
		return nil, err
	}

	lambda.allowedUsernames = strings.Split(params[allowedUsernamesParam], ",")
	lambda.tgClient = tg.MakeClient(params[telegramTokenParam], httpClient)
	lambda.dfClient = df.Make(params[dialogFlowTokenParam], params[dialogFlowLangParam], httpClient)

	return &lambda, nil
}

func (l Lambda) Handler(update tg.Update) (interface{}, error) {
	log.Printf("%+v\n", update)
	if err := l.checkUsername(&update); err != nil {
		return nil, err
	}

	if update.Message.Text == "ping" {
		if err := l.tgClient.SendMessage(update.Message.Chat.ID, "pong"); err != nil {
			return nil, err
		}
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

	return errors.New("username is not in list")
}
