package element

import (
	"context"
	"fmt"

	"github.com/meehighlov/eventor/internal/clients/telegram"
	"github.com/meehighlov/eventor/internal/repositories/element"
)

func (s *Service) Delete(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	elementToDelete, err := s.repositories.Element.Get(ctx, &element.Filter{ID: params.ID}, nil)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()
	keyboard.AppendAsLine(
		keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(elementToDelete.ID.String(), s.constants.COMMAND_INFO_ELEMENT).String()),
		keyboard.NewButton(s.constants.BUTTON_TEXT_DELETE, s.builders.CallbackDataBuilder.Build(elementToDelete.ID.String(), s.constants.COMMAND_DELETE_ELEMENT_CONFIRM).String()),
	)

	header := fmt.Sprintf("Вы уверены, что хотите удалить %s?", elementToDelete.Name)
	s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}

func (s *Service) DeleteConfirm(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	elementToDelete, err := s.repositories.Element.Get(ctx, &element.Filter{ID: params.ID}, nil)
	if err != nil {
		return err
	}

	err = s.repositories.Element.Delete(ctx, elementToDelete.ID, nil)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(elementToDelete.ID.String(), s.constants.COMMAND_LIST_ELEMENT).String()))

	header := fmt.Sprintf("Элемент %s удален", elementToDelete.Name)
	s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
