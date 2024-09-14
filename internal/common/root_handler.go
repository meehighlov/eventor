package common

import (
	"log/slog"

	"github.com/meehighlov/eventor/pkg/telegram"
)

type HandlerType func(Event) error


func CreateRootHandler(logger *slog.Logger, chatCahe *ChatCache, handlers map[string]HandlerType) telegram.UpdateHandler {
	return func(update telegram.Update, client telegram.ApiCaller) error {
		chatContext := chatCahe.GetOrCreateChatContext(update.GetChatIdStr())
		command_ := update.Message.GetCommand()
		command := ""

		if command_ != "" {
			command = command_
			chatContext.Reset()
			logger.Debug("resetting context due to message command priority", "command:", command)
		} else {
			if update.CallbackQuery.Id != "" {
				params := CallbackFromString(update.CallbackQuery.Data)

				logger.Debug("CallbackQueryHandler", "command", params.Command, "entity", params.Entity)
				command = params.Command
			} else {
				command_ = chatContext.GetCommandInProgress()
				logger.Debug("command in progress from chat context", "command:", command_)
				if command_ != "" {
					command = command_
				}
			}
		}

		event := newEvent(client, update, chatContext, command)

		logger.Debug("invoking event", "command", command)

		handler, found := handlers[command]
		if found {
			handler(event)
		} else {
			logger.Debug("handler not found", "command:", command)
		}

		return nil
	}
}
