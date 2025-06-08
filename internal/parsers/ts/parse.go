package ts

import (
	"errors"
	"strings"
	"time"
)

func (p *Parser) ParseNotifyAtDate(eventDateRaw string) (string, error) {
	parts := strings.Fields(eventDateRaw)

	if len(parts) == 0 || len(parts) > 2 {
		return "", nil
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
		p.logger.Debug("len parts == 1 check time only")
		// suppose time only was sepcified
		// date maybe specified as time only, so prepend day
		notifyAt := time.Now().In(p.location)
		toValidate := strings.Join([]string{
			notifyAt.Format("02.01"),
			parts[0],
		}, " ")
		return parseTargetNotifyAt(toValidate)
	}

	// check date was specified as <day hh.mm>
	if notifyAt, err := p.FindNearestDateByDayName(parts[0], false); err == nil {
		toValidate := strings.Join([]string{
			notifyAt,
			parts[1],
		}, " ")
		return parseTargetNotifyAt(toValidate)
	}

	return parseTargetNotifyAt(eventDateRaw)
}

func (p *Parser) ParseScheduleDate(eventDateRaw string) (string, error) {
	_, err := time.Parse("02.01 15:04", eventDateRaw)
	if err == nil {
		return eventDateRaw, nil
	}

	p.logger.Debug("ParseScheduleDate", "trying parse by day name", eventDateRaw)

	days_map := map[string]time.Weekday{
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
		p.logger.Error("ParseScheduleDate: not expected parts len")
		return "", errors.New("ParseScheduleDate: not expected parts len")
	}

	day := parts[0]
	if _, found := days_map[day]; !found {
		p.logger.Error("ParseScheduleDate: not expected day")
		return "", errors.New("ParseScheduleDate: not expected day")
	}

	timeRaw := parts[1]
	_, err = time.Parse("15:04", timeRaw)
	if err != nil {
		p.logger.Error("ParseScheduleDate", "parse time error", err.Error())
		return "", err
	}

	return strings.Join([]string{day, timeRaw}, " "), nil
}
