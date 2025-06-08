package callbackdata

import "strings"

type callbackData struct {
	ID      string
	Command string
}

func (b *Builder) Build(id, command string) *callbackData {
	return &callbackData{
		ID:      id,
		Command: command,
	}
}

func (b *Builder) FromString(raw string) *callbackData {
	params := strings.Split(raw, ";")
	return &callbackData{
		Command: params[0],
		ID:      params[1],
	}
}

func (cd *callbackData) String() string {
	separator := ";"
	return strings.Join(
		[]string{
			cd.Command,
			cd.ID,
		},
		separator,
	)
}
