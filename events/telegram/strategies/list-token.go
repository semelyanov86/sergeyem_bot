package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/lists"
	"bot/settings"
)

const MsgListToken = "Токен успешно сохранён. Теперь вы можете пользоваться сервисом EasyList 👍"

const MsgErrorListToken = "Невозможно сохранить токен. Произошла непредвиденная ошибка 😔 "

const MsgAskEasyListId = "Введите идентификатор пользователя в сервисе EasyList"

type ListTokenHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	listService     lists.ListService
}

func NewListTokenHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, listService lists.ListService) ListTokenHandler {
	return ListTokenHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		listService:     listService,
	}
}

func (h ListTokenHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == AskedEasylistToken
}

func (h ListTokenHandler) Handle(text string, setting *settings.Setting) error {
	setting.EasylistToken = text
	err := h.settingsService.UpdateSetting(setting)
	if err != nil {
		err := h.tg.SendMessage(setting.ChatId, MsgErrorListToken+"while saving setting")
		if err != nil {
			return err
		}
		return e.Wrap("failed to save new settings", err)
	}
	if setting.EasylistId != 0 {
		err = h.settingsService.ChangeMode(setting.Username, Root)
		if err != nil {
			return e.Wrap("failed to change mode to root", err)
		}
		return h.tg.SendMessage(setting.ChatId, MsgListToken)
	}

	err = h.settingsService.ChangeMode(setting.Username, AskedEasylistId)
	if err != nil {
		return e.Wrap("failed to change mode to AskedEasylistId", err)
	}
	return h.tg.SendMessage(setting.ChatId, MsgAskEasyListId)
}
