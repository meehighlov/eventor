package handlers

import (
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/internal/parsers"
	"github.com/meehighlov/eventor/pkg/telegram"
)


func ParseAndBuildEvent(message *telegram.Message) *db.Event {
	notifyAtList := parsers.FindAllTimestampsByMeta(message.Text, "@", parsers.ParseNotifyAtDate)
	scheduleList := parsers.FindAllTimestampsByMeta(message.Text, "&", parsers.ParseScheduleDate)

	notifyAt := ""
	if len(notifyAtList) > 0 {
		notifyAt = notifyAtList[0]
	}

	schedule := ""
	if len(scheduleList) > 0 {
		schedule = scheduleList[0]
	}

	e := db.NewEvent(
		message.From.Id,
		message.GetChatIdStr(),
		message.Text,
		notifyAt,
		schedule,
		"h",
	)

	return e
}
