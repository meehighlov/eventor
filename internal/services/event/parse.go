package event

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meehighlov/eventor/internal/clients/telegram"
	"github.com/meehighlov/eventor/internal/repositories/models"
)

func (s *Service) ParseAndBuildEvent(message *telegram.Message) (*models.Event, error) {
	notifyAtList := s.parsers.Timestamp.FindAllTimestampsByMeta(message.Text, "@", s.parsers.Timestamp.ParseNotifyAtDate)
	scheduleList := s.parsers.Timestamp.FindAllTimestampsByMeta(message.Text, "&", s.parsers.Timestamp.ParseScheduleDate)

	notifyAt := ""
	if len(notifyAtList) > 0 {
		notifyAt = notifyAtList[0]
	}

	schedule := ""
	if len(scheduleList) > 0 {
		schedule = scheduleList[0]
	}

	if len(notifyAt) == 0 && len(schedule) == 0 {
		return nil, errors.New("no notify at or schedule found")
	}

	e := models.Event{
		ID:       uuid.New(),
		OwnerID:  strconv.Itoa(message.From.Id),
		ChatID:   message.GetChatIdStr(),
		Text:     message.Text,
		NotifyAt: notifyAt,
		Schedule: schedule,
		Delta:    "h",
	}

	return &e, nil
}

func (s *Service) GetScheduleNearestOrActualDate(e *models.Event) (string, error) {
	if !e.IsScheduled() {
		return "", errors.New("event is not schedule, event id: " + e.ID.String())
	}
	if _, err := time.Parse("02.01 15:04", e.Schedule); err == nil {
		return strings.Fields(e.Schedule)[0], nil
	}

	day := strings.Fields(e.Schedule)[0]
	return s.parsers.Timestamp.FindNearestDateByDayName(day, true)
}
