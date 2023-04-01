package events

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Fetcher interface {
	Fetch(limit int) ([]Event[TelegramMeta], error)
}

type Processor interface {
	Process(e Event[TelegramMeta]) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event[Meta TelegramMeta] struct {
	Type Type
	Text string
	Meta Meta
}

type TelegramMeta struct {
	ChatID   int64
	Username string
	Message  *tgbotapi.Message
}
