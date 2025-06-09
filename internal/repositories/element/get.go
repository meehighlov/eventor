package element

import (
	"context"
	"errors"

	"github.com/meehighlov/eventor/internal/repositories/models"
	"gorm.io/gorm"
)

func (r *Repository) Get(ctx context.Context, filter *Filter, tx *gorm.DB) (*models.Element, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	element := &models.Element{}

	if filter == nil {
		return nil, errors.New("filter is nil")
	}
	
	q := db.Order("created_at DESC")

	if filter.ID != "" {
		q = q.Where("id = ?", filter.ID)
	}

	err := q.First(element).Error
	if err != nil {
		return nil, err
	}

	return element, nil
}
