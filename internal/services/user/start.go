package user

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/meehighlov/eventor/internal/clients/telegram"
	"github.com/meehighlov/eventor/internal/repositories/models"
)

func (s *Service) Start(ctx context.Context, update *telegram.Update) error {
	err := s.repositories.User.DB().Transaction(func(tx *gorm.DB) error {
		err := s.repositories.User.Save(ctx, &models.User{
			ID:         uuid.New(),
			TgID:       strconv.Itoa(update.Message.From.Id),
			TgUsername: update.Message.From.Username,
			TgChatID:   update.GetChatIdStr(),
		}, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.clients.Telegram.Reply(ctx, s.constants.START_MESSAGE, update)
	return nil
}
