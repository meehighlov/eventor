package event

import (
	"context"

	"github.com/meehighlov/eventor/internal/repositories/models"
)

type Filter struct {
	ID      string
	OwnerID string
	ChatID  string
}

func (r *Repository) List(ctx context.Context, filter *Filter) ([]*models.Event, error) {
	var events []*models.Event

	query := r.DB().WithContext(ctx)

	if filter != nil && filter.ChatID != "" {
		query = query.Where("chat_id = ?", filter.ChatID)
	}

	err := query.Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}
