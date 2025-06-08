package watcher

import (
	"log/slog"
	"time"

	"github.com/meehighlov/eventor/internal/builders"
	"github.com/meehighlov/eventor/internal/clients"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/constants"
	"github.com/meehighlov/eventor/internal/repositories"
	"github.com/meehighlov/eventor/internal/validators"
)

type Service struct {
	logger           *slog.Logger
	repositories     *repositories.Repositories
	clients          *clients.Clients
	builders         *builders.Builders
	validators       *validators.Validators
	constants        *constants.Constants
	CheckIntervalSec int
	location         *time.Location
	ReportChatId     string
}

func New(cfg *config.Config, logger *slog.Logger, repositories *repositories.Repositories, clients *clients.Clients, builders *builders.Builders, validators *validators.Validators, constants *constants.Constants) *Service {
	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		panic(err.Error())
	}

	return &Service{
		logger:           logger,
		repositories:     repositories,
		clients:          clients,
		builders:         builders,
		validators:       validators,
		constants:        constants,
		CheckIntervalSec: cfg.WatcherCheckIntervalSec,
		location:         location,
		ReportChatId:     cfg.ReportChatId,
	}
}
