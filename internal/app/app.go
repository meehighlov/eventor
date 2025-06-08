package app

import (
	"github.com/meehighlov/eventor/internal/builders"
	"github.com/meehighlov/eventor/internal/clients"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/constants"
	"github.com/meehighlov/eventor/internal/parsers"
	"github.com/meehighlov/eventor/internal/repositories"
	"github.com/meehighlov/eventor/internal/server"
	"github.com/meehighlov/eventor/internal/services"
	"github.com/meehighlov/eventor/internal/validators"
)

func Run() {
	cfg := config.MustLoad()
	logger := MustSetupLogging(cfg)

	repositories := repositories.New(cfg, logger)
	clients := clients.New(cfg, logger)
	builders := builders.New(cfg, logger)
	validators := validators.New(cfg, logger)
	constants := constants.New(cfg)
	parsers := parsers.New(cfg, logger)
	services := services.New(cfg, logger, repositories, clients, builders, validators, constants, parsers)

	server := server.New(cfg, logger, services, clients, constants, builders)
	server.Watch()
	server.Serve()
}
