package event

import (
	"context"
	"fmt"
	"strings"

	"github.com/meehighlov/eventor/internal/clients/telegram"
)

func (s *Service) Info(ctx context.Context, update *telegram.Update) error {
	callbackQuery := update.CallbackQuery
	params := s.builders.CallbackDataBuilder.FromString(callbackQuery.Data)

	event, err := s.repositories.Event.Get(ctx, params.ID, nil)
	if err != nil {
		return err
	}

	if params.Command == s.constants.COMMAND_NEXT_DELTA {
		event.Delta = event.NextDelta(false)
		err := s.repositories.Event.Save(ctx, event, nil)
		if err != nil {
			s.clients.Telegram.Reply(ctx, "не удалось обновить повтор", update)
		}
	}

	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()

	msgRows := []string{
		fmt.Sprintf("💬 `%s`", event.Text),
	}

	if event.NotifyNeeded() {
		msgRows = append(msgRows, fmt.Sprintf("🔔 %s", event.NotifyAt))
		msgRows = append(msgRows, fmt.Sprintf("🔁 %s", event.DeltaReadable()))

		nextDeltaButton := keyboard.NewButton(
			"🔁",
			s.builders.CallbackDataBuilder.Build(event.ID.String(), s.constants.COMMAND_NEXT_DELTA).String(),
		)
		keyboard.AppendAsLine(nextDeltaButton)
	}

	if event.IsScheduled() {
		msgRows = append(msgRows, fmt.Sprintf("🗓 в расписании %s", event.Schedule))
	}

	editButton := keyboard.NewButton(
		"✏️",
		s.builders.CallbackDataBuilder.Build(event.ID.String(), s.constants.COMMAND_EDIT_EVENT).String(),
	)

	toListButton := keyboard.NewButton(
		"⬅️",
		s.builders.CallbackDataBuilder.Build(event.ID.String(), s.constants.COMMAND_LIST_EVENT).String(),
	)

	deleteButton := keyboard.NewButton(
		"🗑",
		s.builders.CallbackDataBuilder.Build(event.ID.String(), s.constants.COMMAND_DELETE).String(),
	)

	keyboard.AppendAsLine(toListButton, editButton, deleteButton)

	msg := strings.Join(msgRows, "\n\n")

	s.clients.Telegram.Edit(ctx, msg, update, keyboard.Murkup(), telegram.WithMarkDown())

	return nil
}
