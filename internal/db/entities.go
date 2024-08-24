package db

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BaseFields struct {
	ID        string
	CreatedAt string
	UpdatedAt string
}

func (b *BaseFields) RefresTimestamps() (created string, updated string, _ error) {
	now := time.Now().Format("02.01.2006T15:04:05")
	if b.CreatedAt == "" {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	return b.CreatedAt, b.UpdatedAt, nil
}

func NewBaseFields() BaseFields {
	now := time.Now().Format("02.01.2006T15:04:05")
	return BaseFields{uuid.New().String(), now, now}
}

type User struct {
	// telegram user -> bot's user

	BaseFields

	TGId       int // telegram user id
	Name       string
	TGusername string
	ChatId     int // chatId - id of chat with user, bot uses it to send notification
}

type Event struct {
	BaseFields

	ChatId   string // chatId - id of chat with user, bot uses it to send notification
	OwnerId  int // telegram user id
	Text     string
	NotifyAt string
	Delta    string
}

func (event *Event) NotifyAtAsTimeObject() (time.Time, error) {
	notifyAt, err := time.Parse("02.01.2006 15:04", event.NotifyAt)
	if err != nil {
		return time.Now(), err
	}

	return notifyAt, err
}

func BuildEvent(ownerId int, chatId, text, timestamp string) (*Event, error) {
	notifyAt, delta, _ := ParseTimesatmp(timestamp)
	e, err := (&Event{
		NewBaseFields(),
		chatId,
		ownerId,
		text,
		notifyAt,
		delta,
	}).Validated()

	if err != nil {
		return nil, err
	}

	return e, nil
}

func ParseTimesatmp(timestamp string) (notifyAt string, delta string, err error) {
	parts := strings.Split(timestamp, " ")
	if len(parts) < 3 {
		return "", "", errors.New("not enough parts for timestamp")
	}

	notifyAt = parts[0] + " " + parts[1]
	delta = parts[2]

	return notifyAt, delta, nil
}

func (e *Event) Validated() (*Event, error) {
	month := "01"
	day := "02"
	format := fmt.Sprintf("%s.%s.2006 15:04", day, month)

	_, err := time.Parse(format, e.NotifyAt)

	if err != nil {
		return nil, err
	}

	_, found := map[string]int{"h": 1, "d": 1, "w": 1, "m": 1, "y": 1}[e.Delta]
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
	default:
		slog.Info("delta of value is not supported, notify date is not changed. Delta value:" + e.Delta)
	}

	e.NotifyAt = notifyAt.Format("02.01.2006 15:04")

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
	default:
		slog.Info("delta of value is not supported, notify date is not changed. Delta value:" + e.Delta)
		return "неизвестный интервал"
	}
}

func (user *User) GetTGUserName() string {
	if !strings.HasPrefix("@", user.TGusername) {
		return "@" + user.TGusername
	}
	return user.TGusername
}
