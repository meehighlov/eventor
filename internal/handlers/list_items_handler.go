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

	items, err := (&db.Event{OwnerId: callbackQuery.From.Id}).Filter(ctx)
	if err != nil {
		slog.Error("Error fetching items: " + err.Error())
		return nil
	}

	markup := common.BuildItemListMarkup(
		items,
		common.LIST_LIMIT,
		offset_,params.Pagination.Direction,
	)

	msg := HEADER_MESSAGE_LIST_NOT_EMPTY
	if len(items) == 0 {
		msg = HEADER_MESSAGE_LIST_IS_EMPTY
	}

	event.EditCalbackMessage(
		ctx,
		msg,
		markup,
	)

	return nil
}
