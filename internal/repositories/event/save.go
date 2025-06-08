package event

import (
	"context"

	"github.com/meehighlov/eventor/internal/repositories/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Save(ctx context.Context, event *models.Event, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"notify_at", "updated_at", "schedule", "delta", "text", "chat_id", "owner_id"}),
	}).Create(event).Error
	if err != nil {
		r.logger.Error("SaveEvent error", "reason", err)
		return err
	}

	r.logger.Debug("SaveEvent done", "id", event.ID)
	return nil
}
