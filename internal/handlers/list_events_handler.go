package handlers

import (
	"context"
	"log/slog"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/pkg/telegram"
)

const (
	HEADER_MESSAGE_LIST_NOT_EMPTY = "Нажми, чтобы узнать детали✨"
	HEADER_MESSAGE_LIST_IS_EMPTY = "Записей пока нет✨"
)

func ListEventsHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	message := event.GetMessage()
	events, err := (&db.Event{OwnerId: message.From.Id}).Filter(ctx)

	if err != nil {
		slog.Error("Error fetching events" + err.Error())
		return nil
	}

	if len(events) == 0 {
		event.Reply(ctx, HEADER_MESSAGE_LIST_IS_EMPTY)
		return nil
	}

	event.ReplyWithKeyboard(
		ctx,
		HEADER_MESSAGE_LIST_NOT_EMPTY,
		common.BuildItemListMarkup(
			events,
			common.LIST_LIMIT,
			common.LIST_START_OFFSET,
			"<",
		),
	)

	return nil
}
