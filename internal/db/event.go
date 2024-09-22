package db

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/parsers"
)

type Event struct {
	BaseFields

	// chatId - id of chat with user, bot uses it to send notification
	ChatId string

	// telegram user id
	OwnerId int

	Text     string
	NotifyAt string
	Schedule string
	Delta    string
}

func (e Event) Id() string {
	return e.ID
}

func (e Event) Info() string {
	return e.Text
}

func (e Event) Name() string {
	return "event"
}

func (e Event) Compare() int {
	if e.IsScheduled() {
		return 0
	}
	return 1
}

func (event *Event) NotifyAtAsTimeObject() (time.Time, error) {
	notifyAt, err := time.Parse("02.01 15:04", event.NotifyAt)
	if err != nil {
		return time.Now(), err
	}

	return notifyAt, err
}

func NewEvent(ownerId int, chatId, text, notifyAt, schedule, delta string) *Event {
	e := (&Event{
		NewBaseFields(),
		chatId,
		ownerId,
		text,
		notifyAt,
		schedule,
		delta,
	})

	return e
}

func (e *Event) UpdateNotifyAt() (string, error) {
	notifyAt, err := e.NotifyAtAsTimeObject()
	if err != nil {
		return "", err
	}

	disable := false

	switch e.Delta {
	case "h":
		notifyAt = notifyAt.Add(time.Hour * 1)
	case "d":
		notifyAt = notifyAt.AddDate(0, 0, 1)
	case "w":
		notifyAt = notifyAt.AddDate(0, 0, 7)
	case "m":
		notifyAt = notifyAt.AddDate(0, 1, 0)
	case "y":
		notifyAt = notifyAt.AddDate(1, 0, 0)
	case "0":
		slog.Debug("delta is zero, disabling notification for oneshot event")
		disable = true
	default:
		slog.Info("delta of value is not supported, notify date is not changed. Delta value:" + e.Delta)
	}

	if disable {
		e.NotifyAt = ""
	} else {
		e.NotifyAt = notifyAt.Format("02.01 15:04")
	}

	return e.NotifyAt, nil
}

func deltaReadable(delta string) string {
	switch delta {
	case "h":
		return "раз в час"
	case "d":
		return "раз в день"
	case "w":
		return "раз в неделю"
	case "m":
		return "раз в месяц"
	case "y":
		return "раз в год"
	case "0":
		return "без повторений"
	default:
		slog.Info("delta of value is not supported, notify date is not changed. Delta value:" + delta)
		return "неизвестный интервал"
	}
}

func (e *Event) DeltaReadable() string {
	return deltaReadable(e.Delta)
}

func (e *Event) NextDelta(asReadable bool) string {
	default_ := "0"
	currentToNext := map[string]string{
		"h": "d",
		"d": "w",
		"w": "m",
		"m": "y",
		"y": "0",
		"0": "h",
	}
	next, found := currentToNext[e.Delta]
	if !found {
		slog.Error("not next delta by current: " + e.Delta + " returning " + default_)
		if asReadable {
			return deltaReadable(default_)
		}
		return default_
	}

	if asReadable {
		return deltaReadable(next)
	}
	return next
}

func (e *Event) NotifyNeeded() bool {
	return e.NotifyAt != "" && e.Delta != ""
}

func (e *Event) IsScheduled() bool {
	return e.Schedule != ""
}

func (e *Event) GetScheduleNearestOrActualDate() (string, error) {
	if !e.IsScheduled() {
		return "", errors.New("event is not schedule, event id: " + e.ID)
	}
	if _, err := time.Parse("02.01 15:04", e.Schedule); err == nil {
		return strings.Fields(e.Schedule)[0], nil
	}

	day := strings.Fields(e.Schedule)[0]
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		return "", err
	}
	return parsers.FindNearestDateByDayName(day, true, location)
}
