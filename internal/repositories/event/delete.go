package event

import (
	"context"

	"github.com/google/uuid"
	"github.com/meehighlov/eventor/internal/repositories/models"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB().WithContext(ctx).Delete(&models.Event{ID: id}).Error
}
