package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ELEMENT_STATUS_OPEN = "open"
	ELEMENT_STATUS_IN_PROGRESS = "in_progress"
	ELEMENT_STATUS_MASTERED = "mastered"
)

type Element struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	Name         string    `gorm:"column:name;type:varchar(255);not null"`
	TutorialLink string    `gorm:"column:tutorial_link;type:varchar(255)"`
	Status       string    `gorm:"column:status;type:varchar(50);not null;default:'open'"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()"`
}

func (*Element) TableName() string {
	return "element"
}

func (e *Element) NextStatus() string {
	switch e.Status {
	case ELEMENT_STATUS_OPEN:
		return ELEMENT_STATUS_IN_PROGRESS
	case ELEMENT_STATUS_IN_PROGRESS:
		return ELEMENT_STATUS_MASTERED
	default:
		return ELEMENT_STATUS_OPEN
	}
}

func (e *Element) ElementReadableStatus(status string) string {
	stat := e.Status
	if status != "" {
		stat = status
	}
	switch stat {
	case ELEMENT_STATUS_OPEN:
		return "✨ Открыт"
	case ELEMENT_STATUS_IN_PROGRESS:
		return "🔄 В процессе"
	default:
		return "✅ Освоен"
	}
}
