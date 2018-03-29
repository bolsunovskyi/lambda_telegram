package tglambda

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (l Lambda) getChatsByUsernames(usernames []string) ([]int, error) {
	var rsp []int

	exps := make(map[string]*dynamodb.AttributeValue)
	var placeholders []string
	for k, v := range usernames {
		placeholder := ":u" + strconv.Itoa(k)
		exps[placeholder] = &dynamodb.AttributeValue{S: aws.String(v)}
		placeholders = append(placeholders, placeholder)
	}

	res, err := l.db.Scan(&dynamodb.ScanInput{
		TableName:                 aws.String("chat"),
		ExpressionAttributeValues: exps,
		FilterExpression:          aws.String(fmt.Sprintf("username IN (%s)", strings.Join(placeholders, ","))),
		ProjectionExpression:      aws.String("chat_id, username, user_id"),
	})

	if err != nil {
		return nil, err
	}

	if *res.Count == 0 {
		return nil, errors.New("usernames not found")
	}

	for _, v := range res.Items {
		chatID, err := strconv.Atoi(*v["chat_id"].N)
		if err != nil {
			return nil, err
		}

		rsp = append(rsp, chatID)
	}

	return rsp, nil
}

func (l Lambda) getChatByUsername(username string) (int, error) {
	res, err := l.db.Scan(&dynamodb.ScanInput{
		TableName: aws.String("chat"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
		FilterExpression:     aws.String("username = :username"),
		ProjectionExpression: aws.String("chat_id, username, user_id"),
	})

	if err != nil {
		return 0, err
	}

	if *res.Count == 0 {
		return 0, errors.New("username not found")
	}

	return strconv.Atoi(*res.Items[0]["chat_id"].N)
}

func (l Lambda) saveChat(chatID, fromID int, username string) error {
	res, err := l.db.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":chat_id": {
				N: aws.String(strconv.Itoa(chatID)),
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
		if _, err := l.db.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String("chat"),
			Item: map[string]*dynamodb.AttributeValue{
				"chat_id": {
					N: aws.String(strconv.Itoa(chatID)),
				},
				"username": {
					S: aws.String(username),
				},
				"user_id": {
					N: aws.String(strconv.Itoa(fromID)),
				},
			},
			ReturnConsumedCapacity: aws.String("NONE"),
		}); err != nil {
			return err
		}
	}

	return nil
}
