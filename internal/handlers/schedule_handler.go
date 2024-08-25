package handlers

import (
	"context"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func scheduleEntry(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := "Введи описание события"

	event.Reply(ctx, msg)

	return 2, nil
}

func scheduleAcceptTimestamp(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	event.GetContext().AppendText(event.GetMessage().Text)

	msg := "Введи день недели и время"

	event.Reply(ctx, msg)

	return 3, nil
}

func scheduleSave(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	timestamp := event.GetMessage().Text
	eventText := event.GetContext().GetTexts()[0]

	message := event.GetMessage()

	s := db.NewSchedule(
		message.From.Id,
		message.GetChatIdStr(),
		eventText,
		"d",
		timestamp,
	)

	s.Save(ctx)

	msg := "Расписание обновлено"
	event.Reply(ctx, msg)

	return -1, nil
}

func ScheduleHandler() map[int]telegram.CommandStepHandler {
	return map[int]telegram.CommandStepHandler{
		1: scheduleEntry,
		2: scheduleAcceptTimestamp,
		3: scheduleSave,
	}
}
