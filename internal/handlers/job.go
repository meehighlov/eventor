package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/pkg/telegram"
)

const CHECK_TIMEOUT_SEC = 10

func notify(ctx context.Context, client telegram.ApiCaller, events []db.Event, logger *slog.Logger) error {
	msgTemplate := "🔔 %s"
	for _, event := range events {
		msg := fmt.Sprintf(msgTemplate, event.Text)
		_, err := client.SendMessage(ctx, event.ChatId, msg)
		if err != nil {
			logger.Error("Notification not sent:" + err.Error())
		}

		event.UpdateNotifyAt()
		event.Save(ctx)
	}

	return nil
}

func run(ctx context.Context, client telegram.ApiCaller, logger *slog.Logger, cfg *config.Config) {
	logger.Info("Starting job for checking events")

	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		panic(err.Error())
	}

	for {
		// todo move datetime format to utils 
		date := time.Now().In(location).Format("02.01.2006 15:04")

		events, err := (&db.Event{}).Filter(ctx)

		notifyList := []db.Event{}

		for _, event := range events {
			if event.NotifyAt == date {
				notifyList = append(notifyList, event)
			}
		}

		if err != nil {
			logger.Error("Error getting events: " + err.Error())
		} else {
			notify(ctx, client, notifyList, logger)
		}

		time.Sleep(CHECK_TIMEOUT_SEC * time.Second)
	}
}

func RunEventPoller(
	ctx context.Context,
	logger *slog.Logger,
	cfg *config.Config,
) error {
	withCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	client := telegram.NewClient(cfg.BotToken, logger)

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("Паника в фоновой задаче проверки дней событий\n %s", r)

			logger.Error(errMsg)

			reportChatId := cfg.ReportChatId
			_, err := client.SendMessage(context.Background(), reportChatId, errMsg)
			if err != nil {
				logger.Error("panic report error:" + err.Error())
			}
		}
	}()

	run(withCancel, client, logger, cfg)

	return nil
}
