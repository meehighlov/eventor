package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/eventor/internal/db"
)

func getConflicts(ctx context.Context, scheduleId string) []db.Schedule {
	baseFields := db.BaseFields{ID: scheduleId}
	scs, err := (&db.Schedule{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		slog.Error("error occured while searching for conflicts: " + err.Error())
		return []db.Schedule{}
	}

	if len(scs) == 0 {
		slog.Error("no schedules for id: " + scheduleId)
		return []db.Schedule{}
	}

	target := scs[0]

	related_events, err := (&db.Schedule{OwnerId: target.OwnerId}).Filter(ctx)
	if err != nil {
		slog.Error("error occured while searching for conflicts, when filtering by owner id: " + err.Error())
		return []db.Schedule{}
	}

	conflicts := []db.Schedule{}
	for _, event_ := range related_events {
		if target.ID == event_.ID {
			continue
		}
		if target.Day == event_.Day {
			conflicts = append(conflicts, event_)
		}
	}

	return conflicts
}

func buildConflictsMessage(targetId string, conflicts []db.Schedule) string {
	if len(conflicts) == 0 {
		return "Конфликтов не обнаружено"
	}

	var target db.Schedule
	for _, target_ := range conflicts {
		if target_.ID == targetId {
			target = target_
			break
		}
	}

	metas := []string{
		"Обнаружены конфликты в расписании☝️", "\n",
		fmt.Sprintf("Событие %s дата %s", target.Text, target.Day), "\n",
	}

	for _, c := range conflicts {
		metas = append(metas, c.Text + " " + c.Day)
	}

	msg := strings.Join(metas, "\n")

	return msg
}
