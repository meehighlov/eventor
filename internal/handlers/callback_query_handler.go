package handlers

import (
	"context"
	"log/slog"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/models"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func CallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	event.AnswerCallbackQuery(ctx)

	command := models.CallbackFromString(event.GetCallbackQuery().Data).Command

	slog.Debug("handling callback query, command: " + command)

	if command == "list" {
		ListItemCallbackQueryHandler(event)
	}
	if command == "info" {
		EventInfoCallbackQueryHandler(event)
	}
	if command == "delete" {
		DeleteItemCallbackQueryHandler(event)
	}
	if command == "conflicts" {
		CheckConflictsCallbackHandler(event)
	}
	return nil
}
