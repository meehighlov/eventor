package ts

import (
	"errors"
	"math"
	"regexp"
	"strings"
	"time"
)

func (p *Parser) searchByTimestampPatterns(text, meta string) []string {
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

	p.logger.Debug("searchByTimestampPatterns", "hits count (after clean)", len(notifyDates))
	p.logger.Debug("searchByTimestampPatterns", "matches (after clean)", notifyDates)

	return notifyDates
}

func (p *Parser) FindNearestDateByDayName(dayName string, includeToday bool) (string, error) {
	days_map := map[string]time.Weekday{
		"пн": time.Monday,
		"вт": time.Tuesday,
		"ср": time.Wednesday,
		"чт": time.Thursday,
		"пт": time.Friday,
		"сб": time.Saturday,
		"вс": time.Sunday,
	}
	now := time.Now().In(p.location)
	dayNum, found := days_map[dayName]
	if !found {
		return "", errors.New("FindDateByDayName: not found day number by day name: " + dayName)
	}
	diff := int(math.Abs(float64(now.Weekday() - dayNum)))

	p.logger.Debug("creating notifyat", "day", dayName, "daynum", dayNum)
	p.logger.Debug("creating notifyat", "days diff", diff)
	p.logger.Debug("creating notifyat", "now day", now.Day())

	notifyAt := now

	if includeToday {
		if notifyAt.Weekday() == dayNum {
			return notifyAt.Format("02.01"), nil
		}
	}

	for i := 1; i < 8; i++ {
		notifyAt = notifyAt.AddDate(0, 0, 1)
		if notifyAt.Weekday() == dayNum {
			break
		}
	}

	return notifyAt.Format("02.01"), nil
}

func (p *Parser) FindAllTimestampsByMeta(text, meta string, parser func(string) (string, error)) []string {
	timestamps := p.searchByTimestampPatterns(text, meta)
	if len(timestamps) == 0 {
		return []string{}
	}

	tsList := []string{}
	for _, ts := range timestamps {
		p.logger.Debug("FindAllTimestampsByMeta", "raw ts to parse", ts)
		// todo call remove meta symbols
		prepared, err := parser(ts)
		if err != nil {
			p.logger.Debug("FindAllTimestampsByMeta", "timestamp parser error", err.Error())
		} else {
			tsList = append(tsList, prepared)
		}
	}

	p.logger.Debug("FindAllTimestampsByMeta", "timestamps to save count", len(tsList))

	return tsList
}
