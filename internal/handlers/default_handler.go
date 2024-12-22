package handlers

import (
	"context"
	"strings"

	"github.com/meehighlov/eventor/internal/common"
)

func DefaultTextHandler(ctx context.Context, event common.Event) error {
	metas := []string{"@", "&"}

	hasMeta := false
	for _, sym := range metas {
		if strings.Contains(event.GetMessage().Text, sym) {
			hasMeta = true
		}
	}

	if hasMeta {
		return AddEventSave(ctx, event)
	}

	return nil
}
