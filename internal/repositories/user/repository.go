package user

import (
	"log/slog"

	"github.com/meehighlov/eventor/internal/config"

	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func New(cfg *config.Config, db *gorm.DB, logger *slog.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}
