package server

import (
	"context"
	"runtime/debug"
	"slices"

	"github.com/meehighlov/eventor/internal/clients/telegram"
)

func (s *Server) HandleUpdate(ctx context.Context, update *telegram.Update) error {
	ctx, cancel := context.WithTimeout(ctx, s.handleTimeout)
	defer cancel()

	s.logger.Info("HandleUpdate", "update", update)

	username := update.Message.From.Username
	if username == "" {
		username = update.CallbackQuery.From.Username
	}

	if !slices.Contains(s.allowedUsers, username) {
		s.logger.Info("Anauthorized user", "username", username)
		return nil
	}

	chatContext := s.clients.Cache.GetOrCreateChatContext(update.GetChatIdStr())

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error(
				"Root handler",
				"recovered from panic, error", r,
				"stack", string(debug.Stack()),
				"update", update,
			)
			s.clients.Cache.Reset(update.GetChatIdStr())

			chatId := update.GetChatIdStr()
			if chatId != "" {
				s.clients.Telegram.SendMessage(ctx, chatId, s.constants.ERROR_MESSAGE)
				return
			}

			s.logger.Error(
				"Root handler",
				"recover from panic", "chatId was empty",
				"update", update,
			)
		}
	}()

	command_ := update.Message.GetCommand()
	command := ""

	if command_ != "" {
		command = command_
		s.clients.Cache.Reset(update.GetChatIdStr())
	} else {
		if update.CallbackQuery.Id != "" {
			params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

			s.clients.Telegram.AnswerCallbackQuery(ctx, update.CallbackQuery.Id)

			s.logger.Info("CallbackQueryHandler", "command", params.Command, "chat id", update.GetChatIdStr())
			command = params.Command
		} else {
			command_ = chatContext.GetNextHandler()
			if command_ != "" {
				command = command_
			}
		}
	}

	s.logger.Info("Root handler", "handling", command, "with update", update)

	err := s.handle(ctx, update, command)
	if err != nil {
		s.clients.Cache.Reset(update.GetChatIdStr())
		s.logger.Error("Root handler", "error", err.Error(), "chat id", update.GetChatIdStr(), "update id", update.UpdateId)
	} else {
		s.logger.Info("Root handler", "success", command, "chat id", update.GetChatIdStr(), "update id", update.UpdateId)
	}

	return nil
}

func (s *Server) handle(ctx context.Context, update *telegram.Update, command string) error {
	switch command {
	case s.constants.COMMAND_START:
		return s.services.User.Start(ctx, update)
	case s.constants.COMMAND_EVENTS:
		return s.services.Event.List(ctx, update)
	case s.constants.COMMAND_DELETE:
		return s.services.Event.Delete(ctx, update)
	case s.constants.COMMAND_LIST_EVENT:
		return s.services.Event.List(ctx, update)
	case s.constants.COMMAND_EDIT_EVENT:
		return s.services.Event.Edit(ctx, update)
	case s.constants.COMMAND_EDIT_EVENT_SAVE:
		return s.services.Event.EditSave(ctx, update)
	case s.constants.COMMAND_INFO_EVENT:
		return s.services.Event.Info(ctx, update)
	case s.constants.COMMAND_NEXT_DELTA:
		return s.services.Event.Info(ctx, update)
	case s.constants.COMMAND_ELEMENTS, s.constants.COMMAND_LIST_ELEMENT:
		return s.services.Element.List(ctx, update)
	case s.constants.COMMAND_ADD_ELEMENT:
		return s.services.Element.Add(ctx, update)
	case s.constants.COMMAND_ADD_ELEMENT_SAVE:
		return s.services.Element.AddSave(ctx, update)
	case s.constants.COMMAND_INFO_ELEMENT, s.constants.COMMAND_ELEMENT_SWITCH_STATUS:
		return s.services.Element.Info(ctx, update)
	case s.constants.COMMAND_EDIT_ELEMENT:
		return s.services.Element.Edit(ctx, update)
	case s.constants.COMMAND_EDIT_ELEMENT_NAME_SAVE:
		return s.services.Element.EditNameSave(ctx, update)
	case s.constants.COMMAND_EDIT_ELEMENT_LINK_SAVE:
		return s.services.Element.EditLinkSave(ctx, update)
	case s.constants.COMMAND_EDIT_ELEMENT_REQUEST:
		return s.services.Element.EditRequest(ctx, update)
	case s.constants.COMMAND_DELETE_ELEMENT:
		return s.services.Element.Delete(ctx, update)
	case s.constants.COMMAND_DELETE_ELEMENT_CONFIRM:
		return s.services.Element.DeleteConfirm(ctx, update)
	default:
		return s.services.Event.Add(ctx, update)
	}
}
