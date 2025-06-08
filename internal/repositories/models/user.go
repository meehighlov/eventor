package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	TgID       string    `gorm:"column:tg_id;type:varchar(255);unique;not null"`
	TgUsername string    `gorm:"column:tg_username;type:varchar(255)"`
	TgChatID   string    `gorm:"column:chat_id;type:varchar(255)"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()"`
}

func (*User) TableName() string {
	return "user"
}
