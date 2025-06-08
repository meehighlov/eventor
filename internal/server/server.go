package server

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/meehighlov/eventor/internal/builders"
	"github.com/meehighlov/eventor/internal/clients"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/constants"
	"github.com/meehighlov/eventor/internal/services"
)

type Server struct {
	logger        *slog.Logger
	handleTimeout time.Duration
	constants     *constants.Constants
	services      *services.Services
	clients       *clients.Clients
	builders      *builders.Builders
	allowedUsers  []string
	webServer     *http.Server
	wg            sync.WaitGroup
	shutdownChan  chan struct{}
	cfg           *config.Config
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	services *services.Services,
	clients *clients.Clients,
	constants *constants.Constants,
	builders *builders.Builders,
) *Server {
	return &Server{
		logger:        logger,
		services:      services,
		clients:       clients,
		constants:     constants,
		builders:      builders,
		handleTimeout: time.Duration(cfg.TelegramHandlerTimeoutSec) * time.Second,
		allowedUsers:  cfg.AllowedUsers(),
		cfg:           cfg,
		shutdownChan:  make(chan struct{}),
		wg:            sync.WaitGroup{},
	}
}
