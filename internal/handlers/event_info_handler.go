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

	if params.Command == "next_delta" {
		event_.Delta = event_.NextDelta(false)
		err := event_.Save(ctx)
		if err != nil {
			slog.Error("error nexting delta: " + err.Error())
			event.ReplyCallbackQuery(ctx, "не удалось обновить повтор")
		}
	}

	markup := [][]map[string]string{}

	msgRows := []string{
		fmt.Sprintf("💬 %s", event_.Text),
	}

	if event_.NotifyNeeded() {
		msgRows = append(msgRows, fmt.Sprintf("🔔 %s", event_.NotifyAt))
		msgRows = append(msgRows, fmt.Sprintf("🔁 %s", event_.DeltaReadable()))
		nextDeltaButton := []map[string]string{
			{
				"text": event_.NextDelta(true),
				"callback_data": models.CallNextDelta(params.Id, params.Pagination.Offset).String(),
			},
		}
		markup = append(markup, nextDeltaButton)
	}

	toListButton := []map[string]string{
		{
			"text": "к списку",
			"callback_data": models.CallList(params.Pagination.Offset, "<", "event").String(),
		},
	}
	deleteButton := []map[string]string{
		{
			"text": "удалить",
			"callback_data": models.CallDelete(params.Id, "event").String(),
		},
	}

	markup = append(markup, toListButton)
	markup = append(markup, deleteButton)

	msg := strings.Join(msgRows, "\n\n")

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
