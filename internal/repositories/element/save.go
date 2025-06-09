package element

import (
	"context"

	"github.com/meehighlov/eventor/internal/repositories/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Save(ctx context.Context, element *models.Element, tx *gorm.DB) error {
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
		DoUpdates: clause.AssignmentColumns([]string{"name", "tutorial_link", "status", "updated_at"}),
	}).Create(element).Error
	if err != nil {
		r.logger.Error("SaveElement error", "reason", err)
		return err
	}

	r.logger.Debug("SaveElement done", "id", element.ID)
	return nil
}
