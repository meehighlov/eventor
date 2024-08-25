package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/internal/models"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func EventInfoCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	callbackQuery := event.GetCallbackQuery()
	params := models.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	events, err := (&db.Event{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		slog.Error("error during fetching event info: " + err.Error())
		return nil
	}

	if len(events) == 0 {
		event.EditCalbackMessage(ctx, "непредвиденная ошибка: события не найдены", [][]map[string]string{})
		return nil
	}

	event_, ok := events[0].(db.Event)
	if !ok {
		slog.Error("cast from Item to Event error")
		event.EditCalbackMessage(ctx, "непредвиденная ошибка", [][]map[string]string{})
		return nil
	}

	msg := strings.Join(
		[]string{
			fmt.Sprintf("💬 %s", event_.Text),
			fmt.Sprintf("🔔 %s", event_.NotifyAt),
			fmt.Sprintf("🔁 %s", event_.DeltaReadable()),
		},
		"\n\n",
	)

	markup := [][]map[string]string{
		{
			{
				"text": "к списку",
				"callback_data": models.CallList(params.Pagination.Offset, "<", "event").String(),
			},
		},
		{
			{
				"text": "удалить",
				"callback_data": models.CallDelete(params.Id, "event").String(),
			},
		},
	}

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
