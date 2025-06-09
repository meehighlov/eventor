package constants

import "github.com/meehighlov/eventor/internal/config"

type Constants struct {
	START_MESSAGE string
	ERROR_MESSAGE string

	BUTTON_TEXT_COPY_EVENT   string
	BUTTON_TEXT_ADD_ELEMENT  string
	BUTTON_TEXT_INFO_ELEMENT string
	BUTTON_TEXT_BACK         string
	BUTTON_TEXT_TUTORIAL     string
	BUTTON_TEXT_EDIT         string
	BUTTON_TEXT_DELETE       string
	BUTTON_TEXT_NAME         string
	BUTTON_TEXT_LINK         string

	COMMAND_START           string
	COMMAND_EVENTS          string
	COMMAND_INFO_EVENT      string
	COMMAND_NEXT_DELTA      string
	COMMAND_DELETE          string
	COMMAND_EDIT_EVENT      string
	COMMAND_EDIT_EVENT_SAVE string
	COMMAND_LIST_EVENT      string

	COMMAND_LIST_ELEMENT          string
	COMMAND_ELEMENT_INFO          string
	COMMAND_ELEMENTS              string
	COMMAND_ADD_ELEMENT           string
	COMMAND_ADD_ELEMENT_SAVE      string
	COMMAND_INFO_ELEMENT          string
	COMMAND_EDIT_ELEMENT          string
	COMMAND_EDIT_ELEMENT_SAVE     string
	COMMAND_DELETE_ELEMENT        string
	COMMAND_ELEMENT_SWITCH_STATUS string
	COMMAND_EDIT_ELEMENT_NAME_SAVE string
	COMMAND_EDIT_ELEMENT_LINK_SAVE string
	COMMAND_EDIT_ELEMENT_REQUEST   string
	COMMAND_DELETE_ELEMENT_CONFIRM string
}

func New(cfg *config.Config) *Constants {
	return &Constants{
		START_MESSAGE:                 "Привет!",
		ERROR_MESSAGE:                 "Произошла ошибка",
		BUTTON_TEXT_COPY_EVENT:        "📋",
		COMMAND_START:                 "/start",
		COMMAND_EVENTS:                "/events",
		COMMAND_INFO_EVENT:            "ei",
		COMMAND_LIST_EVENT:            "event_list",
		COMMAND_NEXT_DELTA:            "next_delta",
		COMMAND_DELETE:                "event_delete",
		COMMAND_EDIT_EVENT:            "event_edit",
		COMMAND_EDIT_EVENT_SAVE:       "event_edit_save",
		COMMAND_LIST_ELEMENT:          "element_list",
		COMMAND_ELEMENT_INFO:          "element_info",
		COMMAND_ELEMENTS:              "/elements",
		BUTTON_TEXT_ADD_ELEMENT:       "➕",
		COMMAND_ADD_ELEMENT:           "element_add",
		COMMAND_ADD_ELEMENT_SAVE:      "element_add_save",
		BUTTON_TEXT_INFO_ELEMENT:      "смотреть",
		COMMAND_INFO_ELEMENT:          "element_info",
		BUTTON_TEXT_BACK:              "⬅️",
		BUTTON_TEXT_TUTORIAL:          "туториал",
		BUTTON_TEXT_EDIT:              "✏️",
		COMMAND_EDIT_ELEMENT:          "element_edit",
		COMMAND_EDIT_ELEMENT_SAVE:     "element_edit_save",
		BUTTON_TEXT_DELETE:            "🗑",
		COMMAND_DELETE_ELEMENT:        "element_delete",
		COMMAND_DELETE_ELEMENT_CONFIRM: "el_del_cnfrm",
		COMMAND_ELEMENT_SWITCH_STATUS: "element_switch_status",
		BUTTON_TEXT_NAME:              "название",
		BUTTON_TEXT_LINK:              "ссылка",
		COMMAND_EDIT_ELEMENT_NAME_SAVE: "el_ed_nm_sv",
		COMMAND_EDIT_ELEMENT_LINK_SAVE: "el_ed_lnk_sv",
		COMMAND_EDIT_ELEMENT_REQUEST:   "el_ed_req",
	}
}
