package db

import "log/slog"

type Schedule struct {
	BaseFields

	// chatId - id of chat with user, bot uses it to send notification
	ChatId string

	// telegram user id
	OwnerId int

	// pyload
	Text string

	// period
	Delta string

	// day when schedule is originate
	Day string

	// event id by which notifications are sent
	EventId string
}

func NewSchedule(ownerId int, chatId, text, delta, day string) *Schedule {
	return &Schedule{
		NewBaseFields(),
		chatId,
		ownerId,
		text,
		delta,
		day,
		"notset",
	}
}

func (s Schedule) Id() string {
	return s.ID
}

func (s Schedule) Info() string {
	return s.Text
}

func (Schedule) Name() string {
	return "schedule"
}

func (Schedule) Compare() int {
	return -1
}

func (s *Schedule) DeltaReadable() string {
	switch s.Delta {
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
		slog.Info("delta of value is not supported, notify date is not changed. Delta value:" + s.Delta)
		return "неизвестный интервал"
	}
}

func (s *Schedule) HasNotifications() bool {
	return s.EventId != "notset"
}
