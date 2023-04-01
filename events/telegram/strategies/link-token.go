package strategies

import (
	telegram2 "bot/clients/telegram"
	"bot/events"
	"bot/lib/e"
	"bot/links"
	"bot/settings"
)

const MsgLinkToken = `Токен успешно сохранён. Теперь вы можете дальше пользоваться сервисом LinkAce 👍`
const MsgErrorLinkToken = "Невозможно сохранить токен. Произошла непредвиденная ошибка 😔 "

type LinkTokenHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              *telegram2.Client
	linkService     links.LinkService
}

func NewLinkTokenHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg *telegram2.Client, linkService links.LinkService) LinkTokenHandler {
	return LinkTokenHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		linkService:     linkService,
	}
}

func (h LinkTokenHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == LinksToken
}

func (h LinkTokenHandler) Handle(text string, setting *settings.Setting) error {
	setting.LinkaceToken = text
	err := h.settingsService.UpdateSetting(setting)
	if err != nil {
		err := h.tg.SendMessage(setting.ChatId, MsgErrorLinkToken+"while saving setting")
		if err != nil {
			return err
		}
		return e.Wrap("failed to save new settings", err)
	}

	err = h.settingsService.ChangeMode(setting.Username, Root)
	if err != nil {
		return e.Wrap("failed to change mode to root", err)
	}
	return h.tg.SendMessage(setting.ChatId, MsgLinkToken)
}
