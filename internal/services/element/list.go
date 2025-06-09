package element

import (
	"context"

	"github.com/meehighlov/eventor/internal/clients/telegram"
	"github.com/meehighlov/eventor/internal/repositories/element"
	"github.com/meehighlov/eventor/internal/repositories/user"
)

func (s *Service) List(ctx context.Context, update *telegram.Update) error {
	user, err := s.repositories.User.Get(ctx, &user.Filter{TgChatID: update.GetChatIdStr()}, nil)
	if err != nil {
		return err
	}

	filter := &element.Filter{
		UserID: user.ID.String(),
	}

	elements, err := s.repositories.Element.List(ctx, filter)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()

	addButton := keyboard.NewButton(s.constants.BUTTON_TEXT_ADD_ELEMENT, s.builders.CallbackDataBuilder.Build(user.ID.String(), s.constants.COMMAND_ADD_ELEMENT).String())
	keyboard.AppendAsLine(addButton)

	for _, element := range elements {
		keyboard.AppendAsLine(keyboard.NewButton(element.Name, s.builders.CallbackDataBuilder.Build(element.ID.String(), s.constants.COMMAND_ELEMENT_INFO).String()))
	}

	if update.IsCallback() {
		_, err := s.clients.Telegram.Edit(ctx, "Мои элементы", update, telegram.WithReplyMurkup(keyboard.Murkup()), telegram.WithMarkDown())
		return err
	}

	s.clients.Telegram.Reply(ctx, "Мои элементы", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
