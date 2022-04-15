package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tb "gopkg.in/telebot.v3"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:       os.Getenv("TELEBOT_SECRET"),
		Synchronous: true,
	})
	if err != nil {
		panic(err)
	}

	b.Handle(tb.OnText, func(ctx tb.Context) error {
		return ctx.Send("a")
	})

	lambda.Start(func(req events.APIGatewayProxyRequest) (err error) {
		var u tb.Update
		if err = json.Unmarshal([]byte(req.Body), &u); err == nil {
			b.ProcessUpdate(u)
		}
		return
	})
}
