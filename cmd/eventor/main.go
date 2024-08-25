package main

import (
	"context"

	"github.com/meehighlov/eventor/internal/auth"
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

	bot := telegram.NewBot(cfg.BotToken, nil)

	bot.RegisterCommandHandler("/start", auth.Auth(handlers.StartHandler))
	bot.RegisterCommandHandler("/commands", auth.Auth(handlers.CommandsHandler))
	bot.RegisterCommandHandler("/add_event", auth.Auth(telegram.FSM(handlers.AddEventHandler())))
	bot.RegisterCommandHandler("/events", auth.Auth(handlers.ListEntityHandler("event")))
	bot.RegisterCommandHandler("/schedule", auth.Auth(handlers.ListEntityHandler("schedule")))
	bot.RegisterCommandHandler("/add_schedule", auth.Auth(telegram.FSM(handlers.AddScheduleHandler())))

	bot.RegisterCallbackQueryHandler(handlers.CallbackQueryHandler)

	logger.Info("Starting polling...")
	bot.StartPolling()
}
