package builders

import (
	"log/slog"

	"github.com/meehighlov/eventor/internal/builders/callback_data"
	"github.com/meehighlov/eventor/internal/builders/keyboard"
	"github.com/meehighlov/eventor/internal/builders/short_id"
	"github.com/meehighlov/eventor/internal/config"
)

type Builders struct {
	ShortIdBuilder      *short_id.Builder
	CallbackDataBuilder *callbackdata.Builder
	KeyboardBuilder     *keyboard.Builder
}

func New(cfg *config.Config, logger *slog.Logger) *Builders {
	return &Builders{
		ShortIdBuilder:      short_id.New(6),
		CallbackDataBuilder: callbackdata.New(),
		KeyboardBuilder:     keyboard.New(),
	}
}
