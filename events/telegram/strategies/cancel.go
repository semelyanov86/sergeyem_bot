package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CancelCmd = "cancel"

const MsgCancel = `Последняя команда отменена 👌`

type CancelHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
}

func NewCancelHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client) CancelHandler {
	return CancelHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
	}
}

func (h CancelHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == CancelCmd
}

func (h CancelHandler) Handle(text string, setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, Root)
	if err != nil {
		return e.Wrap(RootModeFailed, err)
	}
	err = h.settingsService.SetContext(setting.Username, "")
	if err != nil {
		return e.Wrap(ContextClearFailed, err)
	}
	_ = tgbotapi.NewRemoveKeyboard(true)
	return h.tg.SendMessage(setting.ChatId, MsgCancel)
}
