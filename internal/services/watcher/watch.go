package watcher

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/meehighlov/eventor/internal/clients/telegram"
	"github.com/meehighlov/eventor/internal/repositories/models"
)

const CHECK_TIMEOUT_SEC = 10

func (s *Service) buildNotificationButtons(eventId string) []*[]map[string]interface{} {
	keyboard := s.builders.KeyboardBuilder.BuildInlineKeyboard()

	detailsButton := keyboard.NewButton(
		"детали",
		s.builders.CallbackDataBuilder.Build(eventId, s.constants.COMMAND_INFO_EVENT).String(),
	)

	deleteButton := keyboard.NewButton(
		"удалить",
		s.builders.CallbackDataBuilder.Build(eventId, s.constants.COMMAND_DELETE).String(),
	)

	keyboard.AppendAsLine(detailsButton, deleteButton)

	return keyboard.Murkup()
}

func (s *Service) notify(ctx context.Context, events []*models.Event) error {
	msgTemplate := "🔔 %s"
	for _, event := range events {
		_, err := event.UpdateNotifyAt()
		if err != nil {
			s.logger.Error("Error updating notify at:" + err.Error())
			continue
		}

		err = s.repositories.Event.Save(ctx, event, nil)
		if err != nil {
			s.logger.Error("Error saving event:" + err.Error())
			continue
		}

		msg := fmt.Sprintf(msgTemplate, event.Text)
		_, err = s.clients.Telegram.SendMessage(
			ctx,
			event.ChatID,
			msg,
			telegram.WithReplyMurkup(s.buildNotificationButtons(event.ID.String())),
		)
		if err != nil {
			s.logger.Error("Notification not sent:" + err.Error())
		}
	}

	return nil
}

func (s *Service) run(ctx context.Context) error {
	s.logger.Debug("Checking events")

	// todo move datetime format to utils
	date := time.Now().In(s.location).Format("02.01 15:04")

	events, err := s.repositories.Event.List(ctx, nil)

	if err != nil {
		s.logger.Error("Error getting events: " + err.Error())
		return err
	}

	notifyList := []*models.Event{}

	for _, event := range events {

		if event.NotifyAt == date {
			notifyList = append(notifyList, event)
		}
	}

	s.notify(ctx, notifyList)

	return nil
}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("Starting job for checking events")

	ticker := time.NewTicker(time.Duration(s.CheckIntervalSec) * time.Second)
	defer ticker.Stop()

	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 8192)
			stack = stack[:runtime.Stack(stack, false)]
			s.logger.Error("Watcher recovered from panic",
				"error", r,
				"stack", string(stack))
		}
	}()

	for {
		select {
		case <-ticker.C:
			checkCtx, cancel := context.WithTimeout(ctx, time.Duration(s.CheckIntervalSec/3)*time.Second)
			err := s.run(checkCtx)
			cancel()

			if err != nil {
				s.logger.Error("Watch", "error checking tracks", err)
			}
		case <-ctx.Done():
			s.logger.Debug("Watch", "context done", ctx.Err())
			return ctx.Err()
		}
	}
}
