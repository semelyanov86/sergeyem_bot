package strategies

import (
	telegram2 "bot/clients/telegram"
	"bot/events"
	"bot/lib/e"
	"bot/links"
	"bot/settings"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

const ListsCmd = "lists"

const MsgLists = "Доступные списки:\n"

type ListsHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              *telegram2.Client
	linkService     links.LinkService
}

func NewListsHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg *telegram2.Client, linkService links.LinkService) ListsHandler {
	return ListsHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		linkService:     linkService,
	}
}

func (h ListsHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == ListsCmd
}

func (h ListsHandler) Handle(msg string, setting *settings.Setting) error {
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

	var text = "<i>" + MsgLists + "</i>"
	linksError := links.ErrTokenNotExist
	allLists, err := h.linkService.GetAllLists()
	if errors.Is(linksError, err) {
		h.tg.SendMessage(setting.ChatId, LinkAceTokenError)
		if err != nil {
			return e.Wrap("error while asking for linkace token", err)
		}
	}
	if len(allLists) < 0 {
		text = MsgNoList
	}
	if err != nil {
		return e.Wrap("error while getting allLists", err)
	}

	for _, list := range allLists {
		text = text + strconv.Itoa(list.Id) + ". <b>" + list.Name + "</b>\n" + list.Description + "\n"
	}

	msgConfig := h.tg.CreateNewMessage(setting.ChatId, text)
	if err := h.tg.Send(msgConfig); err != nil {
		return e.Wrap("error while sending response with allLists", err)
	}
	return nil
}

func GenerateListsButtons(lists []links.List, passIds bool) [][]tgbotapi.KeyboardButton {
	rows := tgbotapi.NewKeyboardButtonRow()
	row := make([]tgbotapi.KeyboardButton, 0)
	allRows := make([][]tgbotapi.KeyboardButton, 0)
	for i, list := range lists {
		if i > 0 && i%3 == 0 {
			rows = tgbotapi.NewKeyboardButtonRow(row...)
			allRows = append(allRows, rows)
			row = make([]tgbotapi.KeyboardButton, 0)
		}
		title := list.Name
		if passIds {
			title = title + "|" + strconv.Itoa(list.Id)
		}
		row = append(row, tgbotapi.NewKeyboardButton(title))
	}
	if len(row) > 0 {
		rows = tgbotapi.NewKeyboardButtonRow(row...)
		allRows = append(allRows, rows)
	}
	return allRows
}
