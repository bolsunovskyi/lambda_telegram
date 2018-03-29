package tglambda

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (l Lambda) PostBackHandler(rq events.APIGatewayProxyRequest) (rsp events.APIGatewayProxyResponse, err error) {
	defer func() {
		if err != nil {
			rsp = events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       err.Error(),
			}
			err = nil
		} else {
			rsp = events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
			}
		}
	}()

	pwdParam, ok := l.params[postbackPasswordParam]
	if !ok {
		err = errors.New("password param not found")
		return
	}

	pwd, ok := rq.QueryStringParameters["pwd"]
	if !ok || pwdParam != pwd {
		err = errors.New("wrong password")
		return
	}

	username, ok := rq.QueryStringParameters["username"]
	if !ok {
		err = errors.New("wrong username")
		return
	}

	var chatID int
	chatID, err = l.getChatByUsername(username)
	if err != nil {
		return
	}

	err = l.tgClient.SendMessage(chatID, rq.Body)
	return
}
