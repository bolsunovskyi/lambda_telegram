package tglambda

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/bolsunovskyi/lambda_telegram/tg"
)

func (l Lambda) saveChat(update *tg.Update) error {
	svc := dynamodb.New(l.sess)

	res, err := svc.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":chat_id": {
				N: aws.String(strconv.Itoa(update.Message.Chat.ID)),
			},
		},
		KeyConditionExpression: aws.String("chat_id = :chat_id"),
		ProjectionExpression:   aws.String("chat_id, username, user_id"),
		TableName:              aws.String("chat"),
	})
	if err != nil {
		return err
	}

	if *res.Count == 0 {
		if _, err := svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String("chat"),
			Item: map[string]*dynamodb.AttributeValue{
				"chat_id": {
					N: aws.String(strconv.Itoa(update.Message.Chat.ID)),
				},
				"username": {
					S: aws.String(update.Message.From.Username),
				},
				"user_id": {
					N: aws.String(strconv.Itoa(update.Message.From.ID)),
				},
			},
			ReturnConsumedCapacity: aws.String("NONE"),
		}); err != nil {
			return err
		}
	}

	return l.tgClient.SendMessage(update.Message.Chat.ID, "welcome ;)")
}
