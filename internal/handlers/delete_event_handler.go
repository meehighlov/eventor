package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
)

func DeleteItemCallbackQueryHandler(event common.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	baseFields := db.BaseFields{ID: params.Id}
	events, err := (&db.Event{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		event.ReplyCallbackQuery(ctx, "Возникла непредвиденная ошибка")
		slog.Error("error serching friend when deleting: " + err.Error())
	}

	if len(events) == 0 {
		slog.Error("not found event row by id: " + params.Id)
		return err
	}

	event_ := events[0]

	event_.Delete(ctx)

	markup := [][]map[string]string{
		{
			{
				"text": "к списку",
				"callback_data": common.CallList(strconv.Itoa(common.LIST_START_OFFSET), "<", "event").String(),
			},
		},
	}

	event.EditCalbackMessage(ctx, "Событие удалено", markup)
	callBackMsg := fmt.Sprintf("Событие %s удалено", event_.Info())
	event.ReplyCallbackQuery(ctx, callBackMsg)

	return nil
}
