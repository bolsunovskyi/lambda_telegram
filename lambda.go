package tglambda

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/bolsunovskyi/lambda_telegram/tg"
)

const (
	allowedUsernamesParam = "telegram_allowed_usernames"
	telegramTokenParam    = "telegram_bot_token"
)

type Lambda struct {
	sess             *session.Session
	tgClient         TgClient
	httpClient       *http.Client
	allowedUsernames []string
}

type TgClient interface {
	SendMessage(chatID int, text string) error
}

func Make(sess *session.Session, httpClient *http.Client) (*Lambda, error) {
	lambda := Lambda{
		sess:       sess,
		httpClient: httpClient,
	}

	ssmClient := ssm.New(sess)

	param, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(allowedUsernamesParam),
	})
	if err != nil {
		return nil, err
	}

	lambda.allowedUsernames = strings.Split(*param.Parameter.Value, ",")

	param, err = ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(telegramTokenParam),
	})
	if err != nil {
		return nil, err
	}

	lambda.tgClient = tg.MakeClient(*param.Parameter.Value, httpClient)

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
