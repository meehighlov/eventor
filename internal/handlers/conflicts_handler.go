package handlers

import (
	"context"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/models"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func CheckConflictsCallbackHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	params := models.CallbackFromString(event.GetCallbackQuery().Data)

	event.ReplyCallbackQuery(
		ctx,
		buildConflictsMessage(ctx, params.Id, getConflicts(ctx, params.Id)),
	)

	return nil
}
