package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type ChatContext struct {
	ChatId        string   `json:"chatId"`
	UserResponses []string `json:"userResponses"`
	NextHandler   string   `json:"nextHandler"`
}

func newChatContext(chatId string) *ChatContext {
	return &ChatContext{ChatId: chatId, UserResponses: []string{}, NextHandler: ""}
}

func (ctx *ChatContext) AppendText(userResponse string) error {
	ctx.UserResponses = append(ctx.UserResponses, userResponse)
	return nil
}

func (ctx *ChatContext) GetTexts() []string {
	return ctx.UserResponses
}

func (ctx *ChatContext) GetNextHandler() string {
	return ctx.NextHandler
}

func (ctx *ChatContext) SetNextHandler(nextHandler string) string {
	ctx.NextHandler = nextHandler
	return ctx.NextHandler
}

func (ctx *ChatContext) reset() error {
	ctx.NextHandler = ""
	ctx.UserResponses = []string{}
	return nil
}

func (c *Client) GetOrCreateChatContext(chatId string) *ChatContext {
	ctxChatCreation, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	val, err := c.Redis.Get(ctxChatCreation, chatId).Result()

	if err == nil {
		var ctx ChatContext
		if err := json.Unmarshal([]byte(val), &ctx); err == nil {
			return &ctx
		}
	}

	newCtx := newChatContext(chatId)

	jsonCtx, _ := json.Marshal(newCtx)
	c.Redis.Set(ctxChatCreation, chatId, jsonCtx, c.CacheExpiration)

	return newCtx
}

func (c *Client) saveChatContext(ctx *ChatContext) error {
	jsonCtx, err := json.Marshal(ctx)
	if err != nil {
		return err
	}

	ctxSave, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	return c.Redis.Set(ctxSave, ctx.ChatId, jsonCtx, c.CacheExpiration).Err()
}

func (c *Client) SetNextHandler(chatId string, nextHandler string) error {
	if nextHandler == "" {
		return c.Reset(chatId)
	}

	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	ctx.NextHandler = nextHandler
	return c.saveChatContext(ctx)
}

func (c *Client) AppendText(chatId string, text string) error {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	ctx.AppendText(text)
	return c.saveChatContext(ctx)
}

func (c *Client) GetTexts(chatId string) []string {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return nil
	}

	return ctx.GetTexts()
}

func (c *Client) Reset(chatId string) error {
	ctx := c.GetOrCreateChatContext(chatId)
	if ctx == nil {
		return errors.New("chat context not found")
	}

	ctx.reset()
	return c.saveChatContext(ctx)
}
