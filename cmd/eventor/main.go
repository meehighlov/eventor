package main

import (
	"context"

	"github.com/meehighlov/eventor/internal/auth"
	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/internal/handlers"
	"github.com/meehighlov/eventor/internal/lib"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func main() {
	cfg := config.MustLoad()

	logger := lib.MustSetupLogging("eventor.log", true, cfg.ENV)

	db.MustSetup("eventor.db", logger)

	go handlers.RunEventPoller(context.Background(), lib.MustSetupLogging("eventor_job.log", false, cfg.ENV), cfg)

	chatCache := common.NewChatCache()

	updateHandlers := map[string]common.HandlerType{
		"/start": auth.Auth(logger, handlers.StartHandler),
		"/add": auth.Auth(logger, common.FSM(logger, chatCache, handlers.AddEventHandler())),
		"/events": auth.Auth(logger, handlers.ListEntityHandler("event")),

		// callback query handlers
		"list": handlers.ListItemCallbackQueryHandler,
		"info_event": handlers.EventInfoCallbackQueryHandler,
		"next_delta": handlers.EventInfoCallbackQueryHandler,
		"delete": handlers.DeleteItemCallbackQueryHandler,
		"conflicts": handlers.CheckConflictsCallbackHandler,
		"edit_event": common.FSM(logger, chatCache, handlers.EditEventHandlers()),
	}

	rootHandler := common.CreateRootHandler(
		logger,
		chatCache,
		updateHandlers,
	)

	logger.Info("Starting polling...")
	telegram.StartPolling(cfg.BotToken, rootHandler)
}
