package handlers

import (
	"context"
	"log/slog"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
)


func editStart(event common.Event) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	event.AnswerCallbackQuery(ctx)

	event.ReplyCallbackQuery(ctx, "Введи измененный текст")
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	event.GetContext().AppendText(params.Id)

	return "2", nil
}

func editSave(event common.Event) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	// message we accept here is not bound to callback query
	// so get params we saved in previous step from cache
	eventId := event.GetContext().GetTexts()[0]

	baseFields := db.BaseFields{ID: eventId}
	events, err := (&db.Event{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		event.Reply(ctx, "Возникла непредвиденная ошибка")
		slog.Error("error serching friend when deleting: " + err.Error())
		return common.STEPS_DONE, err
	}

	if len(events) == 0 {
		event.Reply(ctx, "Возникла непредвиденная ошибка, не найдены события")
		slog.Error("not found event row by id: " + eventId)
		return common.STEPS_DONE, err
	}

	event_ := events[0].(db.Event)

	message := event.GetMessage()

	updatedEvent := ParseAndBuildEvent(message)

	err = event_.Delete(ctx)
	if err != nil {
		event.Reply(ctx, err.Error())
		return common.STEPS_DONE, err
	}

	updatedEvent.Save(ctx)

	event.Reply(ctx, "событие обновлено")

	return common.STEPS_DONE, nil
}

func EditEventHandlers() map[string]common.CommandStepHandler {
	return map[string]common.CommandStepHandler{
		"1": editStart,
		"2": editSave,
	}
}
