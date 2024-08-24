package handlers

import (
	"context"
	"strings"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func CommandsHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	commands := []string{
		"Это список моих команд🙌\n",
		"/add - добавить новое событие",
		"/list - все события",
	}

	msg := strings.Join(commands, "\n")

	event.Reply(ctx, msg)

	return nil
}
