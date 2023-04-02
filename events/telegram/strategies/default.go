package strategies

import (
	"bot/events"
	"bot/settings"
)

const MsgDefault = "Мне не удалось распознать вашу команду. Пожалуйста, исправьте ошибку и попробуйте снова ✍️"

type DefaultHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
}

func NewDefaultHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client) DefaultHandler {
	return DefaultHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
	}
}

func (h DefaultHandler) IsSupported(mode int) bool {
	return true
}

func (h DefaultHandler) Handle(text string, setting *settings.Setting) error {
	return h.tg.SendMessage(setting.ChatId, MsgDefault)
}
