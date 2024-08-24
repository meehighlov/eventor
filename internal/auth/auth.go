package auth

import (
	"fmt"
	"log/slog"

	"github.com/meehighlov/eventor/internal/config"

	"github.com/meehighlov/eventor/pkg/telegram"
)

func isAuth(tgusername string) bool {
	for _, auth_user_name := range config.Cfg().AuthList() {
		if auth_user_name == tgusername {
			return true
		}
	}

	return false
}

func Auth(handler telegram.CommandHandler) telegram.CommandHandler {
	return func(event telegram.Event) error {
		message := event.GetMessage()
		if isAuth(message.From.Username) {
			return handler(event)
		}

		msg := fmt.Sprintf("Unauthorized access attempt by user: id=%d usernmae=%s", message.From.Id, message.From.Username)
		slog.Info(msg)

		return nil
	}
}
