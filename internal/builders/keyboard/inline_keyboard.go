package keyboard

type Button struct {
	text                            string
	callback_data                   string
	copy_text                       map[string]string
	switch_inline_query             string
	switch_inline_query_chosen_chat map[string]interface{}
	url                             string
	button_type                     string
}

func (b *inlineKeyboard) NewButton(text, callback_data string) *Button {
	return &Button{text: text, callback_data: callback_data, button_type: "button"}
}

func (b *inlineKeyboard) NewCopyButton(text, copy_text string) *Button {
	return &Button{
		text:        text,
		copy_text:   map[string]string{"text": copy_text},
		button_type: "copy",
	}
}

func (b *inlineKeyboard) NewSwitchInlineButton(text string) *Button {
	return &Button{
		text:                text,
		switch_inline_query: "",
		button_type:         "switch_inline_query",
	}
}

func (b *inlineKeyboard) NewAddToChatButton(text, query string) *Button {
	return &Button{
		text: text,
		switch_inline_query_chosen_chat: map[string]interface{}{
			"allow_user_chats":    false,
			"allow_bot_chats":     false,
			"allow_group_chats":   true,
			"allow_channel_chats": false,
			"query":               query,
		},
		button_type: "switch_inline_query_chosen_chat",
	}
}

func (b *inlineKeyboard) NewURLButton(text, url string) *Button {
	return &Button{
		text:        text,
		url:         url,
		button_type: "url",
	}
}

func (b *inlineKeyboard) NewShareLinkButton(text, link, description string) *Button {
	query := link
	if description != "" {
		query = description + "\n\n" + link
	}
	return &Button{
		text: text,
		switch_inline_query_chosen_chat: map[string]interface{}{
			"allow_user_chats":    true,
			"allow_bot_chats":     true,
			"allow_group_chats":   true,
			"allow_channel_chats": true,
			"query":               query,
		},
		button_type: "switch_inline_query_chosen_chat",
	}
}

func (b *Button) Raw() *map[string]interface{} {
	result := map[string]interface{}{"text": b.text}

	if b.callback_data != "" {
		result["callback_data"] = b.callback_data
	}

	if b.copy_text != nil {
		result["copy_text"] = b.copy_text
	}

	if b.isInlineSwitchButton() {
		result["switch_inline_query"] = b.switch_inline_query
	}

	if b.isChosenChatButton() {
		result["switch_inline_query_chosen_chat"] = b.switch_inline_query_chosen_chat
	}

	if b.isURLButton() {
		result["url"] = b.url
	}

	return &result
}

func (b *Button) isInlineSwitchButton() bool {
	return b.button_type == "switch_inline_query"
}

func (b *Button) isChosenChatButton() bool {
	return b.button_type == "switch_inline_query_chosen_chat"
}

func (b *Button) isURLButton() bool {
	return b.button_type == "url"
}

type inlineKeyboard struct {
	markup []*[]map[string]interface{}
}

func (b *Builder) BuildInlineKeyboard() *inlineKeyboard {
	return &inlineKeyboard{markup: []*[]map[string]interface{}{}}
}

// appends button list to representation of keyboard to new row below
func (ik *inlineKeyboard) AppendAsLine(buttons ...*Button) {
	rawButtons := []map[string]interface{}{}
	for _, button := range buttons {
		rawButtons = append(rawButtons, *button.Raw())
	}

	ik.markup = append(ik.markup, &rawButtons)
}

// appends button list as stacked lines
func (ik *inlineKeyboard) AppendAsStack(buttons ...*Button) {
	for _, button := range buttons {
		ik.AppendAsLine(button)
	}
}

func (ik *inlineKeyboard) Murkup() []*[]map[string]interface{} {
	return ik.markup
}
