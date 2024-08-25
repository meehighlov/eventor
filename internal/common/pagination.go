package common

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/meehighlov/eventor/internal/models"
)

const (
	LIST_PAGINATION_SHIFT = 5
	LIST_LIMIT = 5
	LIST_START_OFFSET = 0
)

type Item interface {
	Id() string
	Compare() int
	Info() string
	Name() string
	Delete(context.Context) error
	Filter(context.Context) ([]Item, error)
}

func buildPagiButtons(total, limit, offset int, entity string) [][]map[string]string {
	if total == 0 {
		return [][]map[string]string{}
	}
	if offset == total {
		return [][]map[string]string{{
			{
				"text": "свернуть",
				"callback_data": models.CallList(strconv.Itoa(LIST_START_OFFSET), "<<<", entity).String(),
			},
		}}
	}
	var keyBoard []map[string]string
	if offset + limit >= total {
		previousButton := map[string]string{"text": "назад", "callback_data": models.CallList(strconv.Itoa(offset), "<<", entity).String()}
		keyBoard = []map[string]string{previousButton}
	} else {
		if offset == 0 {
			nextButton := map[string]string{"text": "вперед", "callback_data": models.CallList(strconv.Itoa(offset), ">>", entity).String()}
			keyBoard = []map[string]string{nextButton}
		} else {
			nextButton := map[string]string{"text": "вперед", "callback_data": models.CallList(strconv.Itoa(offset), ">>", entity).String()}
			previousButton := map[string]string{"text": "назад", "callback_data": models.CallList(strconv.Itoa(offset), "<<", entity).String()}
			keyBoard = []map[string]string{previousButton, nextButton}
		}
	}

	allButton := map[string]string{"text": fmt.Sprintf("показать все (%d)", total), "callback_data": models.CallList(strconv.Itoa(offset), "<>", entity).String()}
	allButtonBar := []map[string]string{allButton}

	markup := [][]map[string]string{}
	if total <= limit {
		return markup
	}

	markup = append(markup, keyBoard)
	markup = append(markup, allButtonBar)

	return markup
}

func buildListButtons[T Item](items []T, limit, offset int) []map[string]string {
	sort.Slice(items, func(i, j int) bool { return comparator(items, i, j) })
	var buttons []map[string]string
	for i, item := range items {
		if offset != len(items) {
			if i == limit + offset {
				break
			}
			if i < offset {
				continue
			}
		}
		button := map[string]string{
			"text": item.Info(),
			"callback_data": models.CallInfo(item.Id(), strconv.Itoa(offset), item.Name()).String(),
		}
		buttons = append(buttons, button)
	}

	return buttons
}

func BuildItemListMarkup[T Item](items []T, limit, offset int, direction, entity string) [][]map[string]string {
	newOffset := offset
	if direction == "<" {

	}
	if direction == "<<<" {
		newOffset = 0
	}
	if direction == ">>" {
		newOffset += LIST_PAGINATION_SHIFT
	} 
	if direction == "<<" {
		newOffset -= LIST_PAGINATION_SHIFT
	}
	if direction == "<>" {
		newOffset = len(items)
	}

	itemsListAsButtons := buildListButtons(items, limit, newOffset)
	pagiButtons := buildPagiButtons(len(items), limit, newOffset, entity)

	markup := [][]map[string]string{}

	for _, button := range itemsListAsButtons {
		markup = append(markup, []map[string]string{button})
	}

	markup = append(markup, pagiButtons...)

	return markup
}

func comparator[T Item](items []T, i, j int) bool {
	countI := items[i].Compare()
	countJ := items[j].Compare()
	return countI < countJ
}
