package handlers

import (
	"context"
	"log/slog"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/pkg/telegram"
)

func addEventEntry(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	msg := []string{
		"Введи описание\n",
	}

	event.Reply(ctx, strings.Join(msg, ""))

	return 2, nil
}

func addEventSave(event telegram.Event) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	message := event.GetMessage()

	notifyAtList := findAllTimestampsByMeta(message.Text, "@")
	scheduleList := findAllTimestampsByMeta(message.Text, "#")

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
		"0",
	)

	e.Save(ctx)

	msg := "Событие сохранено"
	event.Reply(ctx, msg)

	return -1, nil
}

func findAllTimestampsByMeta(text, meta string) []string {
	notifyDates := searchByNotifyAtPatterns(text, meta)
	if len(notifyDates) == 0 {
		return []string{}
	}

	notifyAtList := []string{}
	for _, notifyAtRaw := range notifyDates {
		// todo call remove meta symbols
		notifyAt := parseEventDate(notifyAtRaw)
		if notifyAt != "" {
			notifyAtList = append(notifyAtList, notifyAt)
		}
	}

	slog.Debug("findAllNotifyAt", "notifyAt to be set", len(notifyAtList))

	return notifyAtList
}

func searchByNotifyAtPatterns(text, meta string) []string {
	atParserDay, _ := regexp.Compile(meta + `[a-яА-Я]{2} [0-9][0-9]:[0-9][0-9][/s]?`)
	atParserDate, _ := regexp.Compile(meta + `[0-9][0-9].[0-9][0-9] [0-9][0-9]:[0-9][0-9][/s]?`)

	patterns := []regexp.Regexp{*atParserDay, *atParserDate}

	clean := func(hits []string) []string {
		cleaned := []string{}
		for _, hit := range hits {
			cleaned_ := strings.TrimSpace(strings.Replace(hit, meta, "", 1))
			cleaned = append(cleaned, cleaned_)
		}
		return cleaned
	}

	notifyDates := []string{}
	for _, p := range patterns {
		hits := p.FindAllString(text, -1)

		notifyDates = append(notifyDates, clean(hits)...)
	}

	slog.Debug("findAllNotifyAt", "hits count (after clean)", len(notifyDates))
	slog.Debug("findAllNotifyAt", "matches (after clean)", notifyDates)

	return notifyDates
}

func parseEventDate(eventDateRaw string) string {
	parts := strings.Split(eventDateRaw, " ")

	if len(parts) != 2 {
		return ""
	}

	days_map := map[string]time.Weekday {
		"пн": time.Monday,
		"вт": time.Tuesday,
		"ср": time.Wednesday,
		"чт": time.Thursday,
		"пт": time.Friday,
		"сб": time.Saturday,
		"вс": time.Sunday,
	}

	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, while extracting notifyAt error: " + err.Error())
	}

	day := parts[0]

	parseTargetNotifyAt := func(notifyAt string) string {
		layout := "02.01 15:04"
		notifyAtObj, err := time.Parse(layout, notifyAt)
		if err != nil {
			return ""
		}

		return notifyAtObj.Format(layout)
	}

	// check date was specified as <day hh.mm>
	dayNum, found := days_map[day]
	if found {
		now := time.Now().In(location)
		diff := int(math.Abs(float64(now.Weekday() - dayNum)))

		slog.Debug("creating notifyat", "day", day, "daynum", dayNum)
		slog.Debug("creating notifyat", "days diff", diff)
		slog.Debug("creating notifyat", "now day", now.Day())

		notifyAt := now
		for i := 1; i < 8; i ++ {
			notifyAt = notifyAt.AddDate(0, 0, 1)
			if notifyAt.Weekday() == dayNum {
				break
			}
		}

		toValidate := strings.Join([]string{
			notifyAt.Format("02.01"),
			parts[1],
		}, " ")

		return parseTargetNotifyAt(toValidate)
	}

	return parseTargetNotifyAt(eventDateRaw)
}

func AddEventHandler() map[int]telegram.CommandStepHandler {
	return map[int]telegram.CommandStepHandler{
		1: addEventEntry,
		2: addEventSave,
	}
}
