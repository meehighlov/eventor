package handlers

import (
	"context"
	"fmt"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/db"
)

func StartHandler(ctx context.Context, event common.Event) error {
	message := event.GetMessage()

	user := db.User{
		BaseFields: db.NewBaseFields(),
		Name:       message.From.FirstName,
		TGusername: message.From.Username,
		TGId:       message.From.Id,
		ChatId:     message.Chat.Id,
	}

	user.Save(ctx)

	hello := fmt.Sprintf(
		"Привет, %s 👋",
		message.From.Username,
	)

	event.Reply(ctx, hello)

	return nil
}
