package element

import (
	"context"
	"fmt"

	"github.com/meehighlov/eventor/internal/clients/telegram"
	"github.com/meehighlov/eventor/internal/repositories/element"
)

func (s *Service) Info(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	element, err := s.repositories.Element.Get(ctx, &element.Filter{ID: params.ID}, nil)
	if err != nil {
		return err
	}

	if params.Command == s.constants.COMMAND_ELEMENT_SWITCH_STATUS {
		element.Status = element.NextStatus()
		err = s.repositories.Element.Save(ctx, element, nil)
		if err != nil {
			return err
		}
	}

	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()
	buttonBack := keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_LIST_ELEMENT).String())
	buttonEdit := keyboard.NewButton(s.constants.BUTTON_TEXT_EDIT, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_EDIT_ELEMENT).String())
	buttonDelete := keyboard.NewButton(s.constants.BUTTON_TEXT_DELETE, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_DELETE_ELEMENT).String())

	header := fmt.Sprintf("💪 %s\n\n", element.Name)
	header += fmt.Sprintf("%s\n\n", element.ElementReadableStatus(""))
	if element.TutorialLink != "" {
		buttonTutorial := keyboard.NewURLButton(s.constants.BUTTON_TEXT_TUTORIAL, element.TutorialLink)
		keyboard.AppendAsLine(buttonTutorial)
	}

	buttonSwitchStatus := keyboard.NewButton(element.ElementReadableStatus(element.NextStatus()), s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_ELEMENT_SWITCH_STATUS).String())
	keyboard.AppendAsLine(buttonSwitchStatus)

	keyboard.AppendAsLine(buttonBack, buttonEdit, buttonDelete)

	s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
