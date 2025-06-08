package ts

import (
	"log/slog"
	"time"

	"github.com/meehighlov/eventor/internal/config"
)

type Parser struct {
	location *time.Location
	logger   *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Parser {
	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		panic(err)
	}

	return &Parser{
		location: location,
		logger:   logger,
	}
}
