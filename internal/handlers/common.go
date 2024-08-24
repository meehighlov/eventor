package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/eventor/internal/db"
)

// searches all conflicts in owner's events
// NOTE: searh by notify date
// accuracy - one day
// if found events at the same day - then those events will be return
func getConflicts(ctx context.Context, eventId string) []db.Event {
	baseFields := db.BaseFields{ID: eventId}
	events, err := (&db.Event{BaseFields: baseFields}).Filter(ctx)

	if err != nil {
		slog.Error("error occured while searching for conflicts: " + err.Error())
		return []db.Event{}
	}

	if len(events) == 0 {
		slog.Error("no events for id: " + eventId)
		return []db.Event{}
	}

	target := events[0]

	related_events, err := (&db.Event{OwnerId: target.OwnerId}).Filter(ctx)
	if err != nil {
		slog.Error("error occured while searching for conflicts, when filtering by owner id: " + err.Error())
		return []db.Event{}
	}

	conflicts := []db.Event{}
	for _, event_ := range related_events {
		if target.ID == event_.ID {
			continue
		}
		if target.CountDaysToBegin() == event_.CountDaysToBegin() {
			conflicts = append(conflicts, event_)
		}
	}

	return conflicts
}

func buildConflictsMessage(targetId string, conflicts []db.Event) string {
	if len(conflicts) == 0 {
		return "Конфликтов не обнаружено"
	}

	var target db.Event
	for _, target_ := range conflicts {
		if target_.ID == targetId {
			target = target_
			break
		}
	}

	metas := []string{
		"Обнаружены конфликты☝️", "\n",
		fmt.Sprintf("Событие %s дата %s", target.Text, target.NotifyAt), "\n",
		"Важно: конфликт обнаружен в дате уведомления",
		"Обрати внимание на дату события", "\n",
	}

	for _, c := range conflicts {
		metas = append(metas, c.Text + " " + c.NotifyAt)
	}

	msg := strings.Join(metas, "\n")

	return msg
}
