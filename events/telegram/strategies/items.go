package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/lists"
	"bot/settings"
	"errors"
	"strconv"
)

const ItemsCmd = "items"

const MsgErrorToken = "Токен не верный. Получите списки, чтобы получить возможность ввести новый токен."

const MsgErrorWrongList = "Список, который вы передали, не верный. Введите в качестве аргумента корректный идентификатор."

type ItemsHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	listService     lists.ListService
}

func NewItemsHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, listService lists.ListService) ItemsHandler {
	return ItemsHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		listService:     listService,
	}
}

func (h ItemsHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == ItemsCmd
}

func (h ItemsHandler) Handle(msg string, setting *settings.Setting) error {
	setting, err := h.settingsService.GetByUserName(setting.Username)
	if err != nil {
		return e.Wrap("failed to get settings", err)
	}
	if setting == nil {
		err = h.tg.SendMessage(setting.ChatId, NoSettingsError)
		if err != nil {
			return e.Wrap("error while sending response with no setting error", err)
		}
		return e.Wrap("ask for lists with no settings", err)
	}

	listId, err := strconv.Atoi(h.meta.Message.CommandArguments())
	if err != nil {
		h.tg.SendMessage(setting.ChatId, MsgErrorWrongList)
		return e.Wrap("wrong list provided", err)
	}

	if listId < 1 {
		h.tg.SendMessage(setting.ChatId, MsgErrorWrongList)
		return e.Wrap("wrong list provided", err)
	}

	listsError := lists.ErrTokenNotExist
	wrongListError := lists.ErrWrongList
	items, err := h.listService.GetItemsFromList(listId)

	if errors.Is(listsError, err) {
		h.tg.SendMessage(setting.ChatId, MsgErrorToken)
		return e.Wrap("wrong token", err)
	}
	if errors.Is(wrongListError, err) {
		h.tg.SendMessage(setting.ChatId, MsgErrorWrongList)
		return e.Wrap("wrong list provided", err)
	}

	if err != nil {
		return e.Wrap("error while getting items from easylist", err)
	}

	text := h.listService.GenerateMessageFromItems(items)

	msgConfig := h.tg.CreateNewMessage(setting.ChatId, text)
	if err := h.tg.Send(msgConfig); err != nil {
		return e.Wrap("error while sending response with latestLists", err)
	}
	return nil
}
