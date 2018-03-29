package tglambda

import (
	"encoding/json"
	"strings"
)

func (l Lambda) SNSHandler(rq interface{}) (interface{}, error) {
	usernames := strings.Split(l.params[snsTopicUsernames], ",")
	body, _ := json.Marshal(rq)

	chats, err := l.getChatsByUsernames(usernames)
	if err != nil {
		return nil, err
	}

	for _, chat := range chats {
		if err := l.tgClient.SendMessage(chat, string(body)); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
