package tglambda

import (
	"fmt"
	"strconv"

	"github.com/bolsunovskyi/lambda_telegram/tg"
	"gopkg.in/go-playground/validator.v9"
)

func (l *Lambda) WebHookHandler(update tg.Update) (interface{}, error) {
	if err := validator.New().Struct(update); err != nil {
		return nil, err
	}

	if err := l.checkUsername(update.Message.From.Username); err != nil {
		return nil, err
	}

	if update.Message.Text == "/start" {
		if err := l.saveChat(update.Message.Chat.ID, update.Message.From.ID, update.Message.From.Username); err != nil {
			return nil, err
		}
		return nil, l.tgClient.SendMessage(update.Message.Chat.ID, "welcome ;)")
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

func (l Lambda) checkUsername(username string) error {
	for _, u := range l.allowedUsernames {
		if username == u {
			return nil
		}
	}

	return fmt.Errorf("username [%s] is not in list", username)
}
