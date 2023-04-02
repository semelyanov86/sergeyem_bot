package telegram

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processor struct {
	Tg              events.Client
	offset          int
	SettingsService settings.ServiceInterface
	config          settings.Config
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func New(client events.Client, config settings.Config) *Processor {
	return &Processor{
		Tg:     client,
		offset: 0,
		config: config,
		SettingsService: &settings.Service{
			Repository: &settings.Repository{
				DB: config.Db.Sql,
			},
		},
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event[events.TelegramMeta], error) {
	updates, err := p.Tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can not get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event[events.TelegramMeta], 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].UpdateID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event[events.TelegramMeta]) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can not process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event2 events.Event[events.TelegramMeta]) error {
	meta := event2.Meta
	if meta.Username == "" {
		return e.Wrap("can not process message due to empty meta", ErrUnknownMetaType)
	}

	if err := p.doCmd(event2.Text, meta); err != nil {
		return e.Wrap("can not process message", err)
	}

	return nil
}

func event(u tgbotapi.Update) events.Event[events.TelegramMeta] {
	updType := fetchType(u)
	res := events.Event[events.TelegramMeta]{
		Type: updType,
		Text: fetchText(u),
	}
	if u.Message != nil {
		res.Meta = events.TelegramMeta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.Chat.UserName,
			Message:  u.Message,
		}
	}

	return res
}

func fetchText(u tgbotapi.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}

func fetchType(u tgbotapi.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func (p *Processor) getSettingsModel(userName string) (*settings.Setting, error) {
	return p.SettingsService.GetByUserName(userName)
}
