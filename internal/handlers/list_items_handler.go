package handlers

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/internal/models"
	"github.com/meehighlov/eventor/pkg/telegram"
)

const (
	HEADER_MESSAGE_LIST_NOT_EMPTY = "Нажми, чтобы узнать детали✨"
	HEADER_MESSAGE_LIST_IS_EMPTY = "Записей пока нет✨"
)

func ListEntityHandler(entity string) telegram.CommandHandler {
	return func (event telegram.Event) error {
		ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
		defer cancel()
	
		message := event.GetMessage()

		items, err := (buildItem(entity, message.From.Id)).Filter(ctx)
	
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

func ListItemCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()
	callbackQuery := event.GetCallbackQuery()

	params := models.CallbackFromString(event.GetCallbackQuery().Data)

	offset := params.Pagination.Offset

	offset_, err := strconv.Atoi(offset)
	if err != nil {
		slog.Error("error parsing offset in list pagination callback query: " + err.Error())
		return err
	}

	msg, markup := buildResponse(
		ctx,
		params.Entity,
		callbackQuery.From.Id,
		offset_,
		params.Pagination.Direction,
	)

	event.EditCalbackMessage(
		ctx,
		msg,
		markup,
	)

	return nil
}

func buildResponse(
	ctx context.Context,
	entity string,
	ownerId int,
	offset int,
	direction string,
) (string, [][]map[string]string) {
	hideMarkup := [][]map[string]string{}
	var msgByItemsLen = func(itemsLen int) string {
		msg := HEADER_MESSAGE_LIST_NOT_EMPTY
		if itemsLen == 0 {
			msg = HEADER_MESSAGE_LIST_IS_EMPTY
		}
		return msg
	}

	items, err := buildItem(entity, ownerId).Filter(ctx)
	if err != nil {
		slog.Error("Error fetching items: " + err.Error())
		return "Не могу разобрать запрос", hideMarkup
	}
	return msgByItemsLen(len(items)), common.BuildItemListMarkup(
		items,
		common.LIST_LIMIT,
		offset,
		direction,
		entity,
	)
}

func buildItem(entity string, ownerId int) common.Item {
	var item common.Item
	if entity == "event" {
		item = db.Event{OwnerId: ownerId}
	}
	if entity == "schedule" {
		item = db.Schedule{OwnerId: ownerId}
	}
	return item
}
