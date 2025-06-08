package user

import (
	"log/slog"

	"github.com/meehighlov/eventor/internal/builders"
	"github.com/meehighlov/eventor/internal/clients"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/constants"
	"github.com/meehighlov/eventor/internal/repositories"
	"github.com/meehighlov/eventor/internal/validators"
)

type Service struct {
	logger       *slog.Logger
	repositories *repositories.Repositories
	clients      *clients.Clients
	builders     *builders.Builders
	validators   *validators.Validators
	constants    *constants.Constants
}

func New(cfg *config.Config, logger *slog.Logger, repositories *repositories.Repositories, clients *clients.Clients, builders *builders.Builders, validators *validators.Validators, constants *constants.Constants) *Service {
	return &Service{
		logger:       logger,
		repositories: repositories,
		clients:      clients,
		builders:     builders,
		validators:   validators,
		constants:    constants,
	}
}
