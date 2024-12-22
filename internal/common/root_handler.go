package common

import (
	"context"
	"log/slog"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/pkg/telegram"
)

type HandlerType func(context.Context, Event) error


func CreateRootHandler(logger *slog.Logger, chatCahe *ChatCache, handlers map[string]HandlerType) telegram.UpdateHandler {
	return func(update telegram.Update, client telegram.ApiCaller) error {
		ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
		defer cancel()

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

				client.AnswerCallbackQuery(ctx, update.CallbackQuery.Id)

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
			handler(ctx, event)
		} else {
			logger.Debug("handler not found, invoking default handler")
			if defaultHanlder, foundDefault := handlers["default"]; foundDefault {
				defaultHanlder(ctx, event)
			} else {
				logger.Debug("default handler not found, skipping")
			}
		}

		return nil
	}
}
