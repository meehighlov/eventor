package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/internal/models"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func DeleteItemCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	params := models.CallbackFromString(event.GetCallbackQuery().Data)

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

	sc, err := (&db.Schedule{EventId: event_.Id()}).Filter(ctx)
	if err != nil {
		if len(sc) != 0 {
			sc_ := sc[0].(db.Schedule)
			_ = sc_.UnboundEvent()
			sc_.Save(ctx)
		}
	} else {
		slog.Error("unbound schedule error:" + err.Error())
		event.ReplyCallbackQuery(ctx, "не удалось отвязать от расписания: " + sc[0].Info())
	}

	markup := [][]map[string]string{
		{
			{
				"text": "к списку",
				"callback_data": models.CallList(strconv.Itoa(common.LIST_START_OFFSET), "<", "event").String(),
			},
		},
	}

	event.EditCalbackMessage(ctx, "Событие удалено", markup)
	callBackMsg := fmt.Sprintf("Событие %s удалено", event_.Info())
	event.ReplyCallbackQuery(ctx, callBackMsg)

	return nil
}
