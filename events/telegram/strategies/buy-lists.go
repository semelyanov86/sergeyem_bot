package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/lists"
	"bot/settings"
	"errors"
	"strconv"
)

const BuyListsCmd = "buylists"

const MsgBuyLists = "Списки из сервиса EasyList:\n"

const AskForListsToken = "Мы не нашли токен EasyList или он не верный. Пожалуйста, предоставьте нам новый 😗"

type BuyListsHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	listService     lists.ListService
}

func NewBuyListsHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, listService lists.ListService) BuyListsHandler {
	return BuyListsHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		listService:     listService,
	}
}

func (h BuyListsHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == BuyListsCmd
}

func (h BuyListsHandler) Handle(msg string, setting *settings.Setting) error {
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

	perPage, err := strconv.Atoi(h.meta.Message.CommandArguments())
	if err != nil {
		perPage = 0
	}

	var text = "<i>" + MsgBuyLists + "</i>"
	listsError := lists.ErrTokenNotExist
	latestLists, err := h.listService.GetAllLists(perPage)

	if errors.Is(listsError, err) {
		err := h.askListsToken(setting)
		if err != nil {
			return e.Wrap("error while asking for EasyList token", err)
		}
		return err
	}
	if len(latestLists) < 1 {
		text = "Не найдено каких-либо списков 🤔"
	}
	if err != nil {
		return e.Wrap("error while getting lists from easylist", err)
	}

	for _, list := range latestLists {
		text = text + list.Id + ". <b>" + list.Attributes.Name + "</b> (" + strconv.Itoa(list.Attributes.ItemsCount) + ")\n"
	}

	msgConfig := h.tg.CreateNewMessage(setting.ChatId, text)
	if err := h.tg.Send(msgConfig); err != nil {
		return e.Wrap("error while sending response with latestLists", err)
	}
	return nil
}

func (h BuyListsHandler) askListsToken(setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, AskedEasylistToken)
	if err != nil {
		return e.Wrap("failed to change mode in AskedEasylistToken", err)
	}

	if err := h.tg.SendMessage(setting.ChatId, AskForListsToken); err != nil {
		return e.Wrap("error while sending response with lists", err)
	}
	return nil
}
