package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/db"
)

func EventInfoCallbackQueryHandler(ctx context.Context, event common.Event) error {
	callbackQuery := event.GetCallbackQuery()
	params := common.CallbackFromString(callbackQuery.Data)

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
		fmt.Sprintf("💬 `%s`", event_.Text),
	}

	if event_.NotifyNeeded() {
		msgRows = append(msgRows, fmt.Sprintf("🔔 %s", event_.NotifyAt))
		msgRows = append(msgRows, fmt.Sprintf("🔁 %s", event_.DeltaReadable()))
		nextDeltaButton := []map[string]string{
			{
				"text": event_.NextDelta(true),
				"callback_data": common.CallNextDelta(params.Id, params.Pagination.Offset).String(),
			},
		}
		markup = append(markup, nextDeltaButton)
	}

	if event_.IsScheduled() {
		msgRows = append(msgRows, fmt.Sprintf("🗓 в расписании %s", event_.Schedule))
		conflictsButton := []map[string]string{
			{
				"text": "конфликты",
				"callback_data": common.CallConflicts(params.Id).String(),
			},
		}
		markup = append(markup, conflictsButton)
	}

	editButton := []map[string]string{
		{
			"text": "редактировать",
			"callback_data": common.CallEdit(params.Id, "event").String(),
		},
	}
	toListButton := []map[string]string{
		{
			"text": "к списку",
			"callback_data": common.CallList(params.Pagination.Offset, "<", "event").String(),
		},
	}
	deleteButton := []map[string]string{
		{
			"text": "удалить",
			"callback_data": common.CallDelete(params.Id, "event").String(),
		},
	}

	markup = append(markup, editButton)
	markup = append(markup, toListButton)
	markup = append(markup, deleteButton)

	msg := strings.Join(msgRows, "\n\n")

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
