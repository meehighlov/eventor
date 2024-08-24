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

	event_ := events[0]

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
				"callback_data": models.CallList(params.Pagination.Offset, "<").String(),
			},
		},
		{
			{
				"text": "удалить",
				"callback_data": models.CallDelete(params.Id).String(),
			},
		},
	}

	event.EditCalbackMessage(ctx, msg, markup)

	return nil
}
