package handlers

import (
	"context"
	"strings"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func addEventEntry(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := []string{
		"Введи описание\n",
	}

	event.Reply(ctx, strings.Join(msg, ""))

	return 2, nil
}

func addEventSave(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	message := event.GetMessage()

	notifyAt := findNotifyAt(message.Text)

	e := db.NewEvent(
		message.From.Id,
		message.GetChatIdStr(),
		message.Text,
		notifyAt,
		"0",
	)

	e.Save(ctx)

	msg := "Событие сохранено"
	event.Reply(ctx, msg)

	return -1, nil
}

func findNotifyAt(text string) string {
	return ""
}

func AddEventHandler() map[int]telegram.CommandStepHandler {
	return map[int]telegram.CommandStepHandler{
		1: addEventEntry,
		2: addEventSave,
	}
}
