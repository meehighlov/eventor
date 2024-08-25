package handlers

import (
	"context"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/internal/models"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func CreateEventForSchedulerHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	callbackQuery := event.GetCallbackQuery()
	params := models.CallbackFromString(callbackQuery.Data)

	baseFields := db.BaseFields{ID: params.Id}
	scs, err := db.Schedule{BaseFields: baseFields}.Filter(ctx)

	if err != nil {
		event.ReplyCallbackQuery(ctx, "ошибка поиска расписания " + err.Error())
		return nil
	}
	if len(scs) == 0 {
		event.ReplyCallbackQuery(ctx, "расписание не найдено")
		return nil
	}

	sc := scs[0].(db.Schedule)

	message := callbackQuery.Message

	e, err := db.BuildEvent(
		message.From.Id,
		message.GetChatIdStr(),
		sc.Text,
		"26.08 00:53", // todo
	)

	if err != nil {
		event.ReplyCallbackQuery(ctx, "ошибка создания напоминания " + err.Error())
		return nil
	}

	e.OwnerId = sc.OwnerId
	e.Save(ctx)

	sc.EventId = e.ID
	sc.Save(ctx)

	event.ReplyCallbackQuery(ctx, "событие добавлено")

	return nil
}
