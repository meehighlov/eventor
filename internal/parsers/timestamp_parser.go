package parsers

import (
	"log/slog"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/meehighlov/eventor/internal/config"
)


func FindAllTimestampsByMeta(text, meta string) []string {
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

	slog.Debug("findAllNotifyAt", "hits count (after clean)", len(notifyDates))
	slog.Debug("findAllNotifyAt", "matches (after clean)", notifyDates)

	return notifyDates
}

func parseEventDate(eventDateRaw string) string {
	parts := strings.Split(eventDateRaw, " ")

	if len(parts) == 0 || len(parts) > 2{
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

	parseTargetNotifyAt := func(notifyAt string) string {
		layout := "02.01 15:04"
		notifyAtObj, err := time.Parse(layout, notifyAt)
		if err != nil {
			return ""
		}

		return notifyAtObj.Format(layout)
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

	day := parts[0]

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