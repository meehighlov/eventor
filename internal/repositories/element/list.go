package element

import (
	"context"

	"github.com/meehighlov/eventor/internal/repositories/models"
)

type Filter struct {
	ID     string
	UserID string
	Status string
}

func (r *Repository) List(ctx context.Context, filter *Filter) ([]*models.Element, error) {
	var elements []*models.Element

	query := r.DB().WithContext(ctx)

	if filter != nil {
		if filter.UserID != "" {
			query = query.Where("user_id = ?", filter.UserID)
		}
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
	}

	err := query.Find(&elements).Error
	if err != nil {
		return nil, err
	}
	return elements, nil
}
