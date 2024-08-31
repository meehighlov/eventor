package db

import (
	"log/slog"
	"time"

	"github.com/meehighlov/eventor/internal/config"
)

type Event struct {
	BaseFields

	// chatId - id of chat with user, bot uses it to send notification
	ChatId string

	// telegram user id
	OwnerId int

	Text     string
	NotifyAt string
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
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " eventId: " + e.ID)
	}
	now := time.Now().In(location)
	notify, err := time.Parse("02.01 15:04", e.NotifyAt)
	if err != nil {
		slog.Error("error parsing notify during count days to event begining: " + err.Error())
		return -1
	}

	diff := now.Sub(notify)
	diff_days := diff.Hours() / 24

	return int(diff_days)
}

func (event *Event) NotifyAtAsTimeObject() (time.Time, error) {
	notifyAt, err := time.Parse("02.01 15:04", event.NotifyAt)
	if err != nil {
		return time.Now(), err
	}

	return notifyAt, err
}

func NewEvent(ownerId int, chatId, text, notifyAt, delta string) *Event {
	e := (&Event{
		NewBaseFields(),
		chatId,
		ownerId,
		text,
		notifyAt,
		"0",
	})

	return e
}

func (e *Event) UpdateNotifyAt() (string, error) {
	notifyAt, err := e.NotifyAtAsTimeObject()
	if err != nil {
		return "", err
	}

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
		slog.Debug("delta is zero, notify timestampt not changed")
	default:
		slog.Info("delta of value is not supported, notify date is not changed. Delta value:" + e.Delta)
	}

	e.NotifyAt = notifyAt.Format("02.01 15:04")

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
