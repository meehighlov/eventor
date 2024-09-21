package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/meehighlov/eventor/internal/common"
	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
)

func CheckConflictsCallbackHandler(event common.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	event.ReplyCallbackQuery(
		ctx,
		buildConflictsMessage(ctx, params.Id, getConflicts(ctx, params.Id)),
	)

	return nil
}

func ConflictsCommandHandler(event common.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	userEvents, err := db.Event{OwnerId: event.GetMessage().From.Id}.Filter(ctx)
	if err != nil {
		event.Reply(ctx, "Произошла ошибка выгрузки событий: " + err.Error())
	}

	all_conflicts := []db.Event{}
	seen := map[string]interface{}{}
	for _, userEvent := range userEvents {
		if _, found := seen[userEvent.Id()]; !found {
			conflicts := getConflicts(ctx, userEvent.Id())
			all_conflicts = append(all_conflicts, conflicts...)
			for _, c := range conflicts {
				seen[c.ID] = 1
			}
		}
		seen[userEvent.Id()] = 1
	}

	if len(all_conflicts) == 0 {
		event.Reply(ctx, "Конфликтов нет")
		return nil
	}

	texts := []string{}
	for _, conflict := range all_conflicts {
		texts = append(texts, fmt.Sprintf("🔴 %s", conflict.Text))
	}

	msg := strings.Join([]string{
		"Обнаружены конфликты в расписании\n",
		strings.Join(texts, "\n"),
	}, "\n")

	event.Reply(ctx, msg)

	return nil
}
