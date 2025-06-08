package user

import (
	"context"

	"github.com/meehighlov/eventor/internal/repositories/models"

	"gorm.io/gorm"
)

type Filter struct {
	TgChatID string
	TgID     string
}

func (r *Repository) List(ctx context.Context, filter *Filter, tx *gorm.DB) ([]*models.User, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	var users []*models.User

	q := db.Order("created_at DESC")

	if filter.TgChatID != "" {
		q = q.Where("chat_id = ?", filter.TgChatID)
	}

	if filter.TgID != "" {
		q = q.Where("tg_id = ?", filter.TgID)
	}

	err := q.Find(&users).Error

	if err != nil {
		r.logger.Error("List error", "reason", err)
		return nil, err
	}

	r.logger.Debug("List done", "count", len(users))
	return users, nil
}
