package db

type Schedule struct {
	BaseFields

	// chatId - id of chat with user, bot uses it to send notification
	ChatId string

	// telegram user id
	OwnerId int

	// pyload
	Text string

	// period
	Delta string

	// day when schedule is originate
	Day string
}

func NewSchedule(ownerId int, chatId, text, delta, day string) *Schedule {
	return &Schedule{
		NewBaseFields(),
		chatId,
		ownerId,
		text,
		delta,
		day,
	}
}
