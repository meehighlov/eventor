package element

import (
	"context"

	"github.com/google/uuid"
	"github.com/meehighlov/eventor/internal/clients/telegram"
	"github.com/meehighlov/eventor/internal/repositories/models"
	"github.com/meehighlov/eventor/internal/repositories/user"
)

func (s *Service) Add(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(
		ctx,
		"Введите название элемента",
		update,
		telegram.WithReplyMurkup(s.builders.KeyboardBuilder.BuildInlineKeyboard().Murkup()),
	)

	s.clients.Cache.SetNextHandler(update.GetChatIdStr(), s.constants.COMMAND_ADD_ELEMENT_SAVE)

	return nil
}

func (s *Service) AddSave(ctx context.Context, update *telegram.Update) error {
	user, err := s.repositories.User.Get(ctx, &user.Filter{TgChatID: update.GetChatIdStr()}, nil)
	if err != nil {
		return err
	}

	element := &models.Element{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   update.Message.Text,
	}

	err = s.repositories.Element.Save(ctx, element, nil)
	if err != nil {
		return err
	}

	s.clients.Cache.SetNextHandler(update.GetChatIdStr(), "")

	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_INFO_ELEMENT, s.builders.CallbackDataBuilder.Build(element.ID.String(), s.constants.COMMAND_INFO_ELEMENT).String()))
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_ADD_ELEMENT, s.builders.CallbackDataBuilder.Build(user.ID.String(), s.constants.COMMAND_ADD_ELEMENT).String()))

	s.clients.Telegram.Reply(ctx, "Элемент добавлен", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
