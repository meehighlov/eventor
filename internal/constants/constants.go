package constants

import "github.com/meehighlov/eventor/internal/config"

type Constants struct {
	START_MESSAGE string
	ERROR_MESSAGE string

	BUTTON_TEXT_COPY_EVENT string

	COMMAND_START           string
	COMMAND_EVENTS          string
	COMMAND_INFO_EVENT      string
	COMMAND_NEXT_DELTA      string
	COMMAND_DELETE          string
	COMMAND_EDIT_EVENT      string
	COMMAND_EDIT_EVENT_SAVE string
	COMMAND_LIST_EVENT      string
}

func New(cfg *config.Config) *Constants {
	return &Constants{
		START_MESSAGE:           "Привет!",
		ERROR_MESSAGE:           "Произошла ошибка",
		BUTTON_TEXT_COPY_EVENT:  "📋",
		COMMAND_START:           "/start",
		COMMAND_EVENTS:          "/events",
		COMMAND_INFO_EVENT:      "ei",
		COMMAND_LIST_EVENT:      "event_list",
		COMMAND_NEXT_DELTA:      "next_delta",
		COMMAND_DELETE:          "event_delete",
		COMMAND_EDIT_EVENT:      "event_edit",
		COMMAND_EDIT_EVENT_SAVE: "event_edit_save",
	}
}
