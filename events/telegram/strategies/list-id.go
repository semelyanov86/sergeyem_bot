package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/lists"
	"bot/settings"
	"strconv"
)

const MsgListId = "Идентификатор успешно сохранён. Вы можете пользоваться сервисом EasyList 👍"

const MsgErrorListId = "Невозможно сохранить ID. Произошла непредвиденная ошибка 😔 "

type ListIdHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	listService     lists.ListService
}

func NewListIdHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, listService lists.ListService) ListIdHandler {
	return ListIdHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		listService:     listService,
	}
}

func (h ListIdHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == AskedEasylistId
}

func (h ListIdHandler) Handle(text string, setting *settings.Setting) error {
	listId, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		h.tg.SendMessage(setting.ChatId, MsgErrorListId+"Wrong ID is provided")
		return e.Wrap("wrong id provided", err)
	}
	setting.EasylistId = listId
	err = h.settingsService.UpdateSetting(setting)
	if err != nil {
		err := h.tg.SendMessage(setting.ChatId, MsgErrorListId+"while saving setting")
		if err != nil {
			return err
		}
		return e.Wrap("failed to save new settings", err)
	}

	err = h.settingsService.ChangeMode(setting.Username, Root)
	if err != nil {
		return e.Wrap("failed to change mode to root", err)
	}
	return h.tg.SendMessage(setting.ChatId, MsgListId)
}
