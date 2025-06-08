package short_id

import (
	"math/rand"
	"time"
)

type Builder struct {
	length  int
	charset string
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func New(length int) *Builder {
	return &Builder{
		length:  length,
		charset: charset,
	}
}

func (s *Builder) Build() string {
	rand.NewSource(time.Now().UnixNano())
	b := make([]byte, s.length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
