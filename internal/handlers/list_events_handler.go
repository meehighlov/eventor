package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/meehighlov/eventor/internal/config"
	"github.com/meehighlov/eventor/internal/db"
	"github.com/meehighlov/eventor/internal/models"
	"github.com/meehighlov/eventor/pkg/telegram"
)

const (
	LIST_PAGINATION_SHIFT = 5
	LIST_LIMIT = 5
	LIST_START_OFFSET = 0

	HEADER_MESSAGE_LIST_NOT_EMPTY = "Нажми, чтобы узнать детали✨"
	HEADER_MESSAGE_LIST_IS_EMPTY = "Записей пока нет✨"
)

func ListEventsHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()

	message := event.GetMessage()
	events, err := (&db.Event{OwnerId: message.From.Id}).Filter(ctx)

	if err != nil {
		slog.Error("Error fetching events" + err.Error())
		return nil
	}

	if len(events) == 0 {
		event.Reply(ctx, "Событий пока нет✨")
		return nil
	}

	eventsAsButtons := buildEventsButtons(events, LIST_LIMIT, LIST_START_OFFSET)
	pagiButtons := buildPagiButtons(len(events), LIST_LIMIT, LIST_START_OFFSET)

	markup := [][]map[string]string{}

	for _, button := range eventsAsButtons {
		markup = append(markup, []map[string]string{button})
	}

	if len(pagiButtons) > 0 {
		markup = append(markup, pagiButtons[0])
	}

	event.ReplyWithKeyboard(ctx, "Нажми на событие, чтобы узнать детали", markup)

	return nil
}


func buildPagiButtons(total, limit, offset int) [][]map[string]string {
	if total == 0 {
		return [][]map[string]string{}
	}
	if offset == total {
		return [][]map[string]string{{
			{
				"text": "свернуть👆",
				"callback_data": models.CallList(strconv.Itoa(LIST_START_OFFSET), "<<<").String(),
			},
		}}
	}
	var keyBoard []map[string]string
	if offset + limit >= total {
		previousButton := map[string]string{"text": "назад", "callback_data": models.CallList(strconv.Itoa(offset), "<<").String()}
		keyBoard = []map[string]string{previousButton}
	} else {
		if offset == 0 {
			nextButton := map[string]string{"text": "вперед", "callback_data": models.CallList(strconv.Itoa(offset), ">>").String()}
			keyBoard = []map[string]string{nextButton}
		} else {
			nextButton := map[string]string{"text": "вперед", "callback_data": models.CallList(strconv.Itoa(offset), ">>").String()}
			previousButton := map[string]string{"text": "назад", "callback_data": models.CallList(strconv.Itoa(offset), "<<").String()}
			keyBoard = []map[string]string{previousButton, nextButton}
		}
	}

	allButton := map[string]string{"text": fmt.Sprintf("показать все (%d)", total), "callback_data": models.CallList(strconv.Itoa(offset), "<>").String()}
	allButtonBar := []map[string]string{allButton}

	markup := [][]map[string]string{}
	if total <= limit {
		return markup
	}

	markup = append(markup, keyBoard)
	markup = append(markup, allButtonBar)

	return markup
}

func ListEventsCallbackQueryHandler(event telegram.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg().HandlerTmeout())
	defer cancel()
	callbackQuery := event.GetCallbackQuery()

	params := models.CallbackFromString(event.GetCallbackQuery().Data)

	offset := params.Pagination.Offset

	limit_ := LIST_LIMIT
	offset_, err := strconv.Atoi(offset)
	if err != nil {
		slog.Error("error parsing offset in list pagination callback query: " + err.Error())
		return err
	}

	events, err := (&db.Event{OwnerId: callbackQuery.From.Id}).Filter(ctx)

	if err != nil {
		slog.Error("Error fetching events: " + err.Error())
		return nil
	}

	direction := params.Pagination.Direction

	slog.Debug(fmt.Sprintf("direction: %s limit: %d offset: %s", direction, limit_, offset))

	if direction == "<" {
		slog.Debug("back to previous screen, offset not changed")
	}
	if direction == "<<<" {
		offset_ = 0
	}
	if direction == ">>" {
		offset_ += LIST_PAGINATION_SHIFT
	} 
	if direction == "<<" {
		offset_ -= LIST_PAGINATION_SHIFT
	}
	if direction == "<>" {
		offset_ = len(events)
	}

	msg := HEADER_MESSAGE_LIST_NOT_EMPTY
	if len(events) == 0 {
		msg = HEADER_MESSAGE_LIST_IS_EMPTY
	}

	event.EditCalbackMessage(ctx, msg, buildEventListMarkup(events, limit_, offset_))

	return nil
}

func buildEventsButtons(events []db.Event, limit, offset int) []map[string]string {
	var buttons []map[string]string
	for i, event := range events {
		if offset != len(events) {
			if i == limit + offset {
				break
			}
			if i < offset {
				continue
			}
		}
		button := map[string]string{
			"text": event.Text,
			"callback_data": models.CallInfo(event.ID, strconv.Itoa(offset)).String(),
		}
		buttons = append(buttons, button)
	}

	return buttons
}

func buildEventListMarkup(friends []db.Event, limit, offset int) [][]map[string]string {
	friendsListAsButtons := buildEventsButtons(friends, limit, offset)
	pagiButtons := buildPagiButtons(len(friends), limit, offset)

	markup := [][]map[string]string{}

	for _, button := range friendsListAsButtons {
		markup = append(markup, []map[string]string{button})
	}

	markup = append(markup, pagiButtons...)

	return markup
}
