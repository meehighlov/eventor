package repositories

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/repositories/event"
	"github.com/meehighlov/eventor/internal/repositories/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Repositories struct {
	User  *user.Repository
	Event *event.Repository
}

func New(cfg *config.Config, logger *slog.Logger) *Repositories {
	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{
		Logger: WrapAppLogger(logger),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := RunMigrations(context.Background(), cfg, logger, db); err != nil {
		log.Fatal("Migration error:", err)
	}

	return &Repositories{
		User:  user.New(cfg, db, logger),
		Event: event.New(cfg, db, logger),
	}
}
