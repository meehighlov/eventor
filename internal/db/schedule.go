package db

import (
	"errors"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/meehighlov/eventor/internal/config"
)

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

	// time when schedule event is started
	Timestamp string

	// event id by which notifications are sent
	EventId string
}

func BuildSchedule(ownerId int, chatId, text, timestamp string) (*Schedule, error) {
	parts := strings.Split(timestamp, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid timestamp value for schedule")
	}
	_, found := map[string]time.Weekday{
		"пн": time.Monday,
		"вт": time.Tuesday,
		"ср": time.Wednesday,
		"чт": time.Thursday,
		"пт": time.Friday,
		"сб": time.Saturday,
		"вс": time.Sunday,
	}[parts[0]]

	if !found {
		return nil, errors.New("unkwon day value")
	}

	_, found = map[string]bool {
		"h": true, "d": true, "w": true, "m": true, "y": true, "0": true,
	}[parts[2]]

	if !found {
		return nil, errors.New("unkwon delta value")
	}

	return &Schedule{
		NewBaseFields(),
		chatId,
		ownerId,
		text,
		parts[2],
		parts[0],
		parts[1],
		"notset",
	}, nil
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

func (s Schedule) UnboundEvent() string {
	value := "notset"
	s.EventId = value
	return s.EventId
}

func (s *Schedule) DayNum() time.Weekday {
	days_map := map[string]time.Weekday {
		"пн": time.Monday,
		"вт": time.Tuesday,
		"ср": time.Wednesday,
		"чт": time.Thursday,
		"пт": time.Friday,
		"сб": time.Saturday,
		"вс": time.Sunday,
	}
	return days_map[s.Day]
}

func (s *Schedule) BuildNotifyAt() string {
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " scheduleId: " + s.ID)
	}
	now := time.Now().In(location)

	diff := int(math.Abs(float64(now.Weekday() - s.DayNum())))

	notifyAt := now.AddDate(0, 0, diff)

	parts := strings.Split(notifyAt.Format("02.01 15:04"), " ")

	return parts[0] + " " + s.Timestamp
}
