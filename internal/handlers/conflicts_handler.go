package handlers

import (
	"context"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
)

func CheckConflictsCallbackHandler(event common.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	event.ReplyCallbackQuery(
		ctx,
		buildConflictsMessage(ctx, params.Id, getConflicts(ctx, params.Id)),
	)

	return nil
}
