package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/eventor/internal/db"
)

func getConflicts(ctx context.Context, scheduleId string) []db.Event {
	target := getTarget(ctx, scheduleId)
	if target == nil {
		return []db.Event{}
	}

	related_events, err := (&db.Event{OwnerId: target.OwnerId, Schedule: target.Schedule}).Filter(ctx)
	if err != nil {
		slog.Error("error occured while searching for conflicts, when filtering by owner id: " + err.Error())
		return []db.Event{}
	}

	conflicts := []db.Event{}
	for _, event_ := range related_events {
		if target.ID == event_.Id() {
			continue
		}
		conflicts = append(conflicts, event_.(db.Event))
	}

	return conflicts
}

func buildConflictsMessage(ctx context.Context, targetId string, conflicts []db.Event) string {
	if len(conflicts) == 0 {
		return "Конфликтов не обнаружено"
	}

	target := getTarget(ctx, targetId)
	if target == nil {
		return "Ошибка поиска целевого расписания"
	}

	metas := []string{
		"☝️ Обнаружены конфликты в расписании",
		fmt.Sprintf("🗓 %s", target.Schedule),
	}

	for _, c := range conflicts {
		metas = append(metas, "🔴 " + c.Text)
	}

	msg := strings.Join(metas, "\n\n")

	return msg
}

func getTarget(ctx context.Context, scheduleId string) *db.Event {
	baseFields := db.BaseFields{ID: scheduleId}
	scs, err := (&db.Event{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		slog.Error("error occured while searching for conflicts: " + err.Error())
		return nil
	}

	if len(scs) == 0 {
		slog.Error("no schedules for id: " + scheduleId)
		return nil
	}

	target := scs[0].(db.Event)

	return &target
}
