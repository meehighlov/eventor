package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/internal/models"
	"github.com/meehighlov/eventor/pkg/telegram"
)

const CHECK_TIMEOUT_SEC = 10


func buildNotificationButtons(eventId string) [][]map[string]string {
	return [][]map[string]string{{
		{
			"text": "удалить",
			"callback_data": models.CallDelete(eventId, "event").String(),
		},
	},
	}
}

func notify(ctx context.Context, client telegram.ApiCaller, events []db.Event, logger *slog.Logger) error {
	msgTemplate := "🔔 %s"
	for _, event := range events {
		msg := fmt.Sprintf(msgTemplate, event.Text)
		_, err := client.SendMessageWithReplyMarkup(ctx, event.ChatId, msg, buildNotificationButtons(event.ID))
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
		date := time.Now().In(location).Format("02.01 15:04")

		events, err := (&db.Event{}).Filter(ctx)

		notifyList := []db.Event{}

		for _, event := range events {
			event_, ok := event.(db.Event)

			if !ok {
				errMsg := "cannot cast Item to Event in background job for checking events"
				logger.Error(errMsg)
				_, err := client.SendMessage(context.Background(), cfg.ReportChatId, errMsg)
				if err != nil {
					logger.Error("error sending report message: " + errMsg)
				}
				continue
			}

			if event_.NotifyAt == date {
				notifyList = append(notifyList, event_)
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

			_, err := client.SendMessage(context.Background(), cfg.ReportChatId, errMsg)
			if err != nil {
				logger.Error("panic report error:" + err.Error())
			}
		}
	}()

	run(withCancel, client, logger, cfg)

	return nil
}
