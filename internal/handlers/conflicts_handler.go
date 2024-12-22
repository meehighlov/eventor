package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/db"
)

func CheckConflictsCallbackHandler(ctx context.Context, event common.Event) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	event.ReplyCallbackQuery(
		ctx,
		buildConflictsMessage(ctx, params.Id, getConflicts(ctx, params.Id)),
	)

	return nil
}

func ConflictsCommandHandler(ctx context.Context, event common.Event) error {
	userEvents, err := db.Event{OwnerId: event.GetMessage().From.Id}.Filter(ctx)
	if err != nil {
		event.Reply(ctx, "Произошла ошибка выгрузки событий: " + err.Error())
	}

	all_conflicts := []db.Event{}
	for _, userEvent := range userEvents {
		conflicts := getConflicts(ctx, userEvent.Id())
		all_conflicts = append(all_conflicts, conflicts...)
	}

	if len(all_conflicts) == 0 {
		event.Reply(ctx, "Конфликтов не обнаружено")
		return nil
	}

	texts := []string{}
	seen := map[string]struct{}{}
	for _, conflict := range all_conflicts {
		if _, found := seen[conflict.ID]; !found {
			texts = append(texts, fmt.Sprintf("🔴 %s", conflict.Text))
		}
		seen[conflict.ID] = struct{}{}
	}

	msg := strings.Join([]string{
		"Обнаружены конфликты в расписании\n",
		strings.Join(texts, "\n\n"),
	}, "\n")

	event.Reply(ctx, msg)

	return nil
}
