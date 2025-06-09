package element

import (
	"context"

	"github.com/google/uuid"
	"github.com/meehighlov/eventor/internal/repositories/models"
	"gorm.io/gorm"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).Where("id = ?", id).Delete(&models.Element{}).Error
	if err != nil {
		return err
	}

	return nil
}
