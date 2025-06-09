package event

import (
	"context"

	"github.com/meehighlov/eventor/internal/clients/telegram"
	filter "github.com/meehighlov/eventor/internal/repositories/event"
)

func (s *Service) List(ctx context.Context, update *telegram.Update) error {
	filter := &filter.Filter{
		ChatID: update.GetChatIdStr(),
	}

	events, err := s.repositories.Event.List(ctx, filter)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		_, err := s.clients.Telegram.Reply(ctx, "Список событий пуст", update)
		return err
	}

	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()

	for _, event := range events {
		button := keyboard.NewButton(
			event.Text,
			s.builders.CallbackDataBuilder.Build(event.ID.String(), s.constants.COMMAND_INFO_EVENT).String(),
		)
		keyboard.AppendAsLine(button)
	}

	_, err = s.clients.Telegram.Reply(
		ctx,
		"Вcе события",
		update,
		telegram.WithReplyMurkup(keyboard.Murkup()),
	)
	return err
}
