package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
)

const StartCmd = "start"

type StartHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
}

func NewStartHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client) StartHandler {
	return StartHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
	}
}

func (h StartHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == StartCmd
}

func (h StartHandler) Handle(text string, setting *settings.Setting) error {
	_, err := h.settingsService.GetOrCreateSetting(setting.Username, setting.ChatId)

	if err != nil {
		return e.Wrap("Failed to get settings", err)
	}
	return h.tg.SendMessage(setting.ChatId, MsgHelp)
}
