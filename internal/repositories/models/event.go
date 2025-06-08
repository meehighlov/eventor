package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ChatID    string    `gorm:"column:chat_id;type:varchar(255);not null"`
	OwnerID   string    `gorm:"column:owner_id;type:varchar(255);not null"`
	Text      string    `gorm:"column:text;type:text;not null"`
	NotifyAt  string    `gorm:"column:notify_at;type:varchar(255);not null"`
	Schedule  string    `gorm:"column:schedule;type:varchar(255);not null"`
	Delta     string    `gorm:"column:delta;type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()"`
}

func (*Event) TableName() string {
	return "event"
}

func (event *Event) NotifyAtAsTimeObject() (time.Time, error) {
	notifyAt, err := time.Parse("02.01 15:04", event.NotifyAt)
	if err != nil {
		return time.Now(), err
	}

	return notifyAt, err
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
		disable = true
	default:
		return "", errors.New("delta of value is not supported, notify date is not changed. Delta value:" + e.Delta)
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
