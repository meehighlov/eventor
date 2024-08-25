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

func ScheduleInfoCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	callbackQuery := event.GetCallbackQuery()
	params := models.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	scs, err := (&db.Schedule{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		slog.Error("error during fetching event info: " + err.Error())
		return nil
	}

	sc, ok := scs[0].(db.Schedule)
	if !ok {
		slog.Error("cast from Item to Event error")
		event.EditCalbackMessage(ctx, "непредвиденная ошибка", [][]map[string]string{})
		return nil
	}

	msg := strings.Join(
		[]string{
			fmt.Sprintf("💬 %s", sc.Text),
			fmt.Sprintf("🔔 %s", sc.DeltaReadable()),
			fmt.Sprintf("🔁 %s", sc.Day),
		},
		"\n\n",
	)

	markup := [][]map[string]string{
		{
			{
				"text": "к списку",
				"callback_data": models.CallList(params.Pagination.Offset, "<", sc.Name()).String(),
			},
		},
		{
			{
				"text": "удалить",
				"callback_data": models.CallDelete(params.Id, sc.Name()).String(),
			},
		},
	}

	if !sc.HasNotifications() {
		btn := []map[string]string{
			{
				"text": "напомнить",
				"callback_data": models.CallCreateEventForSchedule(sc.Id()).String(),
			},
		}
		markup = append(markup, btn)
	}

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
