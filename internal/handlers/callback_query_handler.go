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

	params := models.CallbackFromString(event.GetCallbackQuery().Data)

	slog.Debug("handling callback query, command: " + params.Command + " entity: " + params.Entity)

	command := params.Command

	if command == "list" {
		ListItemCallbackQueryHandler(event)
	}
	if command == "info_event" || command == "next_delta" {
		EventInfoCallbackQueryHandler(event)
	}
	if command == "delete" {
		DeleteItemCallbackQueryHandler(event)
	}
	return nil
}
