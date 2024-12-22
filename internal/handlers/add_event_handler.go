package handlers

import (
	"context"

	"github.com/meehighlov/eventor/internal/common"
)

func AddEventSave(ctx context.Context, event common.Event) error {
	message := event.GetMessage()

	e := ParseAndBuildEvent(message)

	e.Save(ctx)

	msg := "Событие сохранено"
	event.Reply(ctx, msg)

	return nil
}
