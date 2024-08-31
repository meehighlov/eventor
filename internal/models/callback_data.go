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

func CallList(offset, direction, entity string) *CallbackDataModel {
	return newCallback("list", "", offset, direction, entity)
}

func CallDelete(id, entity string) *CallbackDataModel {
	return newCallback("delete", id, "", "", "event")
}

func CallInfo(id, offset, entity string) *CallbackDataModel {
	return newCallback("info_"+entity, id, offset, "", entity)
}

func CallNextDelta(id, offset string) *CallbackDataModel {
	return newCallback("next_delta", id, offset, "", "event")
}

func CallEdit(id string) *CallbackDataModel {
	return newCallback("edit", id, "", "", "event")
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
