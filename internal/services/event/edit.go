package event

import (
	"context"

	"github.com/meehighlov/eventor/internal/clients/telegram"
)

func (s *Service) Edit(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(ctx, "Введи измененный текст", update)
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	s.clients.Cache.AppendText(update.GetChatIdStr(), params.ID)

	s.clients.Cache.SetNextHandler(update.GetChatIdStr(), s.constants.COMMAND_EDIT_EVENT_SAVE)

	return nil
}

func (s *Service) EditSave(ctx context.Context, update *telegram.Update) error {
	eventId := s.clients.Cache.GetTexts(update.GetChatIdStr())[0]

	event, err := s.repositories.Event.Get(ctx, eventId, nil)

	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	message := update.Message

	updatedEvent, err := s.ParseAndBuildEvent(&message)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	updatedEvent.ID = event.ID
	updatedEvent.Delta = event.Delta

	err = s.repositories.Event.Save(ctx, updatedEvent, nil)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	s.clients.Telegram.Reply(ctx, "событие обновлено", update)

	s.clients.Cache.SetNextHandler(update.GetChatIdStr(), "")

	return nil
}
