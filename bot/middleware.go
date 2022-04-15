package bot

import (
	"github.com/husanmusa/code-learn-bot/service/user"
	tele "gopkg.in/telebot.v3"
)

func auth(userService user.Service) tele.MiddlewareFunc {
	return func(handler tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			user, err := userService.Auth(ctx.Chat().ID)
			if err != nil {
				return err
			}

			ctx.Set(userCtxKey, user)
			return handler(ctx)
		}
	}
}
