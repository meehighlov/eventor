package user

import (
	"context"
	"errors"

	"github.com/meehighlov/eventor/internal/repositories/models"

	"gorm.io/gorm"
)

func (r *Repository) Get(ctx context.Context, filter *Filter, tx *gorm.DB) (*models.User, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	user := &models.User{}

	q := db.Order("created_at DESC")

	if filter == nil {
		return nil, errors.New("filter for user is nil")
	}

	if filter.TgChatID != "" {
		q = q.Where("chat_id = ?", filter.TgChatID)
	}

	err := q.First(user).Error
	if err != nil {
		r.logger.Error("Get error", "reason", err, "id", filter.TgChatID)
		return nil, err
	}

	r.logger.Debug("Get done", "id", user.ID)
	return user, nil
}
