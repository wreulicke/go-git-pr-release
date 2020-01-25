package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"gopkg.in/urfave/cli.v1"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type Request struct {
	Body string `json:"body"`
}

func handleInternal(values url.Values) ([]byte, error) {
	texts, ok := values["text"]
	if !ok || len(texts) == 0 {
		texts = []string{""}
	}
	v := texts[0]
	args := strings.Split(v, " ")
	buf := new(bytes.Buffer)
	err := Run(buf, args...)
	if err != nil {
		return json.Marshal(map[string]interface{}{
			"text":          err.Error(),
			"response_type": "in_channel",
		})
	}
	return json.Marshal(map[string]interface{}{
		"text":          buf.String(),
		"response_type": "in_channel",
	})
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, req Request) (Response, error) {
	values, err := url.ParseQuery(req.Body)
	if err != nil {
		log.Println(err)
		return Response{StatusCode: 500}, err
	}
	body, err := handleInternal(values)
	if err != nil {
		log.Println(err)
		return Response{StatusCode: 500}, err
	}
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func init() {
	cli.AppHelpTemplate = "```" + cli.AppHelpTemplate + "```"
	cli.CommandHelpTemplate = "```" + cli.CommandHelpTemplate + "```"
	cli.SubcommandHelpTemplate = "```" + cli.SubcommandHelpTemplate + "```"
}

func main() {
	lambda.Start(Handler)
}

