package services

import (
	"log/slog"

	"github.com/meehighlov/eventor/internal/builders"
	"github.com/meehighlov/eventor/internal/clients"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/constants"
	"github.com/meehighlov/eventor/internal/parsers"
	"github.com/meehighlov/eventor/internal/repositories"
	"github.com/meehighlov/eventor/internal/services/event"
	"github.com/meehighlov/eventor/internal/services/user"
	"github.com/meehighlov/eventor/internal/services/watcher"
	"github.com/meehighlov/eventor/internal/validators"
)

type Services struct {
	User    *user.Service
	Watcher *watcher.Service
	Event   *event.Service
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	repositories *repositories.Repositories,
	clients *clients.Clients,
	builders *builders.Builders,
	validators *validators.Validators,
	constants *constants.Constants,
	parsers *parsers.Parsers,
) *Services {
	return &Services{
		User:    user.New(cfg, logger, repositories, clients, builders, validators, constants),
		Watcher: watcher.New(cfg, logger, repositories, clients, builders, validators, constants),
		Event:   event.New(cfg, logger, repositories, clients, builders, validators, constants, parsers),
	}
}
