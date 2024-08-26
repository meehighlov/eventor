package db

import (
	"errors"
	"fmt"
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

func BuildEvent(ownerId int, chatId, text, timestamp, delta string) (*Event, error) {
	e, err := (&Event{
		NewBaseFields(),
		chatId,
		ownerId,
		text,
		timestamp,
		delta,
	}).build()

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Event) build() (*Event, error) {
	month := "01"
	day := "02"
	format := fmt.Sprintf("%s.%s 15:04", day, month)

	_, err := time.Parse(format, e.NotifyAt)

	if err != nil {
		slog.Error("build event error: " + err.Error())
		return nil, err
	}

	// additionaly validate
	_, found := map[string]int{"0": 1, "h": 1, "d": 1, "w": 1, "m": 1, "y": 1}[e.Delta]
	if !found {
		return nil, errors.New("delta format is incorrect")
	}

	return e, nil
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

func (e *Event) DeltaReadable() string {
	switch e.Delta {
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
		slog.Info("delta of value is not supported, notify date is not changed. Delta value:" + e.Delta)
		return "неизвестный интервал"
	}
}
