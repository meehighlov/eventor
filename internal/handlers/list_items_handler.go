package handlers

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
)

const (
	HEADER_MESSAGE_LIST_NOT_EMPTY = "Нажми, чтобы узнать детали✨"
	HEADER_MESSAGE_LIST_IS_EMPTY = "Записей пока нет✨"
)

func ListEntityHandler(entity string) common.HandlerType {
	return func (event common.Event) error {
		ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
		defer cancel()
	
		message := event.GetMessage()

		items, err := (common.BuildItem(entity, message.From.Id)).Filter(ctx)
	
		if err != nil {
			slog.Error("Error fetching items: " + err.Error())
			return nil
		}
	
		if len(items) == 0 {
			event.Reply(ctx, HEADER_MESSAGE_LIST_IS_EMPTY)
			return nil
		}
	
		event.ReplyWithKeyboard(
			ctx,
			HEADER_MESSAGE_LIST_NOT_EMPTY,
			common.BuildItemListMarkup(
				items,
				common.LIST_LIMIT,
				common.LIST_START_OFFSET,
				"<",
				entity,
			),
		)
	
		return nil
	}
}

// ----------------------------------------------- List items for CallbackQuery ---------------------------------------------------

func ListItemCallbackQueryHandler(event common.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()
	callbackQuery := event.GetCallbackQuery()

	params := common.CallbackFromString(callbackQuery.Data)

	offset := params.Pagination.Offset

	offset_, err := strconv.Atoi(offset)
	if err != nil {
		slog.Error("error parsing offset in list pagination callback query: " + err.Error())
		return err
	}

	msg, markup := common.BuildPagiResponse(
		ctx,
		db.Event{OwnerId: callbackQuery.From.Id},
		offset_,
		params.Pagination.Direction,
		HEADER_MESSAGE_LIST_IS_EMPTY,
		HEADER_MESSAGE_LIST_NOT_EMPTY,
	)

	event.EditCalbackMessage(
		ctx,
		msg,
		markup,
	)

	return nil
}
