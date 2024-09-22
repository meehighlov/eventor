package parsers

import (
	"errors"
	"log/slog"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/meehighlov/eventor/internal/config"
)

func FindAllTimestampsByMeta(text, meta string, parser func(string) (string, error)) []string {
	timestamps := searchByTimestampPatterns(text, meta)
	if len(timestamps) == 0 {
		return []string{}
	}

	tsList := []string{}
	for _, ts := range timestamps {
		slog.Debug("FindAllTimestampsByMeta", "raw ts to parse", ts)
		// todo call remove meta symbols
		prepared, err := parser(ts)
		if err != nil {
			slog.Debug("FindAllTimestampsByMeta", "timestamp parser error", err.Error())
		} else {
			tsList = append(tsList, prepared)
		}
	}

	slog.Debug("FindAllTimestampsByMeta", "timestamps to save count", len(tsList))

	return tsList
}

func searchByTimestampPatterns(text, meta string) []string {
	atParserDay, _ := regexp.Compile(meta + `[a-яА-Я]{2} [0-9][0-9]:[0-9][0-9][/s]?`)
	atParserDate, _ := regexp.Compile(meta + `[0-9][0-9].[0-9][0-9] [0-9][0-9]:[0-9][0-9][/s]?`)
	atParserTime, _ := regexp.Compile(meta + `[0-9][0-9]:[0-9][0-9][/s]?`)

	patterns := []regexp.Regexp{*atParserDay, *atParserDate, *atParserTime}

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

	slog.Debug("searchByTimestampPatterns", "hits count (after clean)", len(notifyDates))
	slog.Debug("searchByTimestampPatterns", "matches (after clean)", notifyDates)

	return notifyDates
}

func ParseNotifyAtDate(eventDateRaw string) (string, error) {
	parts := strings.Fields(eventDateRaw)

	if len(parts) == 0 || len(parts) > 2 {
		return "", nil
	}

	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, while extracting notifyAt error: " + err.Error())
		return "", err
	}

	parseTargetNotifyAt := func(notifyAt string) (string, error) {
		layout := "02.01 15:04"
		notifyAtObj, err := time.Parse(layout, notifyAt)
		if err != nil {
			return "", err
		}

		return notifyAtObj.Format(layout), nil
	}

	if len(parts) == 1 {
		slog.Debug("len parts == 1 check time only")
		// suppose time only was sepcified
		// date maybe specified as time only, so prepend day
		notifyAt := time.Now().In(location)
		toValidate := strings.Join([]string{
			notifyAt.Format("02.01"),
			parts[0],
		}, " ")
		return parseTargetNotifyAt(toValidate)
	}

	// check date was specified as <day hh.mm>
	if notifyAt, err := FindNearestDateByDayName(parts[0], false, location); err == nil {
		toValidate := strings.Join([]string{
			notifyAt,
			parts[1],
		}, " ")
		return parseTargetNotifyAt(toValidate)
	}

	return parseTargetNotifyAt(eventDateRaw)
}

func ParseScheduleDate(eventDateRaw string) (string, error) {
	_, err := time.Parse("02.01 15:04", eventDateRaw)
	if err == nil {
		return eventDateRaw, nil
	}

	slog.Debug("ParseScheduleDate", "trying parse by day name", eventDateRaw)

	days_map := map[string]time.Weekday {
		"пн": time.Monday,
		"вт": time.Tuesday,
		"ср": time.Wednesday,
		"чт": time.Thursday,
		"пт": time.Friday,
		"сб": time.Saturday,
		"вс": time.Sunday,
	}

	parts := strings.Fields(eventDateRaw)
	if len(parts) != 2 {
		slog.Error("ParseScheduleDate: not expected parts len")
		return "", errors.New("ParseScheduleDate: not expected parts len")
	}

	day := parts[0]
	if _, found := days_map[day]; !found {
		slog.Error("ParseScheduleDate: not expected day")
		return "", errors.New("ParseScheduleDate: not expected day")
	}

	timeRaw := parts[1]
	_, err = time.Parse("15:04", timeRaw)
	if err != nil {
		slog.Error("ParseScheduleDate", "parse time error", err.Error())
		return "", err
	}

	return strings.Join([]string{day, timeRaw}, " "), nil
}

func FindNearestDateByDayName(dayName string, includeToday bool, location *time.Location) (string, error) {
	days_map := map[string]time.Weekday {
		"пн": time.Monday,
		"вт": time.Tuesday,
		"ср": time.Wednesday,
		"чт": time.Thursday,
		"пт": time.Friday,
		"сб": time.Saturday,
		"вс": time.Sunday,
	}
	now := time.Now().In(location)
	dayNum, found := days_map[dayName]
	if !found {
		return "", errors.New("FindDateByDayName: not found day number by day name: " + dayName)
	}
	diff := int(math.Abs(float64(now.Weekday() - dayNum)))

	slog.Debug("creating notifyat", "day", dayName, "daynum", dayNum)
	slog.Debug("creating notifyat", "days diff", diff)
	slog.Debug("creating notifyat", "now day", now.Day())

	notifyAt := now

	if includeToday {
		if notifyAt.Weekday() == dayNum {
			return notifyAt.Format("02.01"), nil
		}
	}

	for i := 1; i < 8; i ++ {
		notifyAt = notifyAt.AddDate(0, 0, 1)
		if notifyAt.Weekday() == dayNum {
			break	
		}
	}

	return notifyAt.Format("02.01"), nil
}
