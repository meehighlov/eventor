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

	if !target.IsScheduled() {
		return []db.Event{}
	}

	conflicts, err := findConflictsByDaysComparison(ctx, target)

	if err != nil {
		slog.Error("error occured while searching for conflicts, when filtering by owner id: " + err.Error())
		return []db.Event{}
	}

	return conflicts
}

func findConflictsByDaysComparison(ctx context.Context, target *db.Event) ([]db.Event, error) {
	related_events, err := (&db.Event{OwnerId: target.OwnerId}).Filter(ctx)
	if err != nil {
		return []db.Event{}, err
	}

	conflicts := []db.Event{}

	targetEventDate, err := target.GetScheduleNearestOrActualDate()
	if err != nil {
		slog.Error("findRelatedEvents", "error parsing target date by target day name", err.Error())
		return conflicts, err
	}

	for _, event_ := range related_events {
		if target.ID == event_.Id() {
			continue
		}

		event := event_.(db.Event)

		eventDate, err := event.GetScheduleNearestOrActualDate()
		if err != nil {
			slog.Error("findRelatedEvents", "error parsing event date by event day name", event.Schedule)
			continue
		}
	
		if eventDate == targetEventDate {
			slog.Debug("findConflictsByDaysComparison", "success comarison for event date", eventDate)
			conflicts = append(conflicts, event)
		}
	}

	return conflicts, nil
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
