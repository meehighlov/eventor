package event

import (
	"context"
	"fmt"

	"github.com/meehighlov/eventor/internal/clients/telegram"
)

func (s *Service) Delete(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	event, err := s.repositories.Event.Get(ctx, params.ID, nil)
	if err != nil {
		return err
	}

	err = s.repositories.Event.Delete(ctx, event.ID)
	if err != nil {
		return err
	}

	markup := s.builders.KeyboardBuilder.BuildInlineKeyboard()
	button := markup.NewButton("к списку", s.builders.CallbackDataBuilder.Build(event.ID.String(), "event_list").String())
	markup.AppendAsLine(button)

	_, err = s.clients.Telegram.Edit(ctx, "Событие удалено", update, telegram.WithReplyMurkup(markup.Murkup()))
	if err != nil {
		return err
	}

	callBackMsg := fmt.Sprintf("Событие %s удалено", event.Text)
	s.clients.Telegram.Reply(ctx, callBackMsg, update)

	return nil
}
