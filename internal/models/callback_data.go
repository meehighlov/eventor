package models

import (
	"strings"
)

type pagination struct {
	Offset    string
	Direction string
}

type CallbackDataModel struct {
	Command    string
	Id         string
	Pagination pagination
	Entity     string
}

func CallList(offset, direction string) *CallbackDataModel {
	return newCallback("list", "", offset, direction, "event")
}

func CallDelete(id string) *CallbackDataModel {
	return newCallback("delete", id, "", "", "event")
}

func CallInfo(id, offset string) *CallbackDataModel {
	return newCallback("info", id, offset, "", "event")
}

func CallConflicts(id string) *CallbackDataModel {
	return newCallback("conflicts", id, "0", "", "event")
}

func newCallback(command, id, offset, direction, entity string) *CallbackDataModel {
	return &CallbackDataModel{
		Command: command,
		Id: id,
		Pagination: pagination{
			Offset: offset,
			Direction: direction,
		},
		Entity: entity,
	}
}

func CallbackFromString(raw string) *CallbackDataModel {
	params := strings.Split(raw, ";")
	return &CallbackDataModel{
		Command: params[0],
		Id: params[1],
		Pagination: pagination{
			Offset: params[2],
			Direction: params[3],
		},
		Entity: params[4],
	}
}

func (cd *CallbackDataModel) String() string {
	separator := ";"
	return strings.Join(
		[]string{
			cd.Command,
			cd.Id,
			cd.Pagination.Offset,
			cd.Pagination.Direction,
			cd.Entity,
		},
		separator,
	)
}
