package db

type User struct {
	// telegram user -> bot's user

	BaseFields

	TGId       int // telegram user id
	Name       string
	TGusername string
	ChatId     int // chatId - id of chat with user, bot uses it to send notification
}
