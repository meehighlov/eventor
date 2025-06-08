package event

import (
	"context"

	"github.com/meehighlov/eventor/internal/repositories/models"
	"gorm.io/gorm"
)

func (r *Repository) Get(ctx context.Context, id string, tx *gorm.DB) (*models.Event, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	event := &models.Event{}

	err := db.Where("id = ?", id).First(event).Error
	if err != nil {
		return nil, err
	}

	return event, nil
}
