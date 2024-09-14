package handlers

import (
	"context"
	"strings"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
)

func addEventEntry(event common.Event) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := []string{
		"Введи описание\n",
	}

	event.Reply(ctx, strings.Join(msg, ""))

	return "2", nil
}

func addEventSave(event common.Event) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	message := event.GetMessage()

	e := ParseAndBuildEvent(message)

	e.Save(ctx)

	msg := "Событие сохранено"
	event.Reply(ctx, msg)

	return common.STEPS_DONE, nil
}

func AddEventHandler() map[string]common.CommandStepHandler {
	return map[string]common.CommandStepHandler{
		"1": addEventEntry,
		"2": addEventSave,
	}
}
