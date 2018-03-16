package params

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Client struct {
	lastUpdate int64
	ssmClient  *ssm.SSM
	cache      map[string]string
	cacheTTL   int64
}

func Make(sess *session.Session, cacheTTL int64) *Client {
	return &Client{
		cache:     make(map[string]string),
		ssmClient: ssm.New(sess),
		cacheTTL:  cacheTTL,
	}
}

func (c *Client) GetParams(names []string) (map[string]string, error) {
	if time.Now().Unix()-c.lastUpdate < c.cacheTTL {
		return c.cache, nil
	}

	var awsNames []*string
	for _, v := range names {
		awsNames = append(awsNames, aws.String(v))
	}

	params, err := c.ssmClient.GetParameters(&ssm.GetParametersInput{
		Names: awsNames,
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

	c.lastUpdate = time.Now().Unix()
	c.cache = res

	return res, nil
}
