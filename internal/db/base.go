package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Entity interface {
	Id() string
	Compare() int
	Info() string
	Name() string
	Save(context.Context) error
	Delete(context.Context) error
	Filter(context.Context) ([]Entity, error)
}

type BaseFields struct {
	ID        string
	CreatedAt string
	UpdatedAt string
}

func (b *BaseFields) RefresTimestamps() (created string, updated string, _ error) {
	now := time.Now().Format("02.01.2006T15:04:05")
	if b.CreatedAt == "" {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	return b.CreatedAt, b.UpdatedAt, nil
}

func NewBaseFields() BaseFields {
	now := time.Now().Format("02.01.2006T15:04:05")
	return BaseFields{uuid.New().String(), now, now}
}
