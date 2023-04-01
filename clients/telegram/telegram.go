package telegram

import (
	"bot/lib/e"
	"log"
)
import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Client struct {
	host     string
	basePath string
	client   *tgbotapi.BotAPI
}

func New(host string, token string) *Client {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   bot,
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]tgbotapi.Update, error) {
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(offset)
	updateConfig.Limit = limit

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	updates, err := c.client.GetUpdates(updateConfig)
	if err != nil {
		return nil, e.Wrap("can not do request to receive updates", err)
	}
	return updates, nil
}

func (c *Client) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	// Okay, we're sending our message off! We don't care about the message
	// we just sent, so we'll discard it.
	if _, err := c.client.Send(msg); err != nil {
		return e.Wrap("can not send message", err)
	}

	return nil
}

func (c *Client) Send(msg tgbotapi.MessageConfig) error {
	if _, err := c.client.Send(msg); err != nil {
		return e.Wrap("can not send message", err)
	}
	return nil
}

func (c *Client) CreateNewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msgConfig := tgbotapi.NewMessage(chatID, text)

	msgConfig.ParseMode = tgbotapi.ModeHTML

	return msgConfig
}
