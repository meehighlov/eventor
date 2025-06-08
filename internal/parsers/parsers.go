package parsers

import (
	"log/slog"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/parsers/ts"
)

type Parsers struct {
	Timestamp *ts.Parser
}

func New(cfg *config.Config, logger *slog.Logger) *Parsers {
	return &Parsers{
		Timestamp: ts.New(cfg, logger),
	}
}
