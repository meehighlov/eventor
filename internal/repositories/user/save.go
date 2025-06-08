package user

import (
	"context"

	"github.com/meehighlov/eventor/internal/repositories/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Save(ctx context.Context, user *models.User, tx *gorm.DB) error {
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
		Columns:   []clause.Column{{Name: "chat_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"tg_username", "updated_at"}),
	}).Create(user).Error
	if err != nil {
		r.logger.Error("SaveUser error", "reason", err)
		return err
	}

	r.logger.Debug("SaveUser done", "id", user.ID)
	return nil
}
