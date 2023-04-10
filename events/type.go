package events

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Fetcher interface {
	Fetch(limit int) ([]Event[TelegramMeta], error)
	ListenWebhooks(webhookUrl string) ([]Event[TelegramMeta], error)
}

type Processor interface {
	Process(e Event[TelegramMeta]) error
}

type Client interface {
	Updates(offset int, limit int) ([]tgbotapi.Update, error)
	SendMessage(chatID int64, text string) error
	Send(msg tgbotapi.MessageConfig) error
	CreateNewMessage(chatID int64, text string) tgbotapi.MessageConfig
	GetMessage() string
	GetToken() string
	Request(chat tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	GetWebhookInfo() (tgbotapi.WebhookInfo, error)
	ListenWebhook(url string) tgbotapi.UpdatesChannel
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
