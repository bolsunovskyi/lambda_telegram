package tglambda

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (l Lambda) PostBackHandler(rq events.APIGatewayProxyRequest) (rsp interface{}, err error) {
	defer func() {
		if err != nil {
			rsp = events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
			}
		} else {
			rsp = events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
			}
		}
	}()

	if err := l.loadConfig(); err != nil {
		return nil, err
	}

	pwdParam, ok := l.params[postbackPasswordParam]
	if !ok {
		return nil, errors.New("password param not found")
	}

	pwd, ok := rq.QueryStringParameters["pwd"]
	if !ok || pwdParam != pwd {
		return nil, errors.New("wrong password")
	}

	username, ok := rq.QueryStringParameters["username"]
	if !ok {
		return nil, errors.New("wrong username")
	}

	chatID, err := l.getChatByUsername(username)
	if err != nil {
		return nil, err
	}

	return nil, l.tgClient.SendMessage(chatID, rq.Body)
}
