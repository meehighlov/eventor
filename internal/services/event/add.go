package event

import (
	"context"

	"github.com/meehighlov/eventor/internal/clients/telegram"
)

func (s *Service) Add(ctx context.Context, update *telegram.Update) error {
	message := update.Message

	e, err := s.ParseAndBuildEvent(&message)
	if err != nil {
		return err
	}

	err = s.repositories.Event.Save(ctx, e, nil)
	if err != nil {
		return err
	}

	msg := "Событие сохранено"
	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()
	infoButton := keyboard.NewButton(
		e.Text,
		s.builders.CallbackDataBuilder.Build(e.ID.String(), s.constants.COMMAND_INFO_EVENT).String(),
	)
	keyboard.AppendAsLine(infoButton)
	s.clients.Telegram.Reply(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
