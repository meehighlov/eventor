package handlers

import (
	"context"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func addEventEntry(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := "Введи описание события"

	event.Reply(ctx, msg)

	return 2, nil
}

func addEventAccepTimestamp(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	event.GetContext().AppendText(event.GetMessage().Text)

	msg := "Введи дату и время события"

	event.Reply(ctx, msg)

	return 3, nil
}

func addEventSave(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	timestamp := event.GetMessage().Text
	eventText := event.GetContext().GetTexts()[0]

	message := event.GetMessage()

	e, err := db.BuildEvent(
		message.From.Id,
		message.GetChatIdStr(),
		eventText,
		timestamp,
		"h",
	)

	if err != nil {
		event.Reply(ctx, "Ошибка добавления события: " + err.Error() +"\nпопробуй снова")
		return 3, nil
	}

	e.Save(ctx)

	msg := "Событие сохранено"
	event.Reply(ctx, msg)

	return -1, nil
}

func AddEventHandler() map[int]telegram.CommandStepHandler {
	return map[int]telegram.CommandStepHandler{
		1: addEventEntry,
		2: addEventAccepTimestamp,
		3: addEventSave,
	}
}
